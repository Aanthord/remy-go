package sender

import (
    "fmt"   // Import the fmt package for formatted I/O
    "sync"  // Import the sync package for synchronization primitives
    "time"  // Import the time package for time-related operations

    "github.com/Aanthord/remy-go/pkg/rat" // Import the rat package from the remy project
)

// SenderGang represents a group of senders
type SenderGang struct {
    Senders       []*Sender     // Array of Sender objects
    NumSenders    int           // Total number of senders
    SenderFactory func(id int) *Sender // Function to create new Sender instances
    Mu            sync.Mutex    // Mutex for synchronization
    StopChan      chan struct{} // Channel to signal senders to stop
}

// NewSenderGang is a constructor that creates a new instance of the SenderGang struct
func NewSenderGang(numSenders int, senderFactory func(id int) *Sender) *SenderGang {
    return &SenderGang{
        Senders:       make([]*Sender, numSenders),
        NumSenders:    numSenders,
        SenderFactory: senderFactory,
        StopChan:      make(chan struct{}),
    }
}

// Start starts all the senders in the gang
func (sg *SenderGang) Start() {
    for i := 0; i < sg.NumSenders; i++ {
        sg.Senders[i] = sg.SenderFactory(i)
        go sg.runSender(sg.Senders[i])
    }
}

// Stop stops all the senders in the gang
func (sg *SenderGang) Stop() {
    close(sg.StopChan)
}

// runSender runs a single sender
func (sg *SenderGang) runSender(sender *Sender) {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            sg.Mu.Lock()
            err := sender.Send([]byte(fmt.Sprintf("Data from sender %d", sender.ID)))
            if err != nil {
                fmt.Printf("Sender %d: Error sending data: %v\n", sender.ID, err)
            }
            sg.Mu.Unlock()
        case <-sg.StopChan:
            return
        }
    }
}

// OnAck is called when an acknowledgment (ACK) is received for a specific sender
func (sg *SenderGang) OnAck(senderId int, ack *Ack) {
    sg.Mu.Lock()
    defer sg.Mu.Unlock()

    if senderId >= 0 && senderId < sg.NumSenders {
        sg.Senders[senderId].OnAck(ack)
    }
}

// OnTimeout is called when a timeout occurs for a specific sender
func (sg *SenderGang) OnTimeout(senderId int) {
    sg.Mu.Lock()
    defer sg.Mu.Unlock()

    if senderId >= 0 && senderId < sg.NumSenders {
        sg.Senders[senderId].OnTimeout()
    }
}

// SenderFactory is an example implementation of the sender factory
func SenderFactory(id int) *Sender {
    rat := rat.NewRAT(nil) // Initialize with nil whiskers
    return NewSender(id, rat)
}
