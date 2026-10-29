package tui

import (
	"strings"

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

func (g generalGroup) View(active bool) string {
	s := g.styles

	marker := UNSELECTED_MARKER
	label := lipgloss.NewStyle().Width(GENERAL_LABEL_WIDTH)

	if active && g.focused {
		marker = s.selected.Render(SELECTED_MARKER)
		label = label.Bold(true)
	}

	opts := make([]string, len(g.settings))
	for i, o := range g.settings {
		if i == g.index {
			opts[i] = s.selected.Render("[" + o + "]")
		} else {
			opts[i] = s.dimText.Render(" " + o + " ")
		}
	}

	return marker + label.Render(g.label) + strings.Join(opts, " ")
}

func (g generalGroup) Value() string {
	return g.settings[g.index]
}
