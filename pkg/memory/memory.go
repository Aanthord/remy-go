package memory

import (
    "fmt"
    "math"
    "time"
    "unsafe"

    "github.com/Aanthord/remy-go/pkg/dna"
)

// Define DataType as a type alias for float64
type DataType float64

// Define the Memory struct
type Memory struct {
    RecvRate         DataType // Receive rate
    SendRate         DataType // Send rate
    LatestDelay      DataType // Latest observed delay
    InterPacketDelay DataType // Inter-packet delay
    lastSentTime     time.Time
    lastRecvTime     time.Time
    minRTT           DataType
}

const (
    alpha     = 1.0 / 8.0
    slowAlpha = 1.0 / 256.0
)

// Constructor function for creating a new Memory instance
func NewMemory() *Memory {
    return &Memory{}
}

// Method to update the fields of Memory
func (m *Memory) Update(recvRate, sendRate, latestDelay, interPacketDelay DataType) {
    m.RecvRate = recvRate
    m.SendRate = sendRate
    m.LatestDelay = latestDelay
    m.InterPacketDelay = interPacketDelay
}

// Getter method for the RecvRate field
func (m *Memory) GetRecvRate() DataType {
    return m.RecvRate
}

// Getter method for the SendRate field
func (m *Memory) GetSendRate() DataType {
    return m.SendRate
}

// Getter method for the LatestDelay field
func (m *Memory) GetLatestDelay() DataType {
    return m.LatestDelay
}

// Getter method for the InterPacketDelay field
func (m *Memory) GetInterPacketDelay() DataType {
    return m.InterPacketDelay
}

// MemoryRange represents a range of memory values
type MemoryRange struct {
    Lower *Memory
    Upper *Memory
}

// NewMemoryRange creates a new instance of MemoryRange
func NewMemoryRange(lower, upper *Memory) *MemoryRange {
    return &MemoryRange{
        Lower: lower,
        Upper: upper,
    }
}

// MinMemory returns the minimum possible memory values
func MinMemory() *Memory {
    return &Memory{
        RecvRate:         0,
        SendRate:         0,
        LatestDelay:      0,
        InterPacketDelay: 0,
    }
}

// MaxMemory returns the maximum possible memory values
func MaxMemory() *Memory {
    return &Memory{
        RecvRate:         DataType(math.MaxFloat64),
        SendRate:         DataType(math.MaxFloat64),
        LatestDelay:      DataType(math.MaxFloat64),
        InterPacketDelay: DataType(math.MaxFloat64),
    }
}

// Contains checks if a Memory value is within the MemoryRange
func (mr *MemoryRange) Contains(m *Memory) bool {
    return (m.RecvRate >= mr.Lower.RecvRate && m.RecvRate <= mr.Upper.RecvRate) &&
        (m.SendRate >= mr.Lower.SendRate && m.SendRate <= mr.Upper.SendRate) &&
        (m.LatestDelay >= mr.Lower.LatestDelay && m.LatestDelay <= mr.Upper.LatestDelay) &&
        (m.InterPacketDelay >= mr.Lower.InterPacketDelay && m.InterPacketDelay <= mr.Upper.InterPacketDelay)
}

// Intersects checks if two MemoryRange instances intersect
func (mr *MemoryRange) Intersects(other *MemoryRange) bool {
    return (mr.Lower.RecvRate <= other.Upper.RecvRate && mr.Upper.RecvRate >= other.Lower.RecvRate) &&
        (mr.Lower.SendRate <= other.Upper.SendRate && mr.Upper.SendRate >= other.Lower.SendRate) &&
        (mr.Lower.LatestDelay <= other.Upper.LatestDelay && mr.Upper.LatestDelay >= other.Lower.LatestDelay) &&
        (mr.Lower.InterPacketDelay <= other.Upper.InterPacketDelay && mr.Upper.InterPacketDelay >= other.Lower.InterPacketDelay)
}

// FromDNAMemory converts a dna.Memory to a memory.Memory.
func FromDNAMemory(dnaMem *dna.Memory) *Memory {
    return &Memory{
        RecvRate:         DataType(dnaMem.RecvRate),
        SendRate:         DataType(dnaMem.SendRate),
        LatestDelay:      DataType(dnaMem.LatestDelay),
        InterPacketDelay: DataType(dnaMem.InterPacketDelay),
    }
}

// ToDNAMemory converts a memory.Memory to a dna.Memory.
func (m *Memory) ToDNAMemory() *dna.Memory {
    return &dna.Memory{
        RecvRate:         float32(m.RecvRate),
        SendRate:         float32(m.SendRate),
        LatestDelay:      float32(m.LatestDelay),
        InterPacketDelay: float32(m.InterPacketDelay),
    }
}

// UpdateSentPacket updates the memory state with the information from the sent packet
func (m *Memory) UpdateSentPacket(packet *Packet) {
    if m.lastSentTime.IsZero() || m.lastRecvTime.IsZero() {
        m.lastSentTime = packet.Sent
        m.lastRecvTime = packet.Sent
        m.minRTT = 0
    } else {
        m.SendRate = (1 - alpha) * m.SendRate + alpha*(DataType(packet.Sent.Sub(m.lastSentTime).Seconds()))
        m.RecvRate = (1 - alpha) * m.RecvRate + alpha*(DataType(packet.Sent.Sub(m.lastRecvTime).Seconds()))
        m.minRTT = 0
        m.lastSentTime = packet.Sent
        m.lastRecvTime = packet.Sent
    }
}

// UpdateReceivedPackets updates the memory state with the received packets for the given flow ID
func (m *Memory) UpdateReceivedPackets(packets []*Packet, flowID uint) {
    for _, packet := range packets {
        if packet.FlowID != flowID {
            continue
        }

        rtt := DataType(packet.Received.Sub(packet.Sent).Seconds())
        if m.lastSentTime.IsZero() || m.lastRecvTime.IsZero() {
            m.lastSentTime = packet.Sent
            m.lastRecvTime = packet.Received
            m.minRTT = rtt
        } else {
            m.RecvRate = (1 - alpha) * m.RecvRate + alpha*(DataType(packet.Received.Sub(m.lastRecvTime).Seconds()))
            m.SendRate = (1 - alpha) * m.SendRate + alpha*(DataType(packet.Sent.Sub(m.lastSentTime).Seconds()))
            m.InterPacketDelay = (1 - slowAlpha) * m.InterPacketDelay + slowAlpha*(DataType(packet.Received.Sub(m.lastRecvTime).Seconds()))
            m.lastSentTime = packet.Sent
            m.lastRecvTime = packet.Received
            m.minRTT = DataType(math.Min(float64(m.minRTT), float64(rtt)))
            m.LatestDelay = rtt / m.minRTT
        }
    }
}

// UpdateRTT updates the memory state with the given round-trip time
func (m *Memory) UpdateRTT(rtt DataType) {
    if m.minRTT == 0 {
        m.minRTT = rtt
    } else {
        m.minRTT = DataType(math.Min(float64(m.minRTT), float64(rtt)))
    }
    m.LatestDelay = rtt / m.minRTT
}

// Packet represents a network packet
type Packet struct {
    SeqNo    int        // Sequence number
    ID       int        // Sender ID
    FlowID   uint       // Flow ID
    Sent     time.Time  // Timestamp when the packet was sent
    Received time.Time  // Timestamp when the packet was received
}

// HashCode returns a hash value for the memory state
func (m *Memory) HashCode() uint64 {
    hash := uint64(0)
    hash = hash*31 + *(*uint64)(unsafe.Pointer(&m.RecvRate))
    hash = hash*31 + *(*uint64)(unsafe.Pointer(&m.SendRate))
    hash = hash*31 + *(*uint64)(unsafe.Pointer(&m.LatestDelay))
    hash = hash*31 + *(*uint64)(unsafe.Pointer(&m.InterPacketDelay))
    return hash
}

// IsGreaterThanOrEqual checks if the current memory state is greater than or equal to another memory state
func (m *Memory) IsGreaterThanOrEqual(other *Memory) bool {
    return m.RecvRate >= other.RecvRate &&
        m.SendRate >= other.SendRate &&
        m.LatestDelay >= other.LatestDelay &&
        m.InterPacketDelay >= other.InterPacketDelay
}

// IsLessThan checks if the current memory state is less than another memory state
func (m *Memory) IsLessThan(other *Memory) bool {
    return m.RecvRate < other.RecvRate &&
        m.SendRate < other.SendRate &&
        m.LatestDelay < other.LatestDelay &&
        m.InterPacketDelay < other.InterPacketDelay
}

// IsEqual checks if the current memory state is equal to another memory state
func (m *Memory) IsEqual(other *Memory) bool {
    return m.RecvRate == other.RecvRate &&
        m.SendRate == other.SendRate &&
        m.LatestDelay == other.LatestDelay &&
        m.InterPacketDelay == other.InterPacketDelay
}

// String returns a string representation of the memory state
func (m *Memory) String() string {
    return fmt.Sprintf("RecvRate=%f, SendRate=%f, LatestDelay=%f, InterPacketDelay=%f",
        m.RecvRate, m.SendRate, m.LatestDelay, m.InterPacketDelay)
}

// Reset resets all the memory state fields to their initial values
func (m *Memory) Reset() {
    m.RecvRate = 0
    m.SendRate = 0
    m.LatestDelay = 0
    m.InterPacketDelay = 0
    m.lastSentTime = time.Time{}
    m.lastRecvTime = time.Time{}
    m.minRTT = 0
}

// AdvanceTo advances the memory state to the given tick (time)
func (m *Memory) AdvanceTo(tick uint64) {
    if !m.lastSentTime.IsZero() {
        m.SendRate = (1-alpha)*m.SendRate + alpha*DataType(float64(tick-uint64(m.lastSentTime.UnixNano()))/1e9)
    }
    if !m.lastRecvTime.IsZero() {
        m.RecvRate = (1-alpha)*m.RecvRate + alpha*DataType(float64(tick-uint64(m.lastRecvTime.UnixNano()))/1e9)
        m.InterPacketDelay = (1-slowAlpha)*m.InterPacketDelay + slowAlpha*DataType(float64(tick-uint64(m.lastRecvTime.UnixNano()))/1e9)
    }
    m.lastSentTime = time.Unix(0, int64(tick))
    m.lastRecvTime = time.Unix(0, int64(tick))
}
