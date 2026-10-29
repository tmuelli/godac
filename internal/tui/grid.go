package tui

import (
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

const (
	COLUMN_GAP  = 2
	EMPTY_VALUE = "-"
)

type gridColumn struct {
	col          table.Column
	editable     bool
	visible      bool
	rightAligned bool
}

type gridData struct {
	columns []gridColumn
	rows    [][]string
	row     int
	col     int
}

func orDash(val string) string {
	if val == "" {
		return EMPTY_VALUE
	}

	return val
}

// set column widths to the largest width in col (header or column value)
func columnWidths(data gridData) []int {
	widths := make([]int, len(data.columns))
	for i, c := range data.columns {
		if !c.visible {
			continue
		}

		w := max(c.col.Width, lipgloss.Width(c.col.Title))
		for _, row := range data.rows {
			w = max(w, lipgloss.Width(orDash(row[i])))
		}

		widths[i] = w
	}

	return widths
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
	widths := columnWidths(data)
	gap := strings.Repeat(" ", COLUMN_GAP)

	cellStyle := func(i int) lipgloss.Style {
		style := lipgloss.NewStyle().Width(widths[i])
		if data.columns[i].rightAligned {
			style = style.Align(lipgloss.Right)
		}

		return style
	}

	// header and seperators
	var header []string
	total := 0
	for i, c := range data.columns {
		if !c.visible {
			continue
		}

		header = append(header, cellStyle(i).Bold(true).Render(c.col.Title))
		total += widths[i]
	}
	total += COLUMN_GAP * (len(header) - 1)

	lines := []string{
		UNSELECTED_MARKER + strings.Join(header, gap),
		UNSELECTED_MARKER + lipgloss.NewStyle().Foreground(s.border).Render(strings.Repeat("─", total)),
	}

	for ri, row := range data.rows {
		activeMarker := UNSELECTED_MARKER
		if ri == data.row {
			activeMarker = s.selected.Render(SELECTED_MARKER)
		}

		var cells []string
		for ci, colVal := range row {
			c := data.columns[ci]
			if !c.visible {
				continue
			}

			style := cellStyle(ci)
			switch {
			case ri == data.row && ci == data.col:
				style = style.Inherit(s.selected)
			case !c.editable:
				style = style.Inherit(s.dimText)
			default:
				style = style.Inherit(s.valueStyle(colVal))
			}

			cells = append(cells, style.Render(orDash(colVal)))
		}
		lines = append(lines, activeMarker+strings.Join(cells, gap))
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
