package tui

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type movesModel struct {
	columns		[]table.Column
	rows		[][]string
	row			int
	col			int
	styles		*styles
}

func newMovesModel(s *styles) movesModel {
	columns := []table.Column{
		{Title: "Move", Width: 15},
		{Title: "Base Power", Width: 15},
		{Title: "Type", Width: 15},
		{Title: "Category", Width: 15},
		{Title: "Crit", Width: 15},
		{Title: "Hits", Width: 15},
	}

	rows := [][]string{
		{"Move1", "BP1", "Type1", "None", "No", "0"},
		{"Move2", "BP2", "Type2", "None", "No", "0"},
		{"Move3", "BP3", "Type3", "None", "No", "0"},
		{"Move4", "BP4", "Type4", "None", "No", "0"},
	}

	return movesModel{
		columns:	columns,
		rows:		rows,
		row:		0,
		col:		0,
		styles:		s,
	}
}

func (m movesModel) Init() tea.Cmd {
	return nil
}

func (m movesModel) Update(msg tea.Msg) (tabContent, tea.Cmd) {
	switch msg := msg.(type) {
	case editConfirmMsg:
		m.rows[m.row][m.col] = msg.value
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.row = max(m.row - 1, 0)
		case "down":
			m.row = min(m.row + 1, len(m.rows) - 1)
		case "left":
			m.col = max(m.col - 1, 0)
		case "right":
			m.col = min(m.col + 1, len(m.columns) - 1)
		case "enter":
			key := m.columns[m.col].Title
			value := m.rows[m.row][m.col]
			return m, requestTextInput(key, value)
		}
	}
	
	return m, nil
}

func (m movesModel) View() string {
	var lines []string

	var headerCells []string
	for _, c := range m.columns {
		headerCells = append(headerCells, lipgloss.NewStyle().Width(c.Width).Bold(true).Render(c.Title))
	}
	lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))

	for ri, row := range m.rows {
		var cells []string
		for ci, colVal := range row {
			style := lipgloss.NewStyle().Width(m.columns[ci].Width)
			if ri == m.row && ci == m.col {
				style = style.Foreground(m.styles.highlight).Bold(true)
			}
			cells = append(cells, style.Render(colVal))
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
