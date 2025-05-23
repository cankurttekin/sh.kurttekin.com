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

// ssh server configuration
type Config struct {
	ListenAddr string
	KeyPath    string
	LogFile    string
}

func DefaultConfig() Config {
	defaultLogPath := "tuiserver_connections.log"

	return Config{
		ListenAddr: ":2222",
		LogFile:    defaultLogPath,
	}
}

func setupLogger(logFilePath string) (*os.File, error) {
	if logFilePath == "" {
		return nil, nil
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(multiWriter)

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	return logFile, nil
}

func Start(config Config) error {
	logFile, err := setupLogger(config.LogFile)
	if err != nil {
		log.Printf("Warning: Could not set up file logging to %s: %v. Continuing with stdout only.", config.LogFile, err)
	}

	// close the log file when the server stops
	if logFile != nil {
		defer logFile.Close()
		log.Printf("Logging connections to file: %s", config.LogFile)
	} else if config.LogFile != "" {
		log.Printf("Failed to open log file: %s - logs will only be shown in stdout", config.LogFile)
	}

	// check if the port is already in use
	if isPortInUse(config.ListenAddr) {
		log.Printf("Port %s is already in use. Try stopping existing SSH server or using a different port.", config.ListenAddr)
		return fmt.Errorf("port %s is already in use", config.ListenAddr)
	}

	server := ssh.Server{
		Addr:    config.ListenAddr,
		Handler: handleSession,
	}

	log.Printf("SSH server started on %s\n", config.ListenAddr)
	return server.ListenAndServe()
}

func isPortInUse(addr string) bool {
	// extract port from listen address
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		if strings.HasPrefix(addr, ":") {
			port = addr[1:]
			host = "localhost"
		} else {
			return false
		}
	}
	if host == "" {
		host = "localhost"
	}

	ln, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return true
	}

	// port is available, close the listener
	ln.Close()
	return false
}

// handleSession is called when a new SSH session is established
func handleSession(s ssh.Session) {
	remoteAddr := s.RemoteAddr().String()
	sessionID := s.Context().SessionID()
	username := s.User()

	startTime := time.Now()
	log.Printf("+ Connection opened | User: %s | IP: %s | Session: %s | Time: %s",
		username, remoteAddr, sessionID, startTime.Format(time.RFC3339))

	// check if we have a valid PTY
	pty, windowChange, isPty := s.Pty()
	if !isPty {
		fmt.Fprintln(s, "No active terminal, please run with ssh -t")
		log.Printf("Connection closed (no PTY) | User: %s | IP: %s | Session: %s | Duration: %s",
			username, remoteAddr, sessionID, time.Since(startTime))
		return
	}

	// logging terminal details 
	log.Printf("Terminal info | Term: %s | Size: %dx%d | Session: %s",
		pty.Term, pty.Window.Width, pty.Window.Height, sessionID)

	// clear the screen and hide the cursor
	fmt.Fprint(s, "\033[2J\033[H\033[?25l") 

	// initialize model with term dimensions
	m := tui.NewModel(pty.Window.Width, pty.Window.Height)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithInput(s),          // Use SSH session for input
		tea.WithOutput(s),         // Use SSH session for output
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// handle window resizing events
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

	if _, err := p.Run(); err != nil {
		// show cursor again before displaying error
		fmt.Fprint(s, "\033[?25h")
		fmt.Fprintf(s, "Error running TUI: %v\n", err)
		log.Printf("Error | Session: %s | User: %s | Error: %v",
			sessionID, username, err)
	} else {
		// show cursor again before exiting
		fmt.Fprint(s, "\033[?25h")
	}

	// logging connection termination
	duration := time.Since(startTime)
	log.Printf("- Connection closed | Session: %s | User: %s | IP: %s | Duration: %s",
		sessionID, username, remoteAddr, duration)
}
