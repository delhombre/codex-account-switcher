// Package styles provides the visual theme for the TUI.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors for the theme
var (
	Primary   = lipgloss.Color("#7C3AED") // Violet
	Secondary = lipgloss.Color("#A78BFA") // Light violet
	Success   = lipgloss.Color("#10B981") // Emerald
	Warning   = lipgloss.Color("#F59E0B") // Amber
	Error     = lipgloss.Color("#EF4444") // Red
	Muted     = lipgloss.Color("#6B7280") // Gray
	Text      = lipgloss.Color("#F9FAFB") // Light text
	TextDim   = lipgloss.Color("#9CA3AF") // Dimmed text
)

// Box styles
var (
	// BoxStyle is the main container style
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	// HeaderStyle for titles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	// SubHeaderStyle for subtitles
	SubHeaderStyle = lipgloss.NewStyle().
			Foreground(TextDim).
			Italic(true)
)

// Text styles
var (
	// BoldStyle for emphasized text
	BoldStyle = lipgloss.NewStyle().Bold(true)

	// SuccessStyle for success messages
	SuccessStyle = lipgloss.NewStyle().Foreground(Success)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().Foreground(Error)

	// WarningStyle for warning messages
	WarningStyle = lipgloss.NewStyle().Foreground(Warning)

	// MutedStyle for dimmed text
	MutedStyle = lipgloss.NewStyle().Foreground(Muted)

	// PrimaryStyle for primary actions
	PrimaryStyle = lipgloss.NewStyle().Foreground(Primary)
)

// List styles
var (
	// SelectedItemStyle for the currently selected item
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true).
				PaddingLeft(2)

	// NormalItemStyle for unselected items
	NormalItemStyle = lipgloss.NewStyle().
			Foreground(Text).
			PaddingLeft(4)

	// CurrentAccountStyle for marking the active account
	CurrentAccountStyle = lipgloss.NewStyle().
				Foreground(Success).
				Bold(true)
)

// Status indicators - clean Unicode, no emojis
var (
	CheckMark = SuccessStyle.Render("✓")
	CrossMark = ErrorStyle.Render("✗")
	Bullet    = PrimaryStyle.Render("●")
	Circle    = MutedStyle.Render("○")
	Arrow     = PrimaryStyle.Render("→")
	Dash      = MutedStyle.Render("─")
	Caret     = PrimaryStyle.Render("›")
)

// Spinner styles
var (
	SpinnerStyle = lipgloss.NewStyle().Foreground(Primary)
)

// RenderTitle creates a styled title
func RenderTitle(title string) string {
	return HeaderStyle.Render(title)
}

// RenderBox wraps content in a styled box
func RenderBox(content string) string {
	return BoxStyle.Render(content)
}

// RenderSuccess renders a success message
func RenderSuccess(msg string) string {
	return CheckMark + " " + SuccessStyle.Render(msg)
}

// RenderError renders an error message
func RenderError(msg string) string {
	return CrossMark + " " + ErrorStyle.Render(msg)
}

// RenderWarning renders a warning message
func RenderWarning(msg string) string {
	return WarningStyle.Render("! " + msg)
}

// RenderInfo renders an info message
func RenderInfo(msg string) string {
	return Caret + " " + msg
}

