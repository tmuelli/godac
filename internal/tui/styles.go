package tui

import (
	"image/color"

	"charm.land/bubbles/v2/help"
	"charm.land/lipgloss/v2"
)

const (
	selectedMarker   = "▸ "
	unselectedMarker = "  "
)

type styles struct {
	// color palette
	accent color.Color // selection, active color
	dim    color.Color // hints, seperators and empty values
	ok     color.Color // success, successful values
	err    color.Color // errors, failed inputs etc.
	warn   color.Color // warnings
	border color.Color // borders, inactive tabs etc.

	// text styles
	title    lipgloss.Style
	dimText  lipgloss.Style
	okText   lipgloss.Style
	errText  lipgloss.Style
	selected lipgloss.Style

	// container styles
	doc    lipgloss.Style
	tab    lipgloss.Style
	window lipgloss.Style

	// help styles
	help help.Styles
}

func newStyles(bgIsDark bool) *styles {
	s := new(styles)
	s.apply(bgIsDark)
	return s
}

func (s *styles) apply(bgIsDark bool) {
	ldFunc := lipgloss.LightDark(bgIsDark)

	s.accent = ldFunc(lipgloss.Color("#0969DA"), lipgloss.Color("#58A6FF"))
	s.dim = ldFunc(lipgloss.Color("#57606A"), lipgloss.Color("#8B949E"))
	s.ok = ldFunc(lipgloss.Color("#1A7F37"), lipgloss.Color("#3FB950"))
	s.err = ldFunc(lipgloss.Color("#CF222E"), lipgloss.Color("#F85149"))
	s.warn = ldFunc(lipgloss.Color("#9A6700"), lipgloss.Color("#D29922"))
	s.border = ldFunc(lipgloss.Color("#AFB8C1"), lipgloss.Color("#484F58"))

	s.title = lipgloss.NewStyle().Bold(true).Foreground(s.accent)
	s.dimText = lipgloss.NewStyle().Foreground(s.dim)
	s.okText = lipgloss.NewStyle().Foreground(s.ok)
	s.errText = lipgloss.NewStyle().Foreground(s.err)
	s.selected = lipgloss.NewStyle().Bold(true).Foreground(s.accent)

	s.doc = lipgloss.NewStyle().Padding(1, 2)

	s.tab = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		BorderForeground(s.border)

	s.window = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(s.border)

	s.help = help.Styles{
		ShortKey:       lipgloss.NewStyle().Foreground(s.accent),
		ShortDesc:      lipgloss.NewStyle().Foreground(s.dim),
		ShortSeparator: lipgloss.NewStyle().Foreground(s.border),
		Ellipsis:       lipgloss.NewStyle().Foreground(s.border),
		FullKey:        lipgloss.NewStyle().Foreground(s.accent),
		FullDesc:       lipgloss.NewStyle().Foreground(s.dim),
		FullSeparator:  lipgloss.NewStyle().Foreground(s.border),
	}
}
