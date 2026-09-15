package tui

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

type gridData struct {
	columns		[]table.Column
	rows		[][]string
	row			int
	col			int
}

func (g *gridData) MoveRow(rowDelta int) {
	g.row = max(0, min(g.row + rowDelta, len(g.rows) - 1))
}

// FIXME: implement editable logic on columns so some columns cannot be edited / focused
func (g *gridData) MoveCol(colDelta int) {
	g.col = max(0, min(g.col + colDelta, len(g.columns) - 1))
}

func (g *gridData) Set(row, col int, value string) {
	g.rows[row][col] = value
}

func renderGrid(data gridData, s *styles) string {
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
				style = style.Foreground(s.highlight).Bold(true)
			}
			cells = append(cells, style.Render(colVal))
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}