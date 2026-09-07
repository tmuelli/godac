package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type styles struct {
	doc			lipgloss.Style
	highlight	color.Color
	tab			lipgloss.Style
	window		lipgloss.Style
}

func newStyles(bgIsDark bool) *styles {
	s := new(styles)

	s.highlight = lipgloss.Color("241")

	s.doc = lipgloss.NewStyle().
		Padding(1, 2)

	s.tab = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	s.window = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.NormalBorder())

	return s
}
