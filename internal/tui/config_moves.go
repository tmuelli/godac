package tui

import (
	"fmt"
	"slices"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
)

type MoveColumnType int

const (
	MOVE_COL MoveColumnType = iota
	BP_COL
	MOVE_TYPE_COL
	CATEGORY_COL
	CRIT_COL
	HITS_COL
	METR_HITS_COL

	MOVE_COL_NAME       = "Move"
	BASE_POWER_COL_NAME = "Base Power"
	TYPE_COL_NAME       = "Type"
	CATEGORY_COL_NAME   = "Category"
	CRIT_COL_NAME       = "Crit"
	HITS_COL_NAME       = "Hits"
	METR_HITS_COL_NAME  = "Metr. Hits"
)

type moveRowSettings struct {
	isSet bool
}

type movesModel struct {
	tab         side[gridData]
	rowSettings side[[]moveRowSettings]
	styles      *styles
	data        *models.Data
}

func newMovesModel(s *styles, data *models.Data) movesModel {
	columns := []gridColumn{
		{col: table.Column{Title: MOVE_COL_NAME, Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: BASE_POWER_COL_NAME, Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: TYPE_COL_NAME, Width: 10}, editable: true, visible: true},
		{col: table.Column{Title: CATEGORY_COL_NAME, Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: CRIT_COL_NAME, Width: 10}, editable: true, visible: true},
		{col: table.Column{Title: HITS_COL_NAME, Width: 10}, editable: true, visible: true},
		{col: table.Column{Title: METR_HITS_COL_NAME, Width: 15}, editable: true, visible: false},
	}

	rows := [][]string{
		{DISPLAY_DEFAULT_MOVE1, DISPLAY_DEFAULT_BASEPOWER1, DISPLAY_DEFAULT_MOVE_TYPE1, DISPLAY_NONE, DISPLAY_NO, "0", "0"},
		{DISPLAY_DEFAULT_MOVE2, DISPLAY_DEFAULT_BASEPOWER2, DISPLAY_DEFAULT_MOVE_TYPE2, DISPLAY_NONE, DISPLAY_NO, "0", "0"},
		{DISPLAY_DEFAULT_MOVE3, DISPLAY_DEFAULT_BASEPOWER3, DISPLAY_DEFAULT_MOVE_TYPE3, DISPLAY_NONE, DISPLAY_NO, "0", "0"},
		{DISPLAY_DEFAULT_MOVE4, DISPLAY_DEFAULT_BASEPOWER4, DISPLAY_DEFAULT_MOVE_TYPE4, DISPLAY_NONE, DISPLAY_NO, "0", "0"},
	}

	initialColIndex := firstEditableCol(columns)

	rowData := make([]moveRowSettings, len(rows))

	return movesModel{
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
		rowSettings: side[[]moveRowSettings]{
			left:  append([]moveRowSettings{}, rowData...),
			right: append([]moveRowSettings{}, rowData...),
		},
		styles: s,
		data:   data,
	}
}

func (m movesModel) Init() tea.Cmd {
	return nil
}

func (m movesModel) Update(msg tea.Msg, side Side) (tabContent, tea.Cmd) {
	data := m.tab.GetPtr(side)
	rowSettings := m.rowSettings.GetPtr(side)

	switch msg := msg.(type) {
	case metronomeColVisibleRequestMsg:
		data.columns[METR_HITS_COL].visible = msg.visible
		data.SetNextValidCol()
		return m, nil

	case editConfirmMsg:
		data.Set(data.row, data.col, msg.value)

		if msg.value != "" {
			(*rowSettings)[data.row].isSet = true

			// automatically fill move row only when move name changes
			if msg.fieldID == MOVE_COL_NAME {
				m.fillMovesDefaults(data, msg.value)
			}
		} else {
			(*rowSettings)[data.row].isSet = false
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

func (m movesModel) View(side Side) string {
	return renderGrid(m.tab.Get(side), m.styles)
}

func (m movesModel) Names(side Side, settings []moveRowSettings) []string {
	var names []string
	data := m.tab.Get(side)
	for i, row := range data.rows {
		if !settings[i].isSet {
			// do not add placeholder moves to move slice
			continue
		}

		names = append(names, row[MOVE_COL])
	}

	return names
}

func (m movesModel) MoveRowSettings(side Side) []moveRowSettings {
	return m.rowSettings.Get(side)
}

func (m movesModel) requestInput(key, value string) tea.Cmd {
	min := 0
	max := 0
	var options []string
	isNumeric := false

	switch key {
	case MOVE_COL_NAME:
		options = m.data.MoveNames()
	case BASE_POWER_COL_NAME:
		isNumeric = true
		min = models.MOVE_BASE_POWER_MIN
		max = models.MOVE_BASE_POWER_MAX
	case TYPE_COL_NAME:
		options = m.data.TypeNames()
	case CATEGORY_COL_NAME:
		options = []string{DISPLAY_PHYSICAL_CATEGORY, DISPLAY_SPECIAL_CATEGORY, DISPLAY_STATUS_CATEGORY}
	case CRIT_COL_NAME:
		options = []string{DISPLAY_YES, DISPLAY_NO}
	case HITS_COL_NAME:
		options = hitLabels(hitOptions)
	case METR_HITS_COL_NAME:
		options = hitLabels(metronomeHitOptions)
	}

	if isNumeric {
		return requestNumberInput(key, value, min, max)
	} else {
		return requestSelectInput(key, value, options)
	}
}

func (m movesModel) fillMovesDefaults(data *gridData, moveName string) {
	moveData, ok := m.data.Moves[moveName]
	if !ok {
		// should not happen => move is validated through editConfirmMsg => if still unknown move here, dont do anything => internal error
		return
	}

	for i := range data.columns {
		if i == int(MOVE_COL) {
			// new move was already set in Update
			continue
		}

		var value string
		switch MoveColumnType(i) {
		case BP_COL:
			value = fmt.Sprintf("%d", moveData.BasePower)
		case MOVE_TYPE_COL:
			value = m.data.TypeName(moveData.Type)
		case CATEGORY_COL:
			value = moveData.Category
		case CRIT_COL:
			value = DISPLAY_NO
			if moveData.Critical {
				value = DISPLAY_YES
			}
		case HITS_COL:
			value = DEFAULT_HITS
			if len(moveData.Multihit) > 0 {
				value = hitLabel(hitOptions, slices.Min(moveData.Multihit))
			}
		default:
			continue
		}

		data.Set(data.row, i, value)
	}
}

func (m movesModel) Slots(side Side, settings []moveRowSettings) []string {
	data := m.tab.Get(side)

	slots := make([]string, len(data.rows))
	for i, row := range data.rows {
		if settings[i].isSet {
			slots[i] = row[MOVE_COL]
		}
	}

	return slots
}
