package memory

// Define DataType as a type alias for float64
type DataType float64

// Define the Memory struct
type Memory struct {
    RecvRate         DataType // Receive rate
    SendRate         DataType // Send rate
    LatestDelay      DataType // Latest observed delay
    InterPacketDelay DataType // Inter-packet delay
}

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
