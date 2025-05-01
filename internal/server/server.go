package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	ssh "github.com/charmbracelet/ssh"

	"github.com/cankurttekin/sh.kurttekin.com/internal/tui"
)

// Config holds the configuration for the SSH server
type Config struct {
	ListenAddr string
	KeyPath    string
}

// DefaultConfig returns the default server configuration
func DefaultConfig() Config {
	return Config{
		ListenAddr: ":2222",
	}
}

// Start initializes and starts the SSH server
func Start(config Config) error {
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
	// Check if we have a valid PTY
	pty, windowChange, isPty := s.Pty()
	if !isPty {
		fmt.Fprintln(s, "No active terminal, please run with ssh -t")
		return
	}

	// Clear the screen and ensure we're starting with a clean slate
	fmt.Fprint(s, "\033[2J\033[H\033[?25l") // Clear screen, move cursor to top-left, hide cursor

	// Log the connection for debugging
	log.Printf("New SSH connection from %s (window size: %d x %d)", s.RemoteAddr(), pty.Window.Width, pty.Window.Height)

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
			}
		}
	}()

	// Run the program
	if _, err := p.Run(); err != nil {
		// Show cursor again before displaying error
		fmt.Fprint(s, "\033[?25h")
		fmt.Fprintf(s, "Error running TUI: %v\n", err)
		log.Printf("TUI error for client %s: %v", s.RemoteAddr(), err)
	} else {
		// Always ensure cursor is visible when exiting
		fmt.Fprint(s, "\033[?25h")
	}
}
