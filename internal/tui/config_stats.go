package tui

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
)

type statsData struct {
	columns		[]table.Column
	rows		[][]string
	row			int
	col			int
}

type statsModel struct {
	stats		side[statsData]
	styles		*styles
	data		*models.Data
}

func newStatsModel(s *styles, data *models.Data) statsModel {
	columns := []table.Column{
		{Title: "Stat", Width: 15},
		{Title: "Base", Width: 15},
		{Title: "IV", Width: 15},
		{Title: "EV", Width: 15},
		{Title: "Total", Width: 15},
		{Title: "Boost", Width: 15},
	}

	rows := [][]string{
		{"Hp", "0", "0", "0", "0", "0"},
		{"Atk", "0", "0", "0", "0", "0"},
		{"Def", "0", "0", "0", "0", "0"},
		{"Spa", "0", "0", "0", "0", "0"},
		{"Spd", "0", "0", "0", "0", "0"},
		{"Spe", "0", "0", "0", "0", "0"},
	}

	return statsModel{
		stats:		side[statsData]{
			left:	statsData{
				columns:	columns,
				rows:		cloneRows(rows),
				row:		0,
				col:		0,
			},
			right:	statsData{
				columns:	columns,
				rows:		cloneRows(rows),
				row:		0,
				col:		0,
			},
		},
		styles:		s,
		data:		data,
	}
}

func (m statsModel) Init() tea.Cmd {
	return nil
}

func (m statsModel) Update(msg tea.Msg, side ConfigSide) (tabContent, tea.Cmd) {
	data := m.stats.Get(side)

	switch msg := msg.(type) {
	case editConfirmMsg:
		data.rows[data.row][data.col] = msg.value
		m.stats.Set(side, data)
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
			return m, m.requestInput(key, value)
		}
		m.stats.Set(side, data)
	}
	
	return m, nil
}

func (m statsModel) View(side ConfigSide) string {
	data := m.stats.Get(side)

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

// TODO: maybe make general requestInput
func (m statsModel) requestInput(key, value string) tea.Cmd {
	min := 0
	max := 0
	var options []string
	isNumeric := true

	// TODO: maybe use SliceFROmMapValue somewhere before to not do this everytime again
	// TODO: change to constants where its reasonable
	switch key {
	case "Base":
		min = 1
		max = 255
	case "IV":
		min = 0
		max = 31
	case "EV":
		min = 0
		max = 252
	case "Boost":
		isNumeric = false
		options = []string{"-6", "-5", "-4", "-3", "-2", "-1", "0", "+1", "+2", "+3", "+4", "+5", "+6"}
	}

	if isNumeric {
		return requestNumberInput(key, value, min, max)
	} else {
		return requestSelectInput(key, value, options)
	}
}