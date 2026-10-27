package tui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
)

var displayBoosts = []string{"-6", "-5", "-4", "-3", "-2", "-1", "0", "+1", "+2", "+3", "+4", "+5", "+6"}

type statsModel struct {
	tab    side[gridData]
	styles *styles
	data   *models.Data
}

const (
	STAT_COL = iota
	BASE_COL
	IV_COL
	EV_COL
	TOTAL_COL
	BOOST_COL

	STAT_COL_NAME  = "Stat"
	BASE_COL_NAME  = "Base"
	IV_COL_NAME    = "IV"
	EV_COL_NAME    = "EV"
	TOTAL_COL_NAME = "Total"
	BOOST_COL_NAME = "Boost"
)

const (
	HP_ROW = iota
	ATK_ROW
	DEF_ROW
	SPA_ROW
	SPD_ROW
	SPE_ROW
)

func newStatsModel(s *styles, data *models.Data) statsModel {
	columns := []gridColumn{
		{col: table.Column{Title: STAT_COL_NAME, Width: 15}, editable: false, visible: true},
		{col: table.Column{Title: BASE_COL_NAME, Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: IV_COL_NAME, Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: EV_COL_NAME, Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: TOTAL_COL_NAME, Width: 15}, editable: false, visible: true},
		{col: table.Column{Title: BOOST_COL_NAME, Width: 15}, editable: true, visible: true},
	}

	rows := [][]string{
		{DISPLAY_HP, "0", "0", "0", "0", "0"},
		{DISPLAY_ATK, "0", "0", "0", "0", "0"},
		{DISPLAY_DEF, "0", "0", "0", "0", "0"},
		{DISPLAY_SPA, "0", "0", "0", "0", "0"},
		{DISPLAY_SPD, "0", "0", "0", "0", "0"},
		{DISPLAY_SPE, "0", "0", "0", "0", "0"},
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
		if data.row == HP_ROW && msg.value != "" {
			return m, requestHpRefresh()
		}
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
	switch key {
	case BASE_COL_NAME:
		min = models.BASE_MIN
		max = models.BASE_MAX
	case IV_COL_NAME:
		min = models.IV_MIN
		max = models.IV_MAX
	case EV_COL_NAME:
		min = models.EV_MIN
		max = models.EV_MAX
	case BOOST_COL_NAME:
		isNumeric = false
		options = displayBoosts
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
		DISPLAY_HP:  0,
		DISPLAY_ATK: 0,
		DISPLAY_DEF: 0,
		DISPLAY_SPA: 0,
		DISPLAY_SPD: 0,
		DISPLAY_SPE: 0,
	}

	for _, row := range data.rows {
		value := atoiOrZero(row[col])
		values[row[STAT_COL]] = value
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
		switch data.rows[rowIndex][STAT_COL] {
		case DISPLAY_HP:
			value = fmt.Sprintf("%d", baseStats.Hp)
		case DISPLAY_ATK:
			value = fmt.Sprintf("%d", baseStats.Atk)
		case DISPLAY_DEF:
			value = fmt.Sprintf("%d", baseStats.Def)
		case DISPLAY_SPA:
			value = fmt.Sprintf("%d", baseStats.Spa)
		case DISPLAY_SPD:
			value = fmt.Sprintf("%d", baseStats.Spd)
		case DISPLAY_SPE:
			value = fmt.Sprintf("%d", baseStats.Spe)
		default:
			continue
		}

		data.Set(rowIndex, 1, value)

		// default to max iv
		data.Set(rowIndex, 2, fmt.Sprintf("%d", models.IV_MAX))
	}
}
