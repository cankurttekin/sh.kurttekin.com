package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	BaseColor       = lipgloss.Color("#282c34")
	PrimaryColor    = lipgloss.Color("#5f87ff") // Vibrant blue
	AccentColor     = lipgloss.Color("#ff6ac1") // Pink
	SuccessColor    = lipgloss.Color("#98c379") // Green
	WarningColor    = lipgloss.Color("#e5c07b") // Yellow
	DangerColor     = lipgloss.Color("#e06c75") // Red
	TextColor       = lipgloss.Color("#abb2bf") // Light gray
	SubtleColor     = lipgloss.Color("#565c64") // Dark gray
	BackgroundColor = lipgloss.Color("#1e222a") // Very dark blue-gray
	HighlightColor  = lipgloss.Color("#61afef") // Light blue for links
	SelectionColor  = lipgloss.Color("#c678dd") // Purple for selections
)

// Application layout values
var (
	TabWidth = 16

	HorizontalMargin = 2
	VerticalMargin   = 1
)

// Base styles
var (
	// Base text style
	BaseStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	// App container style
	AppStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2).
			BorderBottom(true)

	// Title style for the application header
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor).
			PaddingBottom(1).
			MarginBottom(1).
			Italic(true).
			Border(lipgloss.Border{
			Bottom: "─",
		}).
		BorderForeground(PrimaryColor)

	// Tab bar styles
	TabBarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true).
			BorderForeground(PrimaryColor)

	ActiveTabStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Background(BaseColor).
			Bold(true).
			Padding(0, 2).
			Border(lipgloss.Border{
			Bottom: "─",
		}, false, false, true).
		BorderForeground(AccentColor)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Padding(0, 2)

	// Content container style
	ContentStyle = lipgloss.NewStyle().
			Padding(1, 2).
			MarginTop(1)

	// Navigation styles
	FocusedStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	InactiveStyle = lipgloss.NewStyle().
			Foreground(SubtleColor)

	// Link styles
	LinkStyle = lipgloss.NewStyle().
			Foreground(HighlightColor).
			Underline(true)

	SelectedLinkStyle = lipgloss.NewStyle().
				Foreground(SelectionColor).
				Background(lipgloss.Color("#2c323c")).
				Bold(true).
				Underline(true)

	// Status bar styles
	StatusBarStyle = lipgloss.NewStyle().
			Background(PrimaryColor).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			PaddingLeft(2).
			PaddingRight(2)

	ModeIndicatorStyle = lipgloss.NewStyle().
				Background(AccentColor).
				Foreground(lipgloss.Color("#000000")).
				Bold(true).
				Padding(0, 1)

	StatusMessageStyle = lipgloss.NewStyle().
				Background(SubtleColor).
				Foreground(TextColor).
				Italic(true).
				Padding(0, 1)

	// Section content style
	SectionContentStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				MarginTop(1)

	// Item styles
	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	HighlightedItemStyle = lipgloss.NewStyle().
				Foreground(SuccessColor).
				PaddingLeft(2)
)

// TabBorder returns a customized tab border (straight)
func TabBorder() lipgloss.Border {
	return lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
	}
}

// RenderTabs creates a tab bar from section titles
func RenderTabs(titles []string, activeTab int, width int) string {
	availWidth := width - 4 // Account for margins

	var tabs []string

	// Generate tab styles with proper width handling
	for i, title := range titles {
		var style lipgloss.Style
		if i == activeTab {
			style = ActiveTabStyle.Copy()
		} else {
			style = InactiveTabStyle.Copy()
		}

		// Ensure tab text fits
		tabText := title
		if len(title) > TabWidth-4 {
			tabText = title[:TabWidth-7] + "..."
		}

		// Capitalize the first letter to make tabs more visible
		if len(tabText) > 0 {
			tabText = strings.ToUpper(tabText[:1]) + tabText[1:]
		}

		tabs = append(tabs, style.Render(tabText))
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	if lipgloss.Width(tabBar) > availWidth {
		return lipgloss.NewStyle().Width(availWidth).Render(tabBar)
	}

	return TabBarStyle.Width(availWidth).Render(tabBar)
}

// RenderStatusBar creates a Neovim-like status bar
func RenderStatusBar(mode string, message string, width int) string {
	// Default mode is "NORMAL"
	if mode == "" {
		mode = "NORMAL"
	}

	// Render mode indicator
	modeIndicator := ModeIndicatorStyle.Render(mode)

	// Right side information (status message)
	statusMsg := StatusMessageStyle.Render(message)

	// Calculate remaining space
	remainingWidth := width - lipgloss.Width(modeIndicator) - lipgloss.Width(statusMsg)

	// Create the padding
	padding := lipgloss.NewStyle().
		Background(SubtleColor).
		Width(remainingWidth).
		Render("")

	// Join all parts
	return lipgloss.JoinHorizontal(lipgloss.Top, modeIndicator, padding, statusMsg)
}
