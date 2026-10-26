package tui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
)

type statsModel struct {
	tab    side[gridData]
	styles *styles
	data   *models.Data
}

func newStatsModel(s *styles, data *models.Data) statsModel {
	columns := []gridColumn{
		{col: table.Column{Title: "Stat", Width: 15}, editable: false, visible: true},
		{col: table.Column{Title: "Base", Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: "IV", Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: "EV", Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: "Total", Width: 15}, editable: false, visible: true},
		{col: table.Column{Title: "Boost", Width: 15}, editable: true, visible: true},
	}

	rows := [][]string{
		{"Hp", "0", "0", "0", "0", "0"},
		{"Atk", "0", "0", "0", "0", "0"},
		{"Def", "0", "0", "0", "0", "0"},
		{"Spa", "0", "0", "0", "0", "0"},
		{"Spd", "0", "0", "0", "0", "0"},
		{"Spe", "0", "0", "0", "0", "0"},
	}

	initialColIndex := firstEditableCol(columns)

	return statsModel{
		tab: side[gridData]{
			left: gridData{
				columns: cloneColumns(columns),
				rows:    cloneRows(rows),
				row:     0,
				col:     initialColIndex,
			},
			right: gridData{
				columns: cloneColumns(columns),
				rows:    cloneRows(rows),
				row:     0,
				col:     initialColIndex,
			},
		},
		styles: s,
		data:   data,
	}
}

func (m statsModel) Init() tea.Cmd {
	return nil
}

func (m statsModel) Update(msg tea.Msg, side Side) (tabContent, tea.Cmd) {
	data := m.tab.GetPtr(side)

	switch msg := msg.(type) {
	case autoFillForSpecieRequestMsg:
		m.fillStatsDefaults(data, msg.specie)
		return m, nil

	case editConfirmMsg:
		data.Set(data.row, data.col, msg.value)
		return m, nil
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, globalKeys.Up):
			data.MoveRow(-1)
		case key.Matches(msg, globalKeys.Down):
			data.MoveRow(1)
		case key.Matches(msg, globalKeys.Left):
			data.MoveCol(-1)
		case key.Matches(msg, globalKeys.Right):
			data.MoveCol(1)
		case key.Matches(msg, globalKeys.Confirm):
			if !data.columns[data.col].editable {
				return m, nil
			}

			key := data.columns[data.col].col.Title
			value := data.rows[data.row][data.col]
			return m, m.requestInput(key, value)
		}
	}

	return m, nil
}

func (m statsModel) View(side Side) string {
	return renderGrid(m.tab.Get(side), m.styles)
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

func (m statsModel) columnValues(side Side, col int) map[string]int {
	data := m.tab.Get(side)
	values := map[string]int{
		"Hp":  0,
		"Atk": 0,
		"Def": 0,
		"Spa": 0,
		"Spd": 0,
		"Spe": 0,
	}

	for _, row := range data.rows {
		value := atoiOrZero(row[col])
		values[row[0]] = value
	}

	return values
}

func (m statsModel) fillStatsDefaults(data *gridData, specie *models.Specie) {
	if specie == nil {
		// internal error
		return
	}

	baseStats := specie.BaseStats

	for rowIndex := range data.rows {
		// set base stats
		var value string
		switch data.rows[rowIndex][0] {
		case "Hp":
			value = fmt.Sprintf("%d", baseStats.Hp)
		case "Atk":
			value = fmt.Sprintf("%d", baseStats.Atk)
		case "Def":
			value = fmt.Sprintf("%d", baseStats.Def)
		case "Spa":
			value = fmt.Sprintf("%d", baseStats.Spa)
		case "Spd":
			value = fmt.Sprintf("%d", baseStats.Spd)
		case "Spe":
			value = fmt.Sprintf("%d", baseStats.Spe)
		default:
			continue
		}

		data.Set(rowIndex, 1, value)

		// default to max iv
		data.Set(rowIndex, 2, fmt.Sprintf("%d", models.IV_MAX))
	}
}
