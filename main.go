package main

import (
	"fmt"
	"log"
	"strings"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ssh "github.com/charmbracelet/ssh"
)

var (
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).PaddingBottom(1)
	sectionStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	focusedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	inactiveStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	contentStyle    = lipgloss.NewStyle().PaddingLeft(4)
	linkStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	selectedLinkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Underline(true)
)

type model struct {
	sectionCursor int
	linkCursor    int
	inLinkMode    bool
	links         []string
	width         int
	height        int
}

type section struct {
	title   string
	content []string
}

var sections = []section{
	{"About", []string{
		"I'm a developer building TUI apps in Go.",
		"Passionate about clean terminals and elegant CLI interfaces.",
	}},
	{"Projects", []string{
		"🔧 TUI Portfolio (this one!)",
		"📦 CLI Package Manager UI",
	}},
	{"Experience", []string{
		"💼 Software Engineer @ TerminalTech (2022 - Present)",
		"🧪 QA Tester @ ShellOps (2020 - 2022)",
	}},
	{"Links", []string{
		"GitHub: https://github.com/you",
		"LinkedIn: https://linkedin.com/in/you",
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

func (m model) Init() tea.Cmd {
	// Extract all links from content
	var allLinks []string
	for _, sec := range sections {
		links := findLinks(sec.content)
		allLinks = append(allLinks, links...)
	}
	
	// No command to return
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			// Toggle between section navigation and link selection
			if m.inLinkMode {
				m.inLinkMode = false
			} else {
				// Only enter link mode if there are links
				links := []string{}
				for _, sec := range sections {
					secLinks := findLinks(sec.content)
					links = append(links, secLinks...)
				}
				
				if len(links) > 0 {
					m.inLinkMode = true
					m.links = links
					// Reset link cursor if it's out of bounds
					if m.linkCursor >= len(links) {
						m.linkCursor = 0
					}
				}
			}
		case "j", "down":
			if m.inLinkMode {
				// Navigate links
				if m.linkCursor < len(m.links)-1 {
					m.linkCursor++
				}
			} else {
				// Navigate sections
				if m.sectionCursor < len(sections)-1 {
					m.sectionCursor++
				}
			}
		case "k", "up":
			if m.inLinkMode {
				// Navigate links
				if m.linkCursor > 0 {
					m.linkCursor--
				}
			} else {
				// Navigate sections
				if m.sectionCursor > 0 {
					m.sectionCursor--
				}
			}
		case "enter":
			if m.inLinkMode && len(m.links) > 0 {
				// In a real app, you'd handle opening the link here
				// For now we just return to section mode
				m.inLinkMode = false
			}
		}
	case tea.WindowSizeMsg:
		// Update the model with the new window size
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	// Container for all content
	doc := strings.Builder{}
	
	// Title
	doc.WriteString(titleStyle.Render("🌿 My TUI Portfolio") + "\n\n")
	
	// Find all links for highlighting
	allLinks := []string{}
	for _, sec := range sections {
		secLinks := findLinks(sec.content)
		allLinks = append(allLinks, secLinks...)
	}
	m.links = allLinks
	
	// Render each section with proper padding and styling
	for i, sec := range sections {
		// Determine cursor and style based on selection
		cursor := "  "
		style := inactiveStyle
		if i == m.sectionCursor && !m.inLinkMode {
			cursor = "➜ "
			style = focusedStyle
		}
		
		// Render section header
		sectionHeader := style.Render(cursor + sec.title)
		doc.WriteString(sectionHeader + "\n")
		
		// Always render content for all sections
		content := strings.Builder{}
		for _, line := range sec.content {
			// Check if this line contains links that need highlighting
			if m.inLinkMode {
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
						if l == linkText && linkIdx == m.linkCursor && m.inLinkMode {
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
				content.WriteString(line + "\n")
			}
		}
		
		// Apply content styling and add to document
		doc.WriteString(contentStyle.Render(content.String()))
		
		// Add consistent spacing between sections
		doc.WriteString("\n")
	}
	
	// Footer with appropriate instructions
	footerText := "Navigate sections with j/k"
	if len(m.links) > 0 {
		footerText += ", TAB to toggle link selection"
	}
	if m.inLinkMode {
		footerText = "Navigate links with j/k, ENTER to select, TAB to exit link mode"
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
