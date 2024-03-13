package dna

import (
    "fmt" // Import the fmt package for formatted I/O
    "github.com/Aanthord/remy-go/proto" // Import the proto package for Protocol Buffers support
)

// ConfigRange represents the configuration range for the network parameters
type ConfigRange struct {
    LinkPpt           *Range `protobuf:"bytes,1,opt,name=link_ppt,json=linkPpt" json:"link_ppt,omitempty"` // Link packets per second (ppt)
    Rtt               *Range `protobuf:"bytes,2,opt,name=rtt" json:"rtt,omitempty"` // Round-trip time (RTT)
    NumSenders        *Range `protobuf:"bytes,3,opt,name=num_senders,json=numSenders" json:"num_senders,omitempty"` // Number of senders
    MeanOfferDuration *Range `protobuf:"bytes,4,opt,name=mean_offer_duration,json=meanOfferDuration" json:"mean_offer_duration,omitempty"` // Mean offer duration
    XXX_unrecognized  []byte `json:"-"` // Unrecognized fields
}

// Reset resets the ConfigRange struct to its default values
func (m *ConfigRange) Reset()         { *m = ConfigRange{} }

// String returns a string representation of the ConfigRange struct
func (m *ConfigRange) String() string { return proto.CompactTextString(m) }

// ProtoMessage is an empty method to satisfy the proto.Message interface
func (*ConfigRange) ProtoMessage()    {}

// Range represents a range of values with a lower and upper bound
type Range struct {
    Low              *float64 `protobuf:"fixed64,1,opt,name=low" json:"low,omitempty"` // Lower bound
    High             *float64 `protobuf:"fixed64,2,opt,name=high" json:"high,omitempty"` // Upper bound
    XXX_unrecognized []byte   `json:"-"` // Unrecognized fields
}

// Reset resets the Range struct to its default values
func (m *Range) Reset()         { *m = Range{} }

// String returns a string representation of the Range struct
func (m *Range) String() string { return proto.CompactTextString(m) }

// ProtoMessage is an empty method to satisfy the proto.Message interface
func (*Range) ProtoMessage()    {}

// Whisker represents a single whisker in the Remy algorithm
type Whisker struct {
    Generation        *int32        `protobuf:"varint,1,opt,name=generation" json:"generation,omitempty"` // Generation of the whisker
    WindowIncrement   *int32        `protobuf:"varint,2,opt,name=window_increment,json=windowIncrement" json:"window_increment,omitempty"` // Window increment
    WindowMultiple    *float64      `protobuf:"fixed64,3,opt,name=window_multiple,json=windowMultiple" json:"window_multiple,omitempty"` // Window multiple
    Intersend         *float64      `protobuf:"fixed64,4,opt,name=intersend" json:"intersend,omitempty"` // Intersend time
    Domain            *MemoryRange  `protobuf:"bytes,5,opt,name=domain" json:"domain,omitempty"` // Memory domain
    XXX_unrecognized  []byte        `json:"-"` // Unrecognized fields
}

// Reset resets the Whisker struct to its default values
func (m *Whisker) Reset()         { *m = Whisker{} }

// String returns a string representation of the Whisker struct
func (m *Whisker) String() string { return proto.CompactTextString(m) }

// ProtoMessage is an empty method to satisfy the proto.Message interface
func (*Whisker) ProtoMessage()    {}

// MemoryRange represents a range of memory values with a lower and upper bound
type MemoryRange struct {
    Lower            *Memory `protobuf:"bytes,1,opt,name=lower" json:"lower,omitempty"` // Lower bound
    Upper            *Memory `protobuf:"bytes,2,opt,name=upper" json:"upper,omitempty"` // Upper bound
    XXX_unrecognized []byte  `json:"-"` // Unrecognized fields
}

// Reset resets the MemoryRange struct to its default values
func (m *MemoryRange) Reset()         { *m = MemoryRange{} }

// String returns a string representation of the MemoryRange struct
func (m *MemoryRange) String() string { return proto.CompactTextString(m) }

// ProtoMessage is an empty method to satisfy the proto.Message interface
func (*MemoryRange) ProtoMessage()    {}

// Memory represents the memory state used by the Remy algorithm
type Memory struct {
    RecvSendEwma      *float64 `protobuf:"fixed64,1,opt,name=recv_send_ewma,json=recvSendEwma" json:"recv_send_ewma,omitempty"` // Receive-send EWMA
    RecvRecEwma       *float64 `protobuf:"fixed64,2,opt,name=recv_rec_ewma,json=recvRecEwma" json:"recv_rec_ewma,omitempty"` // Receive-receive EWMA
    RttRatio          *float64 `protobuf:"fixed64,3,opt,name=rtt_ratio,json=rttRatio" json:"rtt_ratio,omitempty"` // RTT ratio
    SlowRecEwma       *float64 `protobuf:"fixed64,4,opt,name=slow_rec_ewma,json=slowRecEwma" json:"slow_rec_ewma,omitempty"` // Slow receive EWMA
    XXX_unrecognized  []byte   `json:"-"` // Unrecognized fields
}

// Reset resets the Memory struct to its default values
func (m *Memory) Reset()         { *m = Memory{} }

// String returns a string representation of the Memory struct
func (m *Memory) String() string { return proto.CompactTextString(m) }

// ProtoMessage is an empty method to satisfy the proto.Message interface
func (*Memory) ProtoMessage()    {}
