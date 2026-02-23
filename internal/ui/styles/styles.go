// Package styles defines the lipgloss style palette used by the tandem UI.
package styles

import "github.com/charmbracelet/lipgloss"

var (
	Title      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F59E0B"))
	Word       = lipgloss.NewStyle().Foreground(lipgloss.Color("#CBD5E1"))
	Cursor     = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#F59E0B")).Foreground(lipgloss.Color("#1A1A2E"))
	LineNum    = lipgloss.NewStyle().Foreground(lipgloss.Color("#4B5563"))
	Sep        = lipgloss.NewStyle().Foreground(lipgloss.Color("#374151"))
	Status     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	Dim        = lipgloss.NewStyle().Foreground(lipgloss.Color("#4B5563"))
	Analysis   = lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E8F0"))
	POS        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#06B6D4"))
	PanelTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F59E0B"))

	Main = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(0, 1)

	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1, 2)

	Quit = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#EF4444")).
		Padding(1, 2)
)
