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
	maxLineWidth := 0

	// get max width and fill lines
	for _, o := range l.options {
		line := o + ": " + "[" + l.values[o] + "]"
		maxLineWidth = max(maxLineWidth, len(line))
		lines = append(lines, line)
	}

	var styledLines []string
	for i, line := range lines {
		style := lipgloss.NewStyle().Width(maxLineWidth + 20)
		if i == l.activeIndex {
			style = style.Bold(true)
		}
		styledLines = append(styledLines, style.Render(line))
	}

	rows := lipgloss.JoinVertical(lipgloss.Left, styledLines...)
	return rows
}