package tui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cankurttekin/sh.kurttekin.com/internal/models"
	"github.com/cankurttekin/sh.kurttekin.com/pkg/browser"
)

// Model represents the state of the TUI application
type Model struct {
	SectionCursor int      // Active section
	LinkCursor    int      // Active link
	InLinkMode    bool     // Whether we're in link mode
	Links         []string // Links in the current section
	TabTitles     []string // Tab titles
	Width         int      // Terminal width
	Height        int      // Terminal height
	StatusMode    string   // Status bar mode indicator
	StatusMessage string   // Status bar message
	ShowWelcome   bool     // Whether to show the welcome screen
	Title         string   // Title of the current section
}

// Message when a URL should be opened
type openURLMsg string

// Message to indicate the welcome screen should be dismissed
type welcomeDoneMsg struct{}

// NewModel creates and initializes a new TUI model
func NewModel(width, height int) Model {
	// Create tab titles from section titles
	var tabTitles []string
	for _, sec := range models.Sections {
		tabTitles = append(tabTitles, sec.Title)
	}

	m := Model{
		SectionCursor: 0,
		LinkCursor:    0,
		InLinkMode:    false,
		TabTitles:     tabTitles,
		Width:         width,
		Height:        height,
		StatusMode:    "NORMAL",
		StatusMessage: "Ready",
		ShowWelcome:   true,
		Title:         "cankurttekin",
	}

	// Get links for initial section
	if len(models.Sections) > 0 {
		m.Links = FindLinks(models.Sections[0].Content)
	}

	return m
}

// dismiss the welcome screen after a delay
func welcomeScreenTimer() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return welcomeDoneMsg{}
	})
}

// openURLCommand returns a command to open a URL
func openURLCommand(url string) tea.Cmd {
	return func() tea.Msg {
		err := browser.OpenURL(url)
		if err != nil {
			// Just return nil if there was an error
			return nil
		}
		return openURLMsg(url)
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		welcomeScreenTimer(), // Start the welcome screen timer
	)
}

// Update handles all the business logic and state updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case welcomeDoneMsg:
		// Time to dismiss the welcome screen
		m.ShowWelcome = false
		return m, nil
	case tea.KeyMsg:
		// Dismiss welcome screen immediately on any key press
		if m.ShowWelcome {
			m.ShowWelcome = false
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			// Only toggle link mode if current section has links
			currentSectionLinks := FindLinks(models.Sections[m.SectionCursor].Content)
			if len(currentSectionLinks) > 0 {
				m.InLinkMode = !m.InLinkMode
				m.Links = currentSectionLinks

				if m.InLinkMode {
					m.StatusMode = "LINK"
					m.StatusMessage = fmt.Sprintf("Links: %d", len(m.Links))
					// Reset link cursor when entering link mode
					if m.LinkCursor >= len(m.Links) {
						m.LinkCursor = 0
					}
				} else {
					m.StatusMode = "NORMAL"
					m.StatusMessage = "Ready"
				}
			}
		case "j", "down":
			if m.InLinkMode {
				// Navigate links in current section
				if m.LinkCursor < len(m.Links)-1 {
					m.LinkCursor++
					m.StatusMessage = fmt.Sprintf("Link %d/%d", m.LinkCursor+1, len(m.Links))
				}
			} else {
				// Navigate sections
				if m.SectionCursor < len(models.Sections)-1 {
					m.SectionCursor++
					// Update links for the new section
					m.Links = FindLinks(models.Sections[m.SectionCursor].Content)
					m.InLinkMode = false
					m.StatusMode = "NORMAL"
					m.StatusMessage = fmt.Sprintf("Section: %s", models.Sections[m.SectionCursor].Title)
				}
			}
		case "k", "up":
			if m.InLinkMode {
				// Navigate links in current section
				if m.LinkCursor > 0 {
					m.LinkCursor--
					m.StatusMessage = fmt.Sprintf("Link %d/%d", m.LinkCursor+1, len(m.Links))
				}
			} else {
				// Navigate sections
				if m.SectionCursor > 0 {
					m.SectionCursor--
					// Update links for the new section
					m.Links = FindLinks(models.Sections[m.SectionCursor].Content)
					m.InLinkMode = false
					m.StatusMode = "NORMAL"
					m.StatusMessage = fmt.Sprintf("Section: %s", models.Sections[m.SectionCursor].Title)
				}
			}
		case "enter":
			if m.InLinkMode && m.LinkCursor < len(m.Links) {
				// Open the selected link in a browser
				m.StatusMessage = fmt.Sprintf("Opening: %s", m.Links[m.LinkCursor])
				return m, openURLCommand(m.Links[m.LinkCursor])
			}
		}
	case tea.WindowSizeMsg:
		// Update the model with the new window size
		m.Width = msg.Width
		m.Height = msg.Height
	case openURLMsg:
		// URL was opened
		m.StatusMessage = fmt.Sprintf("Opened: %s", string(msg))
	}
	return m, nil
}

// welcome screen view
func (m Model) renderWelcomeScreen() string {
	// Calculate centered position
	width := m.Width
	height := m.Height

	welcomeMsg := "━━━ " + m.Title + " ━━━"

	welcomeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#61afef"))

	// Render the styled message
	styledMsg := welcomeStyle.Render(welcomeMsg)

	// Center the message in the terminal
	centeredMsg := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		styledMsg,
	)

	return centeredMsg
}

// View renders the current state of the model
func (m Model) View() string {
	if m.Width == 0 {
		// Handle initial rendering when width is not yet known
		return "Loading..."
	}

	// Show welcome screen if needed
	if m.ShowWelcome {
		return m.renderWelcomeScreen()
	}

	// Calculate container dimensions
	containerWidth := m.Width * 2 / 3  // Make the container 2/3 of the terminal width
	contentWidth := containerWidth - 8 // Account for padding and borders

	// Calculate a fixed height for the content area to keep box size consistent
	fixedContentHeight := 16 // A reasonable height that should fit most content

	// Decorative elements for the UI - commented out for now but may be used later
	// divider := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("#5f87ff")).
	// 	Render(strings.Repeat("═", contentWidth))

	// Ornaments for the title
	leftOrnament := lipgloss.NewStyle().Foreground(AccentColor).Render("◇")
	rightOrnament := lipgloss.NewStyle().Foreground(AccentColor).Render("◇")
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#5f87ff")).Render(m.Title)

	// Title with more flair
	titleContent := fmt.Sprintf("%s %s %s", leftOrnament, title, rightOrnament)
	titleStr := TitleStyle.Copy().
		Width(contentWidth).
		BorderStyle(lipgloss.Border{
			Bottom: "━",
		}).
		Align(lipgloss.Center).
		MarginBottom(0).
		PaddingBottom(0).
		Render(titleContent)

	// Ensure tab titles are properly set
	if len(m.TabTitles) == 0 {
		for _, sec := range models.Sections {
			m.TabTitles = append(m.TabTitles, sec.Title)
		}
	}

	// Render tabs with proper width
	tabsStr := lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		Render(RenderTabs(m.TabTitles, m.SectionCursor, contentWidth))

	// Get current section content
	currentSection := models.Sections[m.SectionCursor]

	// Content container with section header
	contentBuilder := strings.Builder{}

	sectionHeader := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render("✦ " + strings.ToUpper(currentSection.Title) + " ✦")

	contentBuilder.WriteString(sectionHeader + "\n")
	contentBuilder.WriteString(lipgloss.NewStyle().
		Foreground(SubtleColor).
		Render(strings.Repeat("─", contentWidth/2)) + "\n\n")

	// Process section content
	for _, line := range currentSection.Content {
		processedLine := line

		// Extract and style links
		if m.InLinkMode {
			// In link mode, highlight and make links selectable
			re := regexp.MustCompile(`(https?://\S+)`)
			matches := re.FindAllStringIndex(line, -1)

			// Process matches from right to left to avoid index shifts
			for j := len(matches) - 1; j >= 0; j-- {
				match := matches[j]
				linkText := line[match[0]:match[1]]

				// Check if this link is selected
				isSelected := false
				for linkIdx, l := range m.Links {
					if l == linkText && linkIdx == m.LinkCursor {
						isSelected = true
						break
					}
				}

				// Apply appropriate style with more emphasis
				var styledLink string
				if isSelected {
					styledLink = SelectedLinkStyle.Copy().
						Background(lipgloss.Color("#2a3040")).
						Foreground(lipgloss.Color("#c678dd")).
						Bold(true).
						Underline(true).
						Render("→ " + linkText)
				} else {
					styledLink = LinkStyle.Render(linkText)
				}

				// Replace the original link with the styled version
				processedLine = processedLine[:match[0]] + styledLink + processedLine[match[1]:]
			}
		} else {
			// Not in link mode, just style links
			re := regexp.MustCompile(`(https?://\S+)`)
			matches := re.FindAllStringIndex(line, -1)

			// Process matches from right to left
			for j := len(matches) - 1; j >= 0; j-- {
				match := matches[j]
				linkText := line[match[0]:match[1]]
				styledLink := LinkStyle.Render(linkText)
				processedLine = processedLine[:match[0]] + styledLink + processedLine[match[1]:]
			}
		}

		// Add the processed line to content
		contentBuilder.WriteString("  " + processedLine + "\n")
	}

	// Style the content area with fixed height
	contentStr := lipgloss.NewStyle().
		Height(fixedContentHeight).
		Render(SectionContentStyle.Render(contentBuilder.String()))

	// Create help text with more flair
	var helpText string
	if m.InLinkMode {
		helpText = "↑/↓: navigate links • ENTER: open link • TAB: exit link mode • q: quit"
	} else {
		helpText = "↑/↓: navigate sections"
		if len(m.Links) > 0 {
			helpText += " • TAB: enter link mode"
		}
		helpText += " • q: quit"
	}

	// Add a nice footer with status
	footer := lipgloss.NewStyle().
		Border(lipgloss.Border{Top: "─"}).
		BorderForeground(SubtleColor).
		Padding(0, 1).
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(helpText)

	// Compose the full view in a more integrated way with consistent spacing
	contentArea := fmt.Sprintf("%s\n%s\n%s\n%s",
		titleStr,
		tabsStr,
		contentStr,
		footer)

	// Wrap everything in a container with proper styling
	wrappedView := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2).
		Width(containerWidth).
		Render(contentArea)

	// Center the wrapped view in the terminal
	centeredView := lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		wrappedView,
	)

	// Return the fully composed and centered view
	return centeredView
}
