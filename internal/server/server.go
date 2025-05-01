package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	ssh "github.com/charmbracelet/ssh"

	"github.com/cankurttekin/sh.kurttekin.com/internal/tui"
)

// Config holds the configuration for the SSH server
type Config struct {
	ListenAddr string
	KeyPath    string
	LogFile    string
}

// DefaultConfig returns the default server configuration
func DefaultConfig() Config {
	// Set default log file in the current directory
	defaultLogPath := "tuiserver_connections.log"

	return Config{
		ListenAddr: ":2222",
		LogFile:    defaultLogPath,
	}
}

// setupLogger configures logging to both stdout and file if specified
func setupLogger(logFilePath string) (*os.File, error) {
	// If no log file is specified, use default logger (stdout)
	if logFilePath == "" {
		return nil, nil
	}

	// Open log file (create if doesn't exist, append if it does)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer to log to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Configure logger to use our multi-writer
	log.SetOutput(multiWriter)

	// Add timestamps to log entries
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	return logFile, nil
}

// Start initializes and starts the SSH server
func Start(config Config) error {
	// Set up the logger if a log file is specified
	logFile, err := setupLogger(config.LogFile)
	if err != nil {
		log.Printf("Warning: Could not set up file logging to %s: %v. Continuing with stdout only.", config.LogFile, err)
	}

	// Close the log file when the server stops
	if logFile != nil {
		defer logFile.Close()
		log.Printf("Logging connections to file: %s", config.LogFile)
	} else if config.LogFile != "" {
		log.Printf("Failed to open log file: %s - logs will only be shown in stdout", config.LogFile)
	}

	// Check if the port is already in use
	if isPortInUse(config.ListenAddr) {
		log.Printf("Port %s is already in use. Try stopping existing SSH server or using a different port.", config.ListenAddr)
		return fmt.Errorf("port %s is already in use", config.ListenAddr)
	}

	server := ssh.Server{
		Addr:    config.ListenAddr,
		Handler: handleSession,
	}

	log.Printf("🔒 TUI SSH server started on %s\n", config.ListenAddr)
	return server.ListenAndServe()
}

// isPortInUse checks if a port is already in use
func isPortInUse(addr string) bool {
	// Extract port from listen address
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		// If the address format is invalid, default to checking the whole string
		if strings.HasPrefix(addr, ":") {
			port = addr[1:]
			host = "localhost"
		} else {
			// Invalid address format, assume port is free
			return false
		}
	}
	if host == "" {
		host = "localhost"
	}

	// Try to open the port
	ln, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		// If we can't listen, the port is in use
		return true
	}

	// Port is available, close the listener
	ln.Close()
	return false
}

// handleSession is called when a new SSH session is established
func handleSession(s ssh.Session) {
	// Get client info for logging
	remoteAddr := s.RemoteAddr().String()
	sessionID := s.Context().SessionID()
	username := s.User()

	// Log connection start with more details
	startTime := time.Now()
	log.Printf("Connection opened | Session: %s | User: %s | IP: %s | Time: %s",
		sessionID, username, remoteAddr, startTime.Format(time.RFC3339))

	// Check if we have a valid PTY
	pty, windowChange, isPty := s.Pty()
	if !isPty {
		fmt.Fprintln(s, "No active terminal, please run with ssh -t")
		log.Printf("Connection closed (no PTY) | Session: %s | User: %s | IP: %s | Duration: %s",
			sessionID, username, remoteAddr, time.Since(startTime))
		return
	}

	// Log terminal details
	log.Printf("Terminal info | Session: %s | Term: %s | Size: %dx%d",
		sessionID, pty.Term, pty.Window.Width, pty.Window.Height)

	// Clear the screen and ensure we're starting with a clean slate
	fmt.Fprint(s, "\033[2J\033[H\033[?25l") // Clear screen, move cursor to top-left, hide cursor

	// Initialize model with terminal dimensions
	m := tui.NewModel(pty.Window.Width, pty.Window.Height)

	// Configure Bubble Tea program with full options
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithInput(s),          // Use SSH session for input
		tea.WithOutput(s),         // Use SSH session for output
		tea.WithMouseCellMotion(), // Enable mouse support for better interaction
	)

	// Handle window resizing events
	go func() {
		for {
			select {
			case <-s.Context().Done():
				return
			case w := <-windowChange:
				p.Send(tea.WindowSizeMsg{
					Width:  w.Width,
					Height: w.Height,
				})
				log.Printf("Terminal resize | Session: %s | New size: %dx%d",
					sessionID, w.Width, w.Height)
			}
		}
	}()

	// Run the program
	if _, err := p.Run(); err != nil {
		// Show cursor again before displaying error
		fmt.Fprint(s, "\033[?25h")
		fmt.Fprintf(s, "Error running TUI: %v\n", err)
		log.Printf("Error | Session: %s | User: %s | Error: %v",
			sessionID, username, err)
	} else {
		// Always ensure cursor is visible when exiting
		fmt.Fprint(s, "\033[?25h")
	}

	// Log connection termination
	duration := time.Since(startTime)
	log.Printf("Connection closed | Session: %s | User: %s | IP: %s | Duration: %s",
		sessionID, username, remoteAddr, duration)
}
