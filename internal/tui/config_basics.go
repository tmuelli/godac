package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"golang.org/x/exp/maps"
)

type basicsModel struct {
	values		side[map[string]string]
	options		[]string
	activeIndex	int
	styles		*styles
}

func newBasicsModel(s *styles) basicsModel {
	template := map[string]string{
		"Name": "",
		"Gender": "",
		"Type1": "",
		"Type2": "",
		"Level": "",
	}

	return basicsModel{
		values:		side[map[string]string]{
			left:	maps.Clone(template),
			right:	maps.Clone(template),
		},
		options:		maps.Keys(template),
		activeIndex:	0,
		styles:			s,
	}
}

func (m basicsModel) Init() tea.Cmd {
	return nil
}

func (m basicsModel) Update(msg tea.Msg, side ConfigSide) (tabContent, tea.Cmd) {
	switch msg := msg.(type) {
	case editConfirmMsg:
		m.values.Get(side)[msg.fieldID] = msg.value
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.activeIndex = max(m.activeIndex - 1, 0)
		case "down":
			m.activeIndex = min(m.activeIndex + 1, len(m.options) - 1)
		case "enter":
			key := m.options[m.activeIndex]
			value := m.values.Get(side)[key]
			return m, requestTextInput(key, value)
		}
	}

	return m, nil
}

func (m basicsModel) View(side ConfigSide) string {
	values := m.values.Get(side)

	var lines []string

	for i, o := range m.options {
		style := lipgloss.NewStyle()
		if i == m.activeIndex {
			style = style.Bold(true)
		}
		line := o + ": " + "[" + values[o] + "]"
		lines = append(lines, style.Render(line))
	}

	rows := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return rows
}