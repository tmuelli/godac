package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type generalGroup struct {
	label    string
	settings []string
	index    int
	focused  bool
	onChange func(value string) tea.Cmd

	styles *styles
}

func newGeneralGroup(s *styles, label string, defaultIndex int, settings []string) generalGroup {
	return generalGroup{
		label:    label,
		settings: settings,
		index:    defaultIndex,
		styles:   s,
	}
}

func (g generalGroup) Update(msg tea.Msg) (generalGroup, tea.Cmd) {
	if !g.focused {
		return g, nil
	}

	prevIndex := g.index

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, globalKeys.Left):
			g.index = max(g.index-1, 0)
		case key.Matches(msg, globalKeys.Right):
			g.index = min(g.index+1, len(g.settings)-1)
		}
	}

	if g.onChange != nil && g.index != prevIndex {
		return g, g.onChange(g.settings[g.index])
	}

	return g, nil
}

func (g generalGroup) View() string {
	var opts []string

	for i, s := range g.settings {
		style := g.styles.tab
		if i == g.index {
			style = style.BorderForeground(g.styles.highlight).Bold(true)
		}
		opts = append(opts, style.Render(s))
	}

	labelStyle := lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.Border{Left: " ", Right: " "})
	if g.focused {
		labelStyle = labelStyle.Bold(true)
	}

	rowContent := append([]string{labelStyle.Render(g.label + ": ")}, opts...)
	row := lipgloss.JoinHorizontal(lipgloss.Top, rowContent...)
	return row
}

func (m fieldModel) FormatGeneralSettings() string {
	var title string
	var lines []string

	titleStyle := lipgloss.NewStyle()
	if m.setting == GeneralSetting {
		titleStyle = titleStyle.Bold(true).Foreground(lipgloss.Color("#FFFF00"))
	} else {
		titleStyle.Foreground(lipgloss.NoColor{})
	}

	title = "General"

	for _, g := range m.general {
		lines = append(lines, g.View())
	}

	return titleStyle.Render(title) + ":\n" + lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (g generalGroup) Value() string {
	return g.settings[g.index]
}
