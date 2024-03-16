package rat

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "net"
    "sync"
    "time"

    "github.com/Aanthord/remy-go/pkg/memory"
    "github.com/Aanthord/remy-go/pkg/whisker"
    "golang.org/x/sys/unix"
)

// RAT represents the Remy Augmented TCP (RAT) congestion control algorithm
type RAT struct {
    whiskers        *whisker.WhiskerTree // Pointer to the WhiskerTree that holds the whiskers
    memory          *memory.Memory       // Pointer to the Memory object that stores the network state
    packetsSent     uint                 // Number of packets sent
    packetsReceived uint                 // Number of packets received
    track           bool                 // Flag to enable tracking (used in the C++ implementation)
    lastSendTime    time.Time            // Timestamp of the last packet sent
    congestionWindow uint                // Current congestion window size
    intersendTime   float64              // Intersend time for the current whisker
    flowID          uint                 // Flow ID
    currentWhisker  *whisker.Whisker     // Pointer to the currently selected Whisker
    mu              sync.Mutex           // Mutex for synchronization
}

// NewRAT is a constructor that creates a new instance of the RAT struct
func NewRAT(whiskers *whisker.WhiskerTree, track bool) *RAT {
    return &RAT{
        whiskers:       whiskers,
        memory:         memory.NewMemory(),
        track:          track,
        flowID:         0,
        currentWhisker: whiskers.Root.Whisker,
    }
}

// Start initializes the RAT algorithm
// It creates a TCP listener and sets the RAT algorithm as the system-wide congestion control algorithm
func (rat *RAT) Start() error {
    // Create a TCP listener
    listener, err := net.Listen("tcp", ":0")
    if err != nil {
        return err
    }
    defer listener.Close()

    // Get the listener's socket file descriptor
    listenerFD, err := listener.(*net.TCPListener).File()
    if err != nil {
        return err
    }
    defer listenerFD.Close()

    // Set the RAT algorithm as the system-wide congestion control algorithm
    err = SetRATAlgorithmSystemWide(int(listenerFD.Fd()))
    if err != nil {
        return err
    }

    return nil
}

// Send sends a packet from the sender to the receiver
func (rat *RAT) Send(id int, next *net.TCPConn, seq int, packetsSentCap uint) error {
    rat.mu.Lock()
    defer rat.mu.Unlock()

    // Assertion to ensure that the number of packets sent is greater than or equal to the number of packets received
    assertCondition(rat.packetsSent >= rat.packetsReceived, "Number of packets sent should be greater than or equal to the number of packets received")

    if rat.congestionWindow == 0 {
        // If the congestion window is zero, initialize the current whisker, congestion window, and intersend time
        rat.currentWhisker = rat.whiskers.Root.Whisker
        rat.congestionWindow = rat.currentWhisker.Window(0)
        rat.intersendTime = rat.currentWhisker.Intersend
    }

    if rat.packetsSent < rat.packetsReceived+rat.congestionWindow &&
        time.Since(rat.lastSendTime) >= time.Duration(rat.intersendTime*float64(time.Second)) {

        // Check if it's time to send a packet based on the congestion window and intersend time

        // Have we reached the end of the flow for now?
        if rat.packetsSent >= packetsSentCap {
            return nil
        }

        // Create a new Packet struct and populate its fields
        packet := &Packet{
            SeqNo:  seq,
            ID:     id,
            FlowID: rat.flowID,
            Sent:   time.Now(),
        }
        rat.packetsSent++
        rat.memory.UpdateSentPacket(&memory.Packet{
            SeqNo:    packet.SeqNo,
            ID:       packet.ID,
            FlowID:   packet.FlowID,
            Sent:     packet.Sent,
            Received: packet.Received,
        })
        err := SendPacket(next, packet)
        if err != nil {
            return err
        }
        rat.lastSendTime = time.Now()
    }

    return nil
}

// ReceivePackets receives a slice of packets and updates the RAT state
func (rat *RAT) ReceivePackets(packets []*Packet) {
    rat.mu.Lock()
    defer rat.mu.Unlock()

    rat.packetsReceived += uint(len(packets))

    var memoryPackets []*memory.Packet
    for _, packet := range packets {
        if packet.FlowID == rat.flowID {
            rtt := time.Since(packet.Sent)
            memoryPackets = append(memoryPackets, &memory.Packet{
                SeqNo:    packet.SeqNo,
                ID:       packet.ID,
                FlowID:   packet.FlowID,
                Sent:     packet.Sent,
                Received: packet.Received,
            })
            whisker, err := rat.whiskers.FindWhisker(rat.memory)
            if err != nil {
                fmt.Println("Error finding whisker:", err)
                continue
            }
            rat.updateState(rtt.Seconds(), packet.SeqNo, whisker)
        }
    }
    rat.memory.UpdateReceivedPackets(memoryPackets, rat.flowID)
}

// updateState updates the RAT state based on the received packet
func (rat *RAT) updateState(rtt float64, seq int, currentWhisker *whisker.Whisker) {
    rat.memory.UpdateRTT(memory.DataType(rtt))
    rat.currentWhisker = currentWhisker
    rat.congestionWindow = rat.currentWhisker.Window(rat.congestionWindow)
    rat.intersendTime = rat.currentWhisker.Intersend
    rat.flowID++
    if rat.flowID == 0 {
        rat.flowID = 1
    }
}

// NextEventTime returns the time of the next event
func (rat *RAT) NextEventTime() time.Time {
    nextSendTime := rat.lastSendTime.Add(time.Duration(rat.intersendTime * float64(time.Second)))
    return nextSendTime
}

// PacketsSent returns the number of packets sent
func (rat *RAT) PacketsSent() uint {
    return rat.packetsSent
}

// Whiskers returns the WhiskerTree
func (rat *RAT) Whiskers() *whisker.WhiskerTree {
    return rat.whiskers
}

// Packet represents a network packet
type Packet struct {
    SeqNo    int        // Sequence number
    ID       int        // Sender ID
    FlowID   uint       // Flow ID
    Sent     time.Time  // Timestamp when the packet was sent
    Received time.Time  // Timestamp when the packet was received
}

// SendPacket sends a packet over the network
func SendPacket(conn *net.TCPConn, packet *Packet) error {
    data, err := encodePacket(packet)
    if err != nil {
        return err
    }
    _, err = conn.Write(data)
    return err
}

// ReceivePacket receives a packet from the network
func ReceivePacket(conn *net.TCPConn) (*Packet, error) {
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        return nil, err
    }
    packet, err := decodePacket(buffer[:n])
    if err != nil {
        return nil, err
    }
    packet.Received = time.Now()
    return packet, nil
}

// encodePacket encodes a packet into a byte slice
func encodePacket(packet *Packet) ([]byte, error) {
    buf := new(bytes.Buffer)

    err := binary.Write(buf, binary.BigEndian, int32(packet.SeqNo))
    if err != nil {
        return nil, err
    }

    err = binary.Write(buf, binary.BigEndian, int32(packet.ID))
    if err != nil {
        return nil, err
    }

    err = binary.Write(buf, binary.BigEndian, packet.FlowID)
    if err != nil {
        return nil, err
    }

    err = binary.Write(buf, binary.BigEndian, packet.Sent.UnixNano())
    if err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

// decodePacket decodes a packet from a byte slice
func decodePacket(buffer []byte) (*Packet, error) {
    buf := bytes.NewReader(buffer)

    packet := &Packet{}

    var seqNo int32
    err := binary.Read(buf, binary.BigEndian, &seqNo)
    if err != nil {
        return nil, err
    }
    packet.SeqNo = int(seqNo)

    var id int32
    err = binary.Read(buf, binary.BigEndian, &id)
    if err != nil {
        return nil, err
    }
    packet.ID = int(id)

    err = binary.Read(buf, binary.BigEndian, &packet.FlowID)
    if err != nil {
        return nil, err
    }

    var sent int64
    err = binary.Read(buf, binary.BigEndian, &sent)
    if err != nil {
        return nil, err
    }
    packet.Sent = time.Unix(0, sent)

    return packet, nil
}

// assertCondition is a function for assertions
func assertCondition(condition bool, message string) {
    if !condition {
        panic(message)
    }
}

// SetRATAlgorithmSystemWide sets the RAT congestion control algorithm as the system-wide default
func SetRATAlgorithmSystemWide(sockFD int) error {
    err := unix.SetsockoptString(sockFD, unix.SOL_TCP, unix.TCP_CONGESTION, "remy")
    if err != nil {
        return fmt.Errorf("failed to set RAT congestion control algorithm: %v", err)
    }
    return nil
}
