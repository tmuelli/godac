package tui

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
	"github.com/tmuelli/godac/internal/utils"
)

type MoveColumnType int

const (
	MOVE_COL		MoveColumnType = iota
	BP_COL
	MOVE_TYPE_COL
	CATEGORY_COL
	CRIT_COL
	HITS_COL
)

type moveRowSettings struct {
	isSet	bool
}

type movesModel struct {
	tab				side[gridData]
	rowSettings		side[[]moveRowSettings]
	styles			*styles
	data			*models.Data
}

func newMovesModel(s *styles, data *models.Data) movesModel {
	columns := []gridColumn{
		{col: table.Column{Title: "Move", Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: "Base Power", Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: "Type", Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: "Category", Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: "Crit", Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: "Hits", Width: 15}, editable: true, visible: true},
	}

	rows := [][]string{
		{"Move1", "BP1", "Type1", "None", "No", "0"},
		{"Move2", "BP2", "Type2", "None", "No", "0"},
		{"Move3", "BP3", "Type3", "None", "No", "0"},
		{"Move4", "BP4", "Type4", "None", "No", "0"},
	}

	initialColIndex := firstEditableCol(columns)

	rowData := make([]moveRowSettings, len(rows))

	return movesModel{
		tab:		side[gridData]{
			left:	gridData{
				columns:	columns,
				rows:		cloneRows(rows),
				row:		0,
				col:		initialColIndex,
			},
			right:	gridData{
				columns:	columns,
				rows:		cloneRows(rows),
				row:		0,
				col:		initialColIndex,
			},
		},
		rowSettings: side[[]moveRowSettings]{
			left:	append([]moveRowSettings{}, rowData...),
			right:	append([]moveRowSettings{}, rowData...),
		},
		styles:		s,
		data:		data,
	}
}

func (m movesModel) Init() tea.Cmd {
	return nil
}

func (m movesModel) Update(msg tea.Msg, side Side) (tabContent, tea.Cmd) {
	data := m.tab.GetPtr(side)
	rowSettings := m.rowSettings.GetPtr(side)
	
	switch msg := msg.(type) {
	case editConfirmMsg:
		data.Set(data.row, data.col, msg.value)

		if msg.value != "" {
			(*rowSettings)[data.row].isSet = true

			// automatically fill move row only when move name changes
			if msg.fieldID == "Move" {
				m.fillMovesDefaults(data, msg.value)
			}
		} else {
			(*rowSettings)[data.row].isSet = false
		}

		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			data.MoveRow(-1)
		case "down":
			data.MoveRow(1)
		case "left":
			data.MoveCol(-1)
		case "right":
			data.MoveCol(1)
		case "enter":
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

		names = append(names, row[0])
	}

	return names
}

func (m movesModel) MoveRowSettings(side Side) []moveRowSettings {
	return m.rowSettings.Get(side)
}

// TODO: maybe make general requestInput
func (m movesModel) requestInput(key, value string) tea.Cmd {
	min := 0
	max := 0
	var options []string
	isNumeric := false

	// TODO: maybe use SliceFROmMapValue somewhere before to not do this everytime again
	// TODO: change to constants where its reasonable
	switch key {
	case "Move":
		options = utils.SliceFromMapValue(m.data.Moves, func(m models.Move) string {
			return m.Name
		})
	case "Base Power":
		isNumeric = true
		min = 0
		// TODO: numeric input with only one limit => set sane limit?
		max = 1000
	case "Type":
		options = utils.SliceFromMapValue(m.data.Types, func(t models.Type) string {
			return t.Name
		})
	case "Category":
		options = []string{"Physical", "Special"}
	case "Crit":
		options = []string{"Yes", "No"}
	case "Hits":
		// TODO: better solution
		options = []string{"Once", "Twice", "3 times", "4 times", "5 times",
			"1 hits", "2 hits", "3 hits", "4 hits", "5 hits",
		}
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
			value = moveData.Type
		case CATEGORY_COL:
			value = moveData.Category
		case HITS_COL:
			value = DEFAULT_HITS
		default:
			continue
		}

		data.Set(data.row, i, value)
	}
}