package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cankurttekin/sh.kurttekin.com/internal/server"
)

func main() {
	// Start with default configuration
	config := server.DefaultConfig()

	// Override defaults with command line flags
	flag.StringVar(&config.ListenAddr, "addr", config.ListenAddr, "SSH server address")
	flag.StringVar(&config.KeyPath, "key", config.KeyPath, "Path to SSH server key (optional)")

	// Create a temporary variable for the log flag to detect if it was explicitly set
	var logFilePath string
	flag.StringVar(&logFilePath, "log", "", "Path to connection log file (optional, default: tuiserver_connections.log)")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// Parse flags
	flag.Parse()

	// If log file was explicitly set, override the default
	if logFilePath != "" {
		config.LogFile = logFilePath
	}

	// Handle relative paths for the log file
	if config.LogFile != "" && !filepath.IsAbs(config.LogFile) {
		// If it's a simple filename without directory separators, place it in the executable directory
		if !strings.Contains(config.LogFile, "/") && !strings.Contains(config.LogFile, "\\") {
			execPath, err := os.Executable()
			if err == nil {
				config.LogFile = filepath.Join(filepath.Dir(execPath), config.LogFile)
			}
		} else {
			// If it contains path separators but is still relative, make it relative to current working directory
			absPath, err := filepath.Abs(config.LogFile)
			if err == nil {
				config.LogFile = absPath
			}
		}
	}

	// Start the SSH server
	if err := server.Start(config); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
