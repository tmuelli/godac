package tui

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
	"github.com/tmuelli/godac/internal/utils"
)

type movesModel struct {
	tab			side[gridData]
	styles		*styles
	data		*models.Data
}

func newMovesModel(s *styles, data *models.Data) movesModel {
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
		tab:		side[gridData]{
			left:	gridData{
				columns:	columns,
				rows:		cloneRows(rows),
				row:		0,
				col:		0,
			},
			right:	gridData{
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

func (m movesModel) Init() tea.Cmd {
	return nil
}

func (m movesModel) Update(msg tea.Msg, side Side) (tabContent, tea.Cmd) {
	data := m.tab.GetPtr(side)
	
	switch msg := msg.(type) {
	case editConfirmMsg:
		data.Set(data.row, data.col, msg.value)
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
			key := data.columns[data.col].Title
			value := data.rows[data.row][data.col]
			return m, m.requestInput(key, value)
		}
	}
	
	return m, nil
}

func (m movesModel) View(side Side) string {
	return renderGrid(m.tab.Get(side), m.styles)
}

func (m movesModel) Names(side Side) []string {
	var names []string
	data := m.tab.Get(side)
	for _, row := range data.rows {
		names = append(names, row[0])
	}

	return names
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