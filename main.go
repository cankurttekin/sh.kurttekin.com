package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ssh "github.com/charmbracelet/ssh"
)

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).PaddingBottom(1)
	sectionStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	inactiveStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	contentStyle      = lipgloss.NewStyle().PaddingLeft(4)
	linkStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	selectedLinkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Underline(true)
)

type model struct {
	sectionCursor int
	linkCursor    int
	inLinkMode    bool
	links         []string  // Links in the current section
	width         int
	height        int
}

type section struct {
	title   string
	content []string
}

var sections = []section{
	{"about", []string{
		"i am a software engineer and full-time observer and tinkerer.",
		"i love all kinds of engineering and development. i love free software, freedom in general.",
	}},
	{"experience", []string{
		"💼 software engineer @ akgun technology (2025 - Present)",
		"🧪 software developer intern @ comp. (2020 - 2022)",
		"🧑‍🎓 software developer intern @ comp. (2020 - 2021)",
		"📚 software engineering student @ canakkale onsekiz mart university -- turkey (2017 - 2023)",
	}},
	{"projects", []string{
		"🔧 ssh tui portfolio",
		"",
		"",
		"",
	}},
	{"links", []string{
		"github: https://github.com/cankurttekin",
		"linkedin: https://linkedin.com/in/cankurttekin",
	}},
}

// Extract links from content
func findLinks(content []string) []string {
	links := []string{}
	re := regexp.MustCompile(`https?://\S+`)
	
	for _, line := range content {
		matches := re.FindAllString(line, -1)
		links = append(links, matches...)
	}
	
	return links
}

func openURL(url string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	
	return cmd.Start()
}

// Message when a URL should be opened
type openURLMsg string

func openURLCommand(url string) tea.Cmd {
	return func() tea.Msg {
		err := openURL(url)
		if err != nil {
			return nil
		}
		return openURLMsg(url)
	}
}

func (m model) Init() tea.Cmd {
	// Get links for current section
	if m.sectionCursor < len(sections) {
		m.links = findLinks(sections[m.sectionCursor].content)
	}
	
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			// Only toggle link mode if current section has links
			currentSectionLinks := findLinks(sections[m.sectionCursor].content)
			if len(currentSectionLinks) > 0 {
				m.inLinkMode = !m.inLinkMode
				m.links = currentSectionLinks
				
				// Reset link cursor when entering link mode
				if m.inLinkMode && m.linkCursor >= len(m.links) {
					m.linkCursor = 0
				}
			}
		case "j", "down":
			if m.inLinkMode {
				// Navigate links in current section
				if m.linkCursor < len(m.links)-1 {
					m.linkCursor++
				}
			} else {
				// Navigate sections
				if m.sectionCursor < len(sections)-1 {
					m.sectionCursor++
					// Update links for the new section
					m.links = findLinks(sections[m.sectionCursor].content)
					m.inLinkMode = false
				}
			}
		case "k", "up":
			if m.inLinkMode {
				// Navigate links in current section
				if m.linkCursor > 0 {
					m.linkCursor--
				}
			} else {
				// Navigate sections
				if m.sectionCursor > 0 {
					m.sectionCursor--
					// Update links for the new section
					m.links = findLinks(sections[m.sectionCursor].content)
					m.inLinkMode = false
				}
			}
		case "enter":
			if m.inLinkMode && m.linkCursor < len(m.links) {
				// Open the selected link in a browser
				return m, openURLCommand(m.links[m.linkCursor])
			}
		}
	case tea.WindowSizeMsg:
		// Update the model with the new window size
		m.width = msg.Width
		m.height = msg.Height
	case openURLMsg:
	}
	return m, nil
}

func (m model) View() string {
	// Container for all content
	doc := strings.Builder{}
	
	// Title
	doc.WriteString(titleStyle.Render("cankurttekin") + "\n\n")
	
	// Render each section with padding and styling
	for i, sec := range sections {
		// Determine cursor and style based on selection
		cursor := "  "
		style := inactiveStyle
		
		// Section is focused if it's selected and we're not in link mode
		// OR if it's selected and we are in link mode (both are the same section)
		if i == m.sectionCursor && !m.inLinkMode {
			cursor = "➜ "
			style = focusedStyle
		}
		
		// Render section header
		sectionHeader := style.Render(cursor + sec.title)
		doc.WriteString(sectionHeader + "\n")
		
		// Always render content for all sections
		content := strings.Builder{}
		
		// Only highlight links in the current section when in link mode
		shouldHighlightLinks := i == m.sectionCursor && m.inLinkMode
		
		for _, line := range sec.content {
			// Check if this line contains links that need highlighting
			if shouldHighlightLinks {
				re := regexp.MustCompile(`(https?://\S+)`)
				formattedLine := line
				
				// Find all links in this line
				matches := re.FindAllStringIndex(line, -1)
				
				// Process the line from right to left to avoid index issues
				for j := len(matches) - 1; j >= 0; j-- {
					match := matches[j]
					linkText := line[match[0]:match[1]]
					
					// Determine if this link is selected
					isSelected := false
					for linkIdx, l := range m.links {
						if l == linkText && linkIdx == m.linkCursor {
							isSelected = true
							break
						}
					}
					
					// Apply appropriate style
					styledLink := ""
					if isSelected {
						styledLink = selectedLinkStyle.Render(linkText)
					} else {
						styledLink = linkStyle.Render(linkText)
					}
					
					// Replace the link in the line
					formattedLine = formattedLine[:match[0]] + styledLink + formattedLine[match[1]:]
				}
				
				content.WriteString(formattedLine + "\n")
			} else {
				// If it's the current section but not in link mode,
				// still show links but in normal link style
				if i == m.sectionCursor && !m.inLinkMode {
					re := regexp.MustCompile(`(https?://\S+)`)
					formattedLine := line
					
					// Find all links in this line
					matches := re.FindAllStringIndex(line, -1)
					
					// Process the line from right to left to avoid index issues
					for j := len(matches) - 1; j >= 0; j-- {
						match := matches[j]
						linkText := line[match[0]:match[1]]
						
						// Apply link style
						styledLink := linkStyle.Render(linkText)
						
						// Replace the link in the line
						formattedLine = formattedLine[:match[0]] + styledLink + formattedLine[match[1]:]
					}
					
					content.WriteString(formattedLine + "\n")
				} else {
					content.WriteString(line + "\n")
				}
			}
		}
		
		// Apply content styling and add to document
		doc.WriteString(contentStyle.Render(content.String()))
		
		// Add consistent spacing between sections
		doc.WriteString("\n")
	}
	
	// Footer with appropriate instructions
	footerText := "Navigate sections with j/k"
	
	// Only show TAB hint if current section has links
	currentSectionLinks := findLinks(sections[m.sectionCursor].content)
	if len(currentSectionLinks) > 0 {
		if m.inLinkMode {
			footerText = "Navigate links with j/k, ENTER to open in browser, TAB to exit link mode"
		} else {
			footerText += ", TAB to select links"
		}
	}
	
	footerText += ", quit with q"
	
	doc.WriteString(inactiveStyle.Render(footerText))
	
	return doc.String()
}

func handleSession(s ssh.Session) {
	pty, _, active := s.Pty()
	if !active {
		fmt.Fprintln(s, "No active terminal, please run with ssh -t")
		return
	}
	
	// Initialize model with terminal dimensions
	m := model{
		sectionCursor: 0,
		linkCursor:    0,
		inLinkMode:    false,
		links:         findLinks(sections[0].content), // Get links for initial section
		width:         pty.Window.Width,
		height:        pty.Window.Height,
	}
	
	p := tea.NewProgram(
		m,
		tea.WithInput(s),
		tea.WithOutput(s),
		tea.WithAltScreen(),
	)
	
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(s, "Error:", err)
	}
}

func main() {
	server := ssh.Server{
		Addr:    ":2222",
		Handler: handleSession,
	}

	log.Println("🔒 TUI SSH server started on port 2222")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
