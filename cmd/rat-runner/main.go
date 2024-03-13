package main

import (
    "flag"  // Import the flag package for command-line argument parsing
    "fmt"   // Import the fmt package for formatted I/O
    "os"    // Import the os package for operating system functionality
    "os/signal" // Import the signal package for signal handling
    "syscall"   // Import the syscall package for system calls

    "github.com/Aanthord/remy-go/pkg/rat"    // Import the rat package from the remy project
    "github.com/Aanthord/remy-go/pkg/whisker" // Import the whisker package from the remy project
)

// Define command-line flags
var (
    // whiskersFile is a string flag for the path to the whiskers file
    whiskersFile = flag.String("whiskers", "", "Path to the whiskers file")
)

func main() {
    // Parse the command-line flags
    flag.Parse()

    // Load whiskers from the specified file
    whiskers, err := whisker.LoadWhiskers(*whiskersFile)
    if err != nil {
        fmt.Println("Error loading whiskers:", err)
        os.Exit(1)
    }

    // Create a new RAT (Remy Augmented TCP) congestion controller with the loaded whiskers
    ratController := rat.NewRemyCongestionController(whiskers)

    // Install the RAT congestion controller
    err = ratController.Install()
    if err != nil {
        fmt.Println("Error installing RAT congestion controller:", err)
        os.Exit(1)
    }

    // Start the RAT congestion controller
    ratController.Start()

    // Wait for termination signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan

    // Stop the RAT congestion controller
    ratController.Stop()

    // Uninstall the RAT congestion controller
    err = ratController.Uninstall()
    if err != nil {
        fmt.Println("Error uninstalling RAT congestion controller:", err)
        os.Exit(1)
    }

    fmt.Println("RAT congestion controller stopped.")
}
