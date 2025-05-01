package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cankurttekin/sh.kurttekin.com/internal/server"
)

func main() {
	// Get configuration from command line flags
	var config server.Config

	// Define command line flags
	flag.StringVar(&config.ListenAddr, "addr", ":2222", "SSH server address")
	flag.StringVar(&config.KeyPath, "key", "", "Path to SSH server key (optional)")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// Parse flags
	flag.Parse()

	// Start the SSH server
	if err := server.Start(config); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
