package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"golang.org/x/exp/maps"
)

type basicsModel struct {
	valueMap	map[string]string
	options		[]string
	activeIndex	int
	styles		*styles
}

func newBasicsModel(s *styles) basicsModel {
	valueMap := map[string]string{
		"Name": "",
		"Gender": "",
		"Type1": "",
		"Type2": "",
		"Level": "",
	}

	return basicsModel{
		valueMap:		valueMap,
		options:		maps.Keys(valueMap),
		activeIndex:	0,
		styles:			s,
	}
}

func (m basicsModel) Init() tea.Cmd {
	return nil
}

func (m basicsModel) Update(msg tea.Msg) (tabContent, tea.Cmd) {
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

func (m basicsModel) View() string {
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