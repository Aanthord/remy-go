package sender

import (
    "fmt"   // Import the fmt package for formatted I/O
    "time"  // Import the time package for time-related operations

    "github.com/Aanthord/remy-go/pkg/memory" // Import the memory package from the remy project
    "github.com/Aanthord/remy-go/pkg/rat"    // Import the rat package from the remy project
)

// Sender represents a sender in the network
type Sender struct {
    ID            int              // Unique identifier of the sender
    Rat           *rat.RAT         // Pointer to the RAT (Remy Augmented TCP) congestion control algorithm instance
    SendRate      float64          // Current sending rate of the sender
    CongestionWnd int              // Current congestion window size
    LastSendTime  time.Time        // Timestamp of the last packet sent
    LastAckTime   time.Time        // Timestamp of the last acknowledgment received
    BytesSent     int              // Total number of bytes sent by the sender
    BytesAcked    int              // Total number of bytes acknowledged by the receiver
    SeqNo         int              // Sequence number of the last sent packet
}

// NewSender is a constructor that creates a new instance of the Sender struct
func NewSender(id int, rat *rat.RAT) *Sender {
    return &Sender{
        ID:            id,
        Rat:           rat,
        SendRate:      0,
        CongestionWnd: 1,
        LastSendTime:  time.Now(),
        LastAckTime:   time.Now(),
        BytesSent:     0,
        BytesAcked:    0,
        SeqNo:         0,
    }
}

// Send sends data from the sender
func (s *Sender) Send(data []byte) error {
    // Check if it's too early to send based on the current send rate
    if time.Since(s.LastSendTime) < time.Duration(1.0/s.SendRate)*time.Second {
        return fmt.Errorf("Sender %d: Too early to send, rate limiting", s.ID)
    }

    // Check if the congestion window limit is reached
    if s.BytesSent-s.BytesAcked >= s.CongestionWnd {
        return fmt.Errorf("Sender %d: Congestion window limit reached", s.ID)
    }

    // Send the data
    s.SeqNo++
    s.BytesSent += len(data)
    s.LastSendTime = time.Now()
    s.Rat.OnPacketSent(s.SeqNo)

    return nil
}

// OnAck is called when an acknowledgment (ACK) is received
func (s *Sender) OnAck(ack *Ack) {
    s.BytesAcked += ack.BytesAcked
    s.LastAckTime = time.Now()
    rtt := time.Since(ack.SentTime)
    s.Rat.OnPacketAcked(ack.SeqNo, rtt)
    s.UpdateSendRate()
    s.UpdateCongestionWnd()
}

// OnTimeout is called when a timeout occurs, indicating a packet loss
func (s *Sender) OnTimeout() {
    s.Rat.OnPacketLost()
    s.UpdateSendRate()
    s.UpdateCongestionWnd()
}

// UpdateSendRate updates the send rate of the sender
func (s *Sender) UpdateSendRate() {
    s.SendRate = s.Rat.GetSendRate()
}

// UpdateCongestionWnd updates the congestion window size of the sender
func (s *Sender) UpdateCongestionWnd() {
    s.CongestionWnd = s.Rat.GetCongestionWindow()
}

// Ack represents an acknowledgment received by the sender
type Ack struct {
    SeqNo      int       // Sequence number of the acknowledged packet
    BytesAcked int       // Number of bytes acknowledged
    SentTime   time.Time // Timestamp when the acknowledged packet was sent
}
