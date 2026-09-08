package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"golang.org/x/exp/maps"
)

type miscModel struct {
	valueMap	map[string]string
	options		[]string
	activeIndex	int
	styles		*styles
}

func newMiscModel(s *styles) miscModel {
	valueMap := map[string]string{
		"Nature": "",
		"Ability": "",
		"Item": "",
		"Status": "",
		"Health": "",
	}

	return miscModel{
		valueMap:		valueMap,
		options:		maps.Keys(valueMap),
		activeIndex:	0,
		styles:			s,
	}
}

func (m miscModel) Init() tea.Cmd {
	return nil
}

func (m miscModel) Update(msg tea.Msg) (tabContent, tea.Cmd) {
	switch msg := msg.(type) {
	case editConfirmMsg:
		m.valueMap[msg.fieldID] = msg.value
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.activeIndex = max(m.activeIndex - 1, 0)
		case "down":
			m.activeIndex = min(m.activeIndex + 1, len(m.options) - 1)
		case "enter":
			key := m.options[m.activeIndex]
			value := m.valueMap[key]
			return m, requestTextInput(key, value)
		}
	}

	return m, nil
}

func (m miscModel) View() string {
	var lines []string

	for i, o := range m.options {
		style := lipgloss.NewStyle()
		if i == m.activeIndex {
			style = style.Bold(true)
		}
		lines = append(lines, style.Render(o, ": ", "[", m.valueMap[o], "]"))
	}

	rows := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return rows
}