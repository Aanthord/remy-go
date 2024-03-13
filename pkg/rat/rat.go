package rat

import (
    "fmt"   // Import the fmt package for formatted I/O
    "math/rand" // Import the rand package for random number generation
    "time"  // Import the time package for time-related operations

    "github.com/Aanthord/remy-go/pkg/memory" // Import the memory package from the remy project
    "github.com/Aanthord/remy-go/pkg/whisker" // Import the whisker package from the remy project
)

// RAT represents the Remy Augmented TCP (RAT) congestion control algorithm
type RAT struct {
    Whiskers *whisker.WhiskerTree // Pointer to the WhiskerTree that holds the whiskers
    Memory   *memory.Memory       // Pointer to the Memory object that stores the network state
    Sender   *Sender              // Pointer to the Sender object
    Receiver *Receiver            // Pointer to the Receiver object

    CurrentWhisker *whisker.Whisker // Pointer to the currently selected Whisker
    LastSendTime   time.Time        // Timestamp of the last packet sent
    FlowSize       int              // Number of packets sent in the current flow
    RTT            time.Duration    // Round-Trip Time (RTT) of the network
}

// NewRAT is a constructor that creates a new instance of the RAT struct
func NewRAT(whiskers *whisker.WhiskerTree) *RAT {
    return &RAT{
        Whiskers: whiskers,
        Memory:   memory.NewMemory(),
    }
}

// Start initializes the RAT algorithm with the given Sender and Receiver objects
func (rat *RAT) Start(sender *Sender, receiver *Receiver) {
    rat.Sender = sender
    rat.Receiver = receiver
    rat.LastSendTime = time.Now()
    rat.FlowSize = 0
    rat.RTT = 0
    rat.selectWhisker()
}

// selectWhisker selects the appropriate whisker based on the current memory state
func (rat *RAT) selectWhisker() {
    var err error
    rat.CurrentWhisker, err = rat.Whiskers.FindWhisker(rat.Memory)
    if err != nil {
        fmt.Printf("Error selecting whisker: %v\n", err)
        rat.CurrentWhisker = rat.Whiskers.Root.Whisker
    }
    rat.FlowSize = 1
}

// OnPacketSent is called when a packet is sent
func (rat *RAT) OnPacketSent(seq int) {
    rat.FlowSize++
    rat.LastSendTime = time.Now()
}

// OnPacketAcked is called when a packet is acknowledged by the receiver
func (rat *RAT) OnPacketAcked(seq int, rtt time.Duration) {
    rat.RTT = rtt
    rat.Memory.UpdateSendRate(float64(rat.FlowSize) / rtt.Seconds())
    rat.Memory.UpdateRecvRate(float64(rat.FlowSize) / rtt.Seconds())
    rat.Memory.UpdateLatestDelay(rtt.Seconds())
    rat.Memory.UpdateInterPacketDelay(rtt.Seconds() / float64(rat.FlowSize))
    rat.selectWhisker()
    rat.FlowSize = 0
}

// OnPacketLost is called when a packet is lost
func (rat *RAT) OnPacketLost() {
    rat.FlowSize = 1
    rat.Memory.UpdateSendRate(0)
    rat.Memory.UpdateRecvRate(0)
    rat.selectWhisker()
}

// TimeToSend determines whether it is time to send the next packet
func (rat *RAT) TimeToSend() bool {
    sendInterval := time.Duration(rat.CurrentWhisker.Intersend * float64(time.Second))
    return time.Since(rat.LastSendTime) >= sendInterval
}

// GetSendRate returns the current send rate
func (rat *RAT) GetSendRate() float64 {
    return rat.Memory.GetSendRate()
}

// GetCongestionWindow returns the congestion window size calculated by the current whisker
func (rat *RAT) GetCongestionWindow() int {
    return rat.CurrentWhisker.Window(rat.FlowSize)
}

// Sender and Receiver are placeholder types that represent the sender and receiver of network packets
type Sender struct{}
type Receiver struct{}
