package main

import (
    "flag"  // Import the flag package for command-line argument parsing
    "fmt"   // Import the fmt package for formatted I/O
    "io/ioutil" // Import the ioutil package for file I/O utilities
    "os"    // Import the os package for operating system functionality

    "github.com/Aanthord/remy-go/pkg/whisker" // Import the whisker package from the remy project
)

// Define command-line flags
var (
    // configFile is a string flag for the path to the configuration file
    configFile = flag.String("config", "", "Path to the configuration file")
    // outputFile is a string flag for the path to save the generated whiskers
    outputFile = flag.String("output", "", "Path to save the generated whiskers")
)

func main() {
    // Parse the command-line flags
    flag.Parse()

    // Read the configuration file
    configData, err := ioutil.ReadFile(*configFile)
    if err != nil {
        fmt.Println("Error reading config file:", err)
        os.Exit(1)
    }

    // Parse the configuration data
    config, err := whisker.ParseConfig(configData)
    if err != nil {
        fmt.Println("Error parsing config:", err)
        os.Exit(1)
    }

    // Generate whiskers based on the configuration
    whiskers := whisker.GenerateWhiskers(config)

    // Save the generated whiskers to the specified output file
    err = whisker.SaveWhiskers(whiskers, *outputFile)
    if err != nil {
        fmt.Println("Error saving whiskers:", err)
        os.Exit(1)
    }

    // Print a success message
    fmt.Println("Whiskers generated and saved successfully.")
}
