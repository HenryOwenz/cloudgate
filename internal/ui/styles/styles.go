package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
)

// Styles holds all the UI styling configurations
type Styles struct {
	App          lipgloss.Style
	Title        lipgloss.Style
	Help         lipgloss.Style
	Context      lipgloss.Style
	Error        lipgloss.Style
	Table        table.Styles
	SearchPrompt lipgloss.Style
	SearchText   lipgloss.Style
}

// DefaultStyles returns the default styling configuration
func DefaultStyles() Styles {
	s := Styles{}

	// Use color constants for consistent styling
	subtle := lipgloss.Color(constants.ColorSubtle)
	highlight := lipgloss.Color(constants.ColorPrimary)
	special := lipgloss.Color(constants.ColorSuccess)
	contextColor := lipgloss.Color(constants.ColorSubtle)
	darkGray := lipgloss.Color(constants.ColorBgAlt)
	titleColor := lipgloss.Color(constants.ColorTitle)
	headerColor := lipgloss.Color(constants.ColorHeader)
	textColor := lipgloss.Color(constants.ColorText)

	s.App = lipgloss.NewStyle().
		Padding(constants.PaddingX, constants.PaddingY)

	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(titleColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(subtle).
		Padding(0, 1)

	s.Help = lipgloss.NewStyle().
		Foreground(subtle).
		MarginTop(1)

	s.Context = lipgloss.NewStyle().
		Foreground(contextColor).
		Height(6).
		Padding(0, 1)

	s.Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color(constants.ColorError)).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(constants.ColorError))

	// Search styles
	s.SearchPrompt = lipgloss.NewStyle().
		Foreground(highlight).
		Bold(true)

	s.SearchText = lipgloss.NewStyle().
		Foreground(textColor)

	// Table styles with fixed height
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(special).
		BorderBottom(true).
		Bold(true).
		Padding(0, 1).
		Foreground(headerColor).
		Align(lipgloss.Center)

	ts.Selected = ts.Selected.
		Foreground(darkGray).
		Background(highlight).
		Bold(true).
		Padding(0, 1)

	ts.Cell = ts.Cell.
		BorderForeground(subtle).
		Padding(0, 1)

	s.Table = ts

	return s
}
