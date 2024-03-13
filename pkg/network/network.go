package network

import (
    "math/rand" // Import the rand package for random number generation
    "time"      // Import the time package for time-related operations

    "github.com/Aanthord/remy-go/pkg/memory" // Import the memory package from the remy project
    "github.com/Aanthord/remy-go/pkg/rat"    // Import the rat package from the remy project
    "github.com/Aanthord/remy-go/pkg/sender" // Import the sender package from the remy project
)

// Network represents the simulated network environment
type Network struct {
    Senders   []*rat.RAT     // Array of RAT (Remy Augmented TCP) senders
    Receivers []*Receiver    // Array of Receiver objects that receive packets
    Links     []*Link        // Array of Link objects representing network links
    Delay     time.Duration  // Propagation delay in the network
}

// NewNetwork is a constructor that creates a new instance of the Network struct
func NewNetwork(numSenders int, numLinks int, delay time.Duration) *Network {
    network := &Network{
        Senders:   make([]*rat.RAT, numSenders),
        Receivers: make([]*Receiver, numSenders),
        Links:     make([]*Link, numLinks),
        Delay:     delay,
    }

    // Initialize senders, receivers, and links
    for i := 0; i < numSenders; i++ {
        network.Senders[i] = rat.NewRAT(nil) // Initialize with nil whiskers
        network.Receivers[i] = NewReceiver()
    }
    for i := 0; i < numLinks; i++ {
        network.Links[i] = NewLink(1000, 10*time.Millisecond)
    }

    return network
}

// Run simulates the network for a given duration
func (n *Network) Run(duration time.Duration) {
    rand.Seed(time.Now().UnixNano()) // Seed the random number generator
    endTime := time.Now().Add(duration)

    for time.Now().Before(endTime) {
        for i, sender := range n.Senders {
            if sender.TimeToSend() {
                seq := rand.Intn(1000) // Generate a random sequence number
                packet := &Packet{
                    SeqNo: seq,
                    Sent:  time.Now(),
                }
                n.SendPacket(packet, sender, n.Receivers[i])
            }
        }
        time.Sleep(n.Delay) // Sleep for the propagation delay
    }
}

// SendPacket simulates the sending of a packet from a sender to a receiver
func (n *Network) SendPacket(packet *Packet, sender *rat.RAT, receiver *Receiver) {
    sender.OnPacketSent(packet.SeqNo)

    // Send the packet through all links
    for _, link := range n.Links {
        if !link.Enqueue(packet) {
            sender.OnPacketLost()
            return
        }
    }

    // Schedule packet delivery after the propagation delay
    go func() {
        time.Sleep(n.Delay)
        if rand.Float64() < 0.1 { // Simulate packet loss with 10% probability
            sender.OnPacketLost()
            return
        }
        receiver.ReceivePacket(packet)
        sender.OnPacketAcked(packet.SeqNo, time.Since(packet.Sent))
    }()
}

// Packet represents a network packet
type Packet struct {
    SeqNo int        // Sequence number
    Sent  time.Time  // Timestamp when the packet was sent
}

// Receiver represents a receiver that receives packets
type Receiver struct {
    ReceivedPackets []*Packet // Slice of received packets
}

// NewReceiver is a constructor that creates a new instance of the Receiver struct
func NewReceiver() *Receiver {
    return &Receiver{
        ReceivedPackets: make([]*Packet, 0),
    }
}

// ReceivePacket adds a received packet to the receiver's ReceivedPackets slice
func (r *Receiver) ReceivePacket(packet *Packet) {
    r.ReceivedPackets = append(r.ReceivedPackets, packet)
}
