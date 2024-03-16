package main

import (
    "flag"
    "fmt"
    "net"
    "strconv"
    "sync"

    "github.com/Aanthord/remy-go/pkg/rat"
    "github.com/Aanthord/remy-go/pkg/whisker"
)

var (
    whiskersFile   = flag.String("if", "", "Path to the file containing the pre-trained WhiskerTree")
    linkPPTString  = flag.String("link", "1.0", "Link packets per millisecond")
    rttString      = flag.String("rtt", "150.0", "Round-trip time in milliseconds")
    numSendersInt  = flag.Int("nsrc", 8, "Maximum number of senders")
    meanOnDuration = flag.Float64("on", 5000.0, "Mean on duration in milliseconds")
    meanOffDuration = flag.Float64("off", 5000.0, "Mean off duration in milliseconds")
)

func main() {
    flag.Parse()

    // Load whiskers from the specified file
    whiskerTree, err := whisker.LoadWhiskers(*whiskersFile)
    if err != nil {
        fmt.Printf("Error loading whiskers: %v\n", err)
        return
    }

    // Parse the link packets per millisecond
    linkPPT, err := strconv.ParseFloat(*linkPPTString, 64)
    if err != nil {
        fmt.Printf("Error parsing link packets per millisecond: %v\n", err)
        return
    }

    // Parse the round-trip time
    rtt, err := strconv.ParseFloat(*rttString, 64)
    if err != nil {
        fmt.Printf("Error parsing round-trip time: %v\n", err)
        return
    }

    // Create a new RAT congestion controller with the loaded whiskers
    ratController := rat.NewRAT(whiskerTree, false)

    // Start the RAT congestion controller
    if err := ratController.Start(); err != nil {
        fmt.Printf("Error starting RAT congestion controller: %v\n", err)
        return
    }

    // Create a TCP connection for sending and receiving packets
    conn, err := net.Dial("tcp", ":0")
    if err != nil {
        fmt.Printf("Error creating TCP connection: %v\n", err)
        return
    }
    defer conn.Close()

    // Create a WaitGroup to wait for all senders to finish
    var wg sync.WaitGroup

    // Create senders
    for i := 0; i < *numSendersInt; i++ {
        wg.Add(1)
        go runSender(i, ratController, conn, &wg, float64(linkPPT), float64(rtt), *meanOnDuration, *meanOffDuration)
    }

    // Create a receiver
    go runReceiver(ratController, conn)

    // Wait for all senders to finish
    wg.Wait()
}

// runSender simulates a sender
func runSender(id int, ratController *rat.RAT, conn net.Conn, wg *sync.WaitGroup, linkPPT, rtt, meanOnDuration, meanOffDuration float64) {
    defer wg.Done()

    // Initialize sender state
    numPacketsSent := 0
    numPacketsReceived := 0
    onDuration := 0.0
    offDuration := 0.0
    isSending := true

    for {
        // Check if it's time to send a packet
        if isSending {
            // Create a packet and send it
            packet := &rat.Packet{
                SeqNo: numPacketsSent, // Set the sequence number
                ID:    id,             // Set the sender ID
            }
            if err := ratController.Send(packet.ID, conn.(*net.TCPConn), packet.SeqNo, uint(numPacketsSent+1)); err != nil {
                fmt.Printf("Error sending packet: %v\n", err)
                return
            }
            numPacketsSent++
        }

        // Receive packets
        packet, err := rat.ReceivePacket(conn.(*net.TCPConn))
        if err != nil {
            fmt.Printf("Error receiving packet: %v\n", err)
            return
        }

        // Update the RAT state with the received packet
        ratController.ReceivePackets([]*rat.Packet{packet})
        numPacketsReceived++

        // Update the sender state
        if isSending {
            onDuration += rtt
            if onDuration >= meanOnDuration {
                isSending = false
                onDuration = 0.0
            }
        } else {
            offDuration += rtt
            if offDuration >= meanOffDuration {
                isSending = true
                offDuration = 0.0
            }
        }
    }
}

// runReceiver simulates a receiver
func runReceiver(ratController *rat.RAT, conn net.Conn) {
    for {
        // Receive packets
        packet, err := rat.ReceivePacket(conn.(*net.TCPConn))
        if err != nil {
            fmt.Printf("Error receiving packet: %v\n", err)
            return
        }

        // Update the RAT state with the received packet
        ratController.ReceivePackets([]*rat.Packet{packet})
    }
}
