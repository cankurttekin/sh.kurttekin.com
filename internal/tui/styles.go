package tui

import (
	"strings"

	"github.com/cankurttekin/sh.kurttekin.com/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// Style management for the entire application
// All styles and colors are defined here to ensure consistency

// Portfolio instance to access theme data
var portfolio = models.GetPortfolio()

// Color palette - all colors used in the application should be defined here
var (
	BaseColor       = lipgloss.Color("#282c34")
	PrimaryColor    = lipgloss.Color(portfolio.Theme.Primary)
	AccentColor     = lipgloss.Color(portfolio.Theme.Accent)
	SuccessColor    = lipgloss.Color("#98c379") // Green
	WarningColor    = lipgloss.Color("#e5c07b") // Yellow
	DangerColor     = lipgloss.Color("#e06c75") // Red
	TextColor       = lipgloss.Color(portfolio.Theme.Text)
	SubtleColor     = lipgloss.Color(portfolio.Theme.Subtle)
	BackgroundColor = lipgloss.Color("#1e222a") // Very dark blue-gray
	HighlightColor  = lipgloss.Color(portfolio.Theme.Links)
	SelectionColor  = lipgloss.Color(portfolio.Theme.Selection)
	LinkBackground  = lipgloss.Color("#2a3040") // Background for selected links
)

// Layout constants
var (
	TabWidth         = 16
	HorizontalMargin = 2
	VerticalMargin   = 1
	ContentHeight    = 16 // Fixed height for content area
)

// Base text styles
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
			Bottom: "━", // More decorative bottom border
		}).
		BorderForeground(PrimaryColor)

	// Content container style
	ContentStyle = lipgloss.NewStyle().
			Padding(1, 2).
			MarginTop(1)
)

// Welcome screen styles
var (
	WelcomeTextStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(HighlightColor)
)

// Tab styles
var (
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
)

// Navigation styles
var (
	FocusedStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	InactiveStyle = lipgloss.NewStyle().
			Foreground(SubtleColor)
)

// Link styles
var (
	LinkStyle = lipgloss.NewStyle().
			Foreground(HighlightColor).
			Underline(true)

	SelectedLinkStyle = lipgloss.NewStyle().
				Foreground(SelectionColor).
				Background(LinkBackground).
				Bold(true).
				Underline(true)
)

// Status bar styles
var (
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
)

// Content styles
var (
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

	// Section header style
	SectionHeaderStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true)

	// Section divider style
	SectionDividerStyle = lipgloss.NewStyle().
				Foreground(SubtleColor)

	// Footer style
	FooterStyle = lipgloss.NewStyle().
			Border(lipgloss.Border{Top: "─"}).
			BorderForeground(SubtleColor).
			Padding(0, 1).
			Align(lipgloss.Center)

	// Title ornament style
	OrnamentStyle = lipgloss.NewStyle().
			Foreground(AccentColor)

	// Main container style
	ContainerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2)
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
