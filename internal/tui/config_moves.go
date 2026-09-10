package tui

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type movesData struct {
	columns		[]table.Column
	rows		[][]string
	row			int
	col			int
}

type movesModel struct {
	moves		side[movesData]
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
		moves:		side[movesData]{
			left:	movesData{
				columns:	columns,
				rows:		cloneRows(rows),
				row:		0,
				col:		0,
			},
			right:	movesData{
				columns:	columns,
				rows:		cloneRows(rows),
				row:		0,
				col:		0,
			},
		},
		styles:		s,
	}
}

func (m movesModel) Init() tea.Cmd {
	return nil
}

func (m movesModel) Update(msg tea.Msg, side ConfigSide) (tabContent, tea.Cmd) {
	data := m.moves.Get(side)
	
	switch msg := msg.(type) {
	case editConfirmMsg:
		data.rows[data.row][data.col] = msg.value
		m.moves.Set(side, data)
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			data.row = max(data.row - 1, 0)
		case "down":
			data.row = min(data.row + 1, len(data.rows) - 1)
		case "left":
			data.col = max(data.col - 1, 0)
		case "right":
			data.col = min(data.col + 1, len(data.columns) - 1)
		case "enter":
			key := data.columns[data.col].Title
			value := data.rows[data.row][data.col]
			return m, requestTextInput(key, value)
		}
		m.moves.Set(side, data)
	}
	
	return m, nil
}

func (m movesModel) View(side ConfigSide) string {
	data := m.moves.Get(side)

	var lines []string

	var headerCells []string
	for _, c := range data.columns {
		headerCells = append(headerCells, lipgloss.NewStyle().Width(c.Width).Bold(true).Render(c.Title))
	}
	lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))

	for ri, row := range data.rows {
		var cells []string
		for ci, colVal := range row {
			style := lipgloss.NewStyle().Width(data.columns[ci].Width)
			if ri == data.row && ci == data.col {
				style = style.Foreground(m.styles.highlight).Bold(true)
			}
			cells = append(cells, style.Render(colVal))
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m movesModel) Names(side ConfigSide) []string {
	var names []string
	data := m.moves.Get(side)
	for _, row := range data.rows {
		names = append(names, row[0])
	}

	return names
}
