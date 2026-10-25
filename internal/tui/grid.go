package tui

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

type gridColumn struct {
	col      table.Column
	editable bool
	visible  bool
}

type gridData struct {
	columns []gridColumn
	rows    [][]string
	row     int
	col     int
}

func (g *gridData) MoveRow(rowDelta int) {
	g.row = max(0, min(g.row+rowDelta, len(g.rows)-1))
}

func (g *gridData) MoveCol(colDelta int) {
	// g.col = max(0, min(g.col + colDelta, len(g.columns) - 1))
	next := g.col
	for {
		next += colDelta

		if next < 0 || next >= len(g.columns) {
			return
		}

		if g.columns[next].editable && g.columns[next].visible {
			g.col = next
			return
		}
	}
}

func (g *gridData) Set(row, col int, value string) {
	g.rows[row][col] = value
}

func (g *gridData) SetNextValidCol() {
	if g.col < len(g.columns) && g.columns[g.col].editable && g.columns[g.col].visible {
		return
	}

	g.MoveCol(1)

	if !g.columns[g.col].editable || !g.columns[g.col].visible {
		g.MoveCol(-1)
	}
}

func renderGrid(data gridData, s *styles) string {
	var lines []string

	var headerCells []string
	for _, c := range data.columns {
		if !c.visible {
			continue
		}

		headerCells = append(headerCells, lipgloss.NewStyle().
			Width(c.col.Width).
			Bold(true).
			Render(c.col.Title))
	}
	lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))

	for ri, row := range data.rows {
		var cells []string
		for ci, colVal := range row {
			c := data.columns[ci]
			if !c.visible {
				continue
			}

			style := lipgloss.NewStyle().Width(c.col.Width)
			if ri == data.row && ci == data.col {
				style = style.Foreground(s.highlight).Bold(true)
			}
			cells = append(cells, style.Render(colVal))
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func firstEditableCol(columns []gridColumn) int {
	for i, c := range columns {
		if c.editable && c.visible {
			return i
		}
	}
	return 0
}
