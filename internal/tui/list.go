package tui

import "charm.land/lipgloss/v2"

type list struct {
	values		map[string]string
	options		[]string
	activeIndex	int
}

func (l *list) SetValue(fieldID, value string) {
	l.values[fieldID] = value
}

func (l *list) SetActiveRow(rowIndex int) {
	l.activeIndex = max(0, min(l.activeIndex + rowIndex, len(l.options) - 1))
}

func renderList(l list, s *styles) string {
	var lines []string

	for i, o := range l.options {
		style := lipgloss.NewStyle()
		if i == l.activeIndex {
			style = style.Bold(true)
		}
		line := o + ": " + "[" + l.values[o] + "]"
		lines = append(lines, style.Render(line))
	}

	rows := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return rows
}