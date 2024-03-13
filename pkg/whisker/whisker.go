package whisker

import (
    "fmt"   // Import the fmt package for formatted I/O
    "io/ioutil" // Import the ioutil package for file I/O utilities
    "math"  // Import the math package for mathematical functions

    "github.com/Aanthord/remy-go/pkg/dna"    // Import the dna package from the remy project
    "github.com/Aanthord/remy-go/pkg/memory" // Import the memory package from the remy project
)

// Whisker represents a single congestion control rule
type Whisker struct {
    Generation      uint            // Generation of the whisker
    WindowIncrement int             // Window increment value
    WindowMultiple  float64         // Window multiple value
    Intersend       float64         // Intersend time
    Domain          *memory.MemoryRange // Memory range associated with the whisker
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
    for generation := uint(0); generation < config.Generations; generation++ {
        for _, windowIncrement := range config.WindowIncrements {
            for _, windowMultiple := range config.WindowMultiples {
                for _, intersend := range config.Intersends {
                    for _, domain := range config.Domains {
                        whisker := NewWhisker(generation, windowIncrement, windowMultiple, intersend, domain)
                        whiskers = append(whiskers, whisker)
                    }
                }
            }
        }
    }

    return whiskers
}

// LoadWhiskers loads whiskers from a file specified by the filename
func LoadWhiskers(filename string) ([]*Whisker, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var whiskers []*Whisker

    // Unmarshal the whiskers from the data
    dnaWhiskers := &dna.Whiskers{}
    err = dnaWhiskers.Unmarshal(data)
    if err != nil {
        return nil, err
    }

    // Convert dna.Whisker to Whisker
    for _, dnaWhisker := range dnaWhiskers.Whiskers {
        domain := &memory.MemoryRange{
            Lower: memory.FromDNAMemory(dnaWhisker.Domain.Lower),
            Upper: memory.FromDNAMemory(dnaWhisker.Domain.Upper),
        }
        whisker := NewWhisker(dnaWhisker.Generation, int(dnaWhisker.WindowIncrement), dnaWhisker.WindowMultiple, dnaWhisker.Intersend, domain)
        whiskers = append(whiskers, whisker)
    }

    return whiskers, nil
}

// SaveWhiskers saves a slice of whiskers to a file specified by the filename
func SaveWhiskers(whiskers []*Whisker, filename string) error {
    dnaWhiskers := &dna.Whiskers{}

    // Convert Whisker to dna.Whisker
    for _, whisker := range whiskers {
        dnaWhisker := &dna.Whisker{
            Generation:      whisker.Generation,
            WindowIncrement: int32(whisker.WindowIncrement),
            WindowMultiple:  whisker.WindowMultiple,
            Intersend:       whisker.Intersend,
            Domain: &dna.MemoryRange{
                Lower: whisker.Domain.Lower.ToDNAMemory(),
                Upper: whisker.Domain.Upper.ToDNAMemory(),
            },
        }
        dnaWhiskers.Whiskers = append(dnaWhiskers.Whiskers, dnaWhisker)
    }

    // Marshal the whiskers to data
    data, err := dnaWhiskers.Marshal()
    if err != nil {
        return err
    }

    // Write the data to the file
    err = ioutil.WriteFile(filename, data, 0644)
    if err != nil {
        return err
    }

    return nil
}
