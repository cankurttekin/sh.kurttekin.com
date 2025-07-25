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

type Model struct {
	SectionCursor int              // Active section
	LinkCursor    int              // Active link
	InLinkMode    bool             // Whether we're in link mode
	Links         []string         // Links in the current section
	TabTitles     []string         // Tab titles
	Width         int              // Terminal width
	Height        int              // Terminal height
	StatusMode    string           // Status bar mode indicator
	StatusMessage string           // Status bar message
	ShowWelcome   bool             // Whether to show the welcome screen
	Portfolio     models.Portfolio // Portfolio data
}

// message when a URL should be opened
type openURLMsg string

// message to indicate the welcome screen should be dismissed
type welcomeDoneMsg struct{}

// initializes a new TUI model
func NewModel(width, height int) Model {
	portfolio := models.GetPortfolio()

	// create tab titles from section titles
	var tabTitles []string
	for _, sec := range portfolio.Sections {
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
		Portfolio:     portfolio,
	}

	// get links for initial section
	if len(portfolio.Sections) > 0 {
		m.Links = FindLinks(portfolio.Sections[0].Content)
	}

	return m
}

// dismiss the welcome screen after a delay
func welcomeScreenTimer() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return welcomeDoneMsg{}
	})
}

func openURLCommand(url string) tea.Cmd {
	return func() tea.Msg {
		err := browser.OpenURL(url)
		if err != nil {
			return nil
		}
		return openURLMsg(url)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		welcomeScreenTimer(),
	)
}

// update handles all the business logic and state updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case welcomeDoneMsg:
		// time to dismiss the welcome screen
		m.ShowWelcome = false
		return m, nil
	case tea.KeyMsg:
		// dismiss welcome screen immediately on any key press
		if m.ShowWelcome {
			m.ShowWelcome = false
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			// only toggle link mode if current section has links
			currentSectionLinks := FindLinks(m.Portfolio.Sections[m.SectionCursor].Content)
			if len(currentSectionLinks) > 0 {
				m.InLinkMode = !m.InLinkMode
				m.Links = currentSectionLinks

				if m.InLinkMode {
					m.StatusMode = "LINK"
					m.StatusMessage = fmt.Sprintf("Links: %d", len(m.Links))
					// reset link cursor when entering link mode
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
				if m.SectionCursor < len(m.Portfolio.Sections)-1 {
					m.SectionCursor++
					// Update links for the new section
					m.Links = FindLinks(m.Portfolio.Sections[m.SectionCursor].Content)
					m.InLinkMode = false
					m.StatusMode = "NORMAL"
					m.StatusMessage = fmt.Sprintf("Section: %s", m.Portfolio.Sections[m.SectionCursor].Title)
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
					m.Links = FindLinks(m.Portfolio.Sections[m.SectionCursor].Content)
					m.InLinkMode = false
					m.StatusMode = "NORMAL"
					m.StatusMessage = fmt.Sprintf("Section: %s", m.Portfolio.Sections[m.SectionCursor].Title)
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

	// Simple welcome message
	welcomeMsg := "━━━ " + m.Portfolio.Title + " ━━━"

	// Use the consolidated WelcomeTextStyle from styles.go
	styledMsg := WelcomeTextStyle.Render(welcomeMsg)

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
	containerWidth := m.Width * 2 / 3
	contentWidth := containerWidth - 8

	// Ornaments for the title using style from styles.go
	leftOrnament := OrnamentStyle.Render("◇")
	rightOrnament := OrnamentStyle.Render("◇")
	title := lipgloss.NewStyle().Foreground(PrimaryColor).Render(m.Portfolio.Title)

	titleContent := fmt.Sprintf("%s %s %s", leftOrnament, title, rightOrnament)
	titleStr := TitleStyle.Copy().
		Width(contentWidth).
		Align(lipgloss.Center).
		MarginBottom(0).
		PaddingBottom(0).
		Render(titleContent)

	// ensuring tab titles are properly set
	if len(m.TabTitles) == 0 {
		for _, sec := range m.Portfolio.Sections {
			m.TabTitles = append(m.TabTitles, sec.Title)
		}
	}

	// render tabs with proper width
	tabsStr := lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		Render(RenderTabs(m.TabTitles, m.SectionCursor, contentWidth))

	// Get current section content
	currentSection := m.Portfolio.Sections[m.SectionCursor]

	// Content container with section header
	contentBuilder := strings.Builder{}

	// Use SectionHeaderStyle from styles.go
	sectionHeader := SectionHeaderStyle.Render("✦" + strings.ToUpper(currentSection.Title) + "✦")
	contentBuilder.WriteString(sectionHeader + "\n")

	// Use SectionDividerStyle from styles.go
	contentBuilder.WriteString(SectionDividerStyle.Render(strings.Repeat("─", contentWidth/2)) + "\n\n")

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

				var styledLink string
				if isSelected {
					styledLink = SelectedLinkStyle.Copy().
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
			re := regexp.MustCompile(`(https?://\S+)`)
			matches := re.FindAllStringIndex(line, -1)

			for j := len(matches) - 1; j >= 0; j-- {
				match := matches[j]
				linkText := line[match[0]:match[1]]
				var styledLink string
				styledLink = LinkStyle.Render(linkText)
				processedLine = processedLine[:match[0]] + styledLink + processedLine[match[1]:]
			}
		}

		// Add the processed line to content
		contentBuilder.WriteString("  " + processedLine + "\n")
	}

	// Style the content area with fixed height from styles.go
	contentStr := lipgloss.NewStyle().
		Height(ContentHeight).
		Render(SectionContentStyle.Render(contentBuilder.String()))

	var helpText string
	if m.InLinkMode {
		helpText = "↑/↓: navigate links • enter: open link • tab: exit link mode • q: quit"
	} else {
		helpText = "↑/↓ : navigate sections"
		if len(m.Links) > 0 {
			helpText += " • tab: enter link mode"
		}
		helpText += " • q: quit"
	}

	footer := FooterStyle.
		Width(contentWidth).
		Render(helpText)

	statusBar := StatusBarStyle.
		Render(RenderStatusBar(m.StatusMode, m.StatusMessage, contentWidth))

	contentArea := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		titleStr,
		tabsStr,
		contentStr,
		footer,
		statusBar)

	wrappedView := ContainerStyle.
		Width(containerWidth).
		Render(contentArea)

	centeredView := lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		wrappedView,
	)

	return centeredView
}
