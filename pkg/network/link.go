package network

import (
    "time" // Import the time package for time-related operations
)

// Link represents a network link
type Link struct {
    Bandwidth    int           // Bandwidth of the link in bytes per second
    Latency      time.Duration // Latency of the link
    PacketQueue  chan *Packet  // Channel acting as a packet queue
}

// NewLink is a constructor that creates a new instance of the Link struct
func NewLink(bandwidth int, latency time.Duration) *Link {
    return &Link{
        Bandwidth:    bandwidth,
        Latency:      latency,
        PacketQueue:  make(chan *Packet, 1000), // Create a buffered channel with capacity 1000
    }
}

// Enqueue attempts to enqueue a packet into the link's packet queue
// It returns true if the packet was successfully enqueued, false otherwise
func (l *Link) Enqueue(packet *Packet) bool {
    select {
    case l.PacketQueue <- packet:
        return true
    default:
        return false
    }
}
