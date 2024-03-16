package crat

/*
#cgo CPPFLAGS: -I../third_party/remy/include
#cgo LDFLAGS: -L../third_party/remy/lib -lremyprotos -lboost_system -lboost_random
#include "rat.hh"
#include "memory.hh"
#include "whiskertree.hh"
#include "packet.hh"
#include "dna.pb.h"
#include "receiver.hh"
#include "sendergang.hh"
#include "network.hh"
#include "configrange.hh"
#include "evaluator.hh"
*/
import "C"
import (
    "fmt"
    "unsafe"

    "github.com/Aanthord/remy-go/pkg/whisker"
)

// RAT is a Go wrapper around the C++ implementation of the RAT algorithm.
type RAT struct {
    cRAT *C.Rat
}

// NewRAT creates a new instance of the RAT struct.
func NewRAT(whiskers *whisker.WhiskerTree) *RAT {
    cWhiskerTree := C.WhiskerTree(whiskers.ToDNAWhiskerTree())
    cRAT := C.Rat(cWhiskerTree, false)
    return &RAT{cRAT: &cRAT}
}

// Evaluate evaluates the performance of the pre-trained model in a simulated network environment.
func (rat *RAT) Evaluate(config *Config) {
    cConfig := C.ConfigRange{
        C.pair_make_pair(C.double(config.LinkPPT), C.double(config.LinkPPT)),
        C.pair_make_pair(C.double(config.RTT), C.double(config.RTT)),
        C.pair_make_pair(C.int(config.NumSenders), C.int(config.NumSenders)),
        C.double(config.MeanOnDuration),
        C.double(config.MeanOffDuration),
        C.bool(false), // LOOnly is always false
    }
    cEvaluator := C.Evaluator(rat.cRAT, cConfig)
    cOutcome := cEvaluator.score(rat.cRAT, C.bool(false), C.uint(1))
    // Print the evaluation results
    fmt.Printf("Score: %f\n", cOutcome.score)
    for _, throughputDelay := range cOutcome.throughputs_delays {
        config := C.NetConfig{
            throughputDelay.first.mean_on_duration,
            throughputDelay.first.mean_off_duration,
            throughputDelay.first.num_senders,
            throughputDelay.first.link_ppt,
            throughputDelay.first.delay,
        }
        fmt.Printf("Config: %s\n", C.GoString(C.NetConfig_str(&config)))
        for _, result := range throughputDelay.second {
            fmt.Printf("Sender: [Throughput=%f, Delay=%f]\n", result.first/config.link_ppt, result.second/config.delay)
        }
    }
    fmt.Printf("Whiskers: %s\n", C.GoString(cOutcome.used_whiskers.str(C.uint(0))))
}

// Config represents the configuration for the network simulation and model evaluation.
type Config struct {
    LinkPPT        float64
    RTT            float64
    NumSenders     int
    MeanOnDuration float64
    MeanOffDuration float64
}
