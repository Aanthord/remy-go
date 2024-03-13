package network

import (
    "time" // Import the time package for time-related operations
)

// Delay represents the propagation delay in the network
type Delay struct {
    Delay time.Duration // Duration of the delay
}

// NewDelay is a constructor that creates a new instance of the Delay struct
func NewDelay(delay time.Duration) *Delay {
    return &Delay{
        Delay: delay,
    }
}
