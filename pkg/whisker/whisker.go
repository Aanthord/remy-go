package whisker

import (
    "fmt"
    "io/ioutil"
    "math"

    "google.golang.org/protobuf/proto"
    "github.com/Aanthord/remy-go/pkg/dna"
    "github.com/Aanthord/remy-go/pkg/memory"
)

// Whisker represents a single congestion control rule
type Whisker struct {
    Generation      uint
    WindowIncrement int
    WindowMultiple  float64
    Intersend       float64
    Domain          *memory.MemoryRange
}

// NewWhisker is a constructor that creates a new instance of the Whisker struct
func NewWhisker(generation uint, windowIncrement int, windowMultiple, intersend float64, domain *memory.MemoryRange) *Whisker {
    return &Whisker{
        Generation:      generation,
        WindowIncrement: windowIncrement,
        WindowMultiple:  windowMultiple,
        Intersend:       intersend,
        Domain:          domain,
    }
}

// Window calculates the new congestion window size based on the previous window and the whisker's properties
func (w *Whisker) Window(prevWindow uint) uint {
    return uint(math.Max(0, math.Min(float64(prevWindow)*w.WindowMultiple+float64(w.WindowIncrement), 1000000)))
}

// String returns a string representation of the whisker
func (w *Whisker) String() string {
    return fmt.Sprintf("Generation=%d, WindowIncrement=%d, WindowMultiple=%f, Intersend=%f, Domain=%v",
        w.Generation, w.WindowIncrement, w.WindowMultiple, w.Intersend, w.Domain)
}

// GenerateWhiskers generates a slice of whiskers based on the provided configuration
func GenerateWhiskers(config *dna.ConfigRange) []*Whisker {
    var whiskers []*Whisker

    // Generate whiskers based on the configuration
    for generation := uint(0); generation < uint(config.Generations); generation++ {
        for _, windowIncrement := range config.WindowIncrements {
            for _, windowMultiple := range config.WindowMultiples {
                for _, intersend := range config.Intersends {
                    for _, domain := range config.Domains {
                        memoryDomain := memory.NewMemoryRange(memory.FromDNAMemory(domain.Lower), memory.FromDNAMemory(domain.Upper))
                        whisker := NewWhisker(generation, int(windowIncrement), float64(windowMultiple), float64(intersend), memoryDomain)
                        whiskers = append(whiskers, whisker)
                    }
                }
            }
        }
    }

    return whiskers
}

// LoadWhiskers loads whiskers from a file and constructs a WhiskerTree
func LoadWhiskers(filename string) (*WhiskerTree, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    // Unmarshal the whiskers from the data
    dnaWhiskers := &dna.Whiskers{}
    if err := proto.Unmarshal(data, dnaWhiskers); err != nil {
        return nil, err
    }

    // Create a new WhiskerTree
    tree := NewWhiskerTree()

    // Insert the loaded whiskers into the WhiskerTree
    for _, dnaWhisker := range dnaWhiskers.Whiskers {
        domain := memory.NewMemoryRange(memory.FromDNAMemory(dnaWhisker.Domain.Lower), memory.FromDNAMemory(dnaWhisker.Domain.Upper))
        whisker := NewWhisker(uint(dnaWhisker.Generation), int(dnaWhisker.WindowIncrement), float64(dnaWhisker.WindowMultiple), float64(dnaWhisker.Intersend), domain)
        if err := tree.Insert(whisker); err != nil {
            return nil, err
        }
    }

    return tree, nil
}

func SaveWhiskers(whiskers []*Whisker, filename string) error {
    dnaWhiskers := &dna.Whiskers{}

    // Convert Whisker to dna.Whisker
    for _, whisker := range whiskers {
        dnaWhisker := &dna.Whisker{
            Generation:      uint32(whisker.Generation),
            WindowIncrement: uint32(whisker.WindowIncrement),
            WindowMultiple:  float32(whisker.WindowMultiple),
            Intersend:       float32(whisker.Intersend),
            Domain: &dna.MemoryRange{
                Lower: whisker.Domain.Lower.ToDNAMemory(),
                Upper: whisker.Domain.Upper.ToDNAMemory(),
            },
        }
        dnaWhiskers.Whiskers = append(dnaWhiskers.Whiskers, dnaWhisker)
    }

    // Marshal the whiskers to data
    data, err := proto.Marshal(dnaWhiskers)
    if err != nil {
        return err
    }
    // Write the data to the file
    if err := ioutil.WriteFile(filename, data, 0644); err != nil {
        return err
    }

    return nil
}

// ParseConfig parses the configuration data and returns a dna.ConfigRange instance
func ParseConfig(configData []byte) (*dna.ConfigRange, error) {
    config := &dna.ConfigRange{}
    err := proto.Unmarshal(configData, config)
    if err != nil {
        return nil, err
    }
    return config, nil
}
