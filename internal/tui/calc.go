package tui

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
)

const (
	LEFT = iota
	RIGHT
)

type calcContext struct {
	leftName	string
	rightName	string
	leftMoves	[]string
	rightMoves	[]string
}

type calcModel struct {
	columns		[]table.Column
	rows		[][]string
	row			int
	col			int
	styles		*styles
	data		*models.Data
}

func newCalcModel(s *styles, data *models.Data) calcModel {
	columns := []table.Column{
		{Title: "Pokemon1 [L]", Width: 15},
		{Title: "Pokemon2 [R]", Width: 15},
	}

	rows := [][]string{
		{"Move 1 [L1]", "Move 1 [R1]"},
		{"Move 2 [L2]", "Move 2 [R2]"},
		{"Move 3 [L3]", "Move 3 [R3]"},
		{"Move 4 [L4]", "Move 4 [R4]"},
	}

	return calcModel{
		columns:	columns,
		rows:		rows,
		row:		0,
		col:		0,
		styles:		s,
		data:		data,
	}
}

func (m calcModel) Init() tea.Cmd {
	return nil
}

func (m calcModel) Update(msg tea.Msg) (calcModel, tea.Cmd) {
	switch msg := msg.(type) {
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
		}
	}

	return m, nil
}

func (m calcModel) View(ctx calcContext) string {
	content := m.Title()
	content += "\n\n"
	content += m.Body(ctx)
	return content
}

func (m calcModel) Title() string {
	return "Calculation view - Press [R] or [L] for config view or [F] for field view"
}

func (m calcModel) Body(ctx calcContext) string {
	m.fillFromContext(ctx)

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

	body := lipgloss.JoinVertical(lipgloss.Left, lines...)
	body += "\n\nChoose one move\n"
	body += "\nPokemon1 (Move x) vs. Pokemon 2: MinDmg - MaxDmg (MinPerc - MaxPerc)\n"
	body += "\nPossible dmg amounts: (x, ...)\n"
	return body
}

func (m calcModel) fillFromContext(ctx calcContext) {
	// left
	m.fillSideFromContext(ctx.leftName, ctx.leftMoves, LEFT, "L")

	// right
	m.fillSideFromContext(ctx.rightName, ctx.rightMoves, RIGHT, "R")
}

func (m calcModel) fillSideFromContext(sideName string, sideMoves []string, side int, sidePrefix string) {
	if sideName != "" {
		m.columns[side].Title = fmt.Sprintf("%s [%s]", sideName, sidePrefix)
	}

	for i, move := range sideMoves {
		if move == "" {
			continue
		}

		m.rows[i][side] = fmt.Sprintf("%s [%s%d]", move, sidePrefix, i+1)
	}
}
