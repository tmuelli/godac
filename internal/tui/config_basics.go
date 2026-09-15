package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
	"github.com/tmuelli/godac/internal/utils"
	"golang.org/x/exp/maps"
)

type basicsModel struct {
	fieldList	side[list]
	styles		*styles
	data		*models.Data
}

func newBasicsModel(s *styles, data *models.Data) basicsModel {
	template := map[string]string{
		"Name": "",
		"Gender": "",
		"Type1": "",
		"Type2": "",
		"Level": "",
	}

	options := []string{"Name", "Gender", "Type1", "Type2", "Level"}

	return basicsModel{
		fieldList:	side[list]{
			left:	list{
				values:			maps.Clone(template),
				options:		options,
				activeIndex:	0,
			},
			right:	list{
				values:			maps.Clone(template),
				options:		options,
				activeIndex:	0,
			},
		},
		styles:			s,
		data:			data,
	}
}

func (m basicsModel) Init() tea.Cmd {
	return nil
}

func (m basicsModel) Update(msg tea.Msg, side Side) (tabContent, tea.Cmd) {
	data := m.fieldList.GetPtr(side)
	
	switch msg := msg.(type) {
	case editConfirmMsg:
		data.SetValue(msg.fieldID, msg.value)
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			data.SetActiveRow(-1)
		case "down":
			data.SetActiveRow(1)
		case "enter":
			key := data.options[data.activeIndex]
			value := data.values[key]
			return m, m.requestInput(key, value)
		}
	}

	return m, nil
}

func (m basicsModel) View(side Side) string {
	return renderList(m.fieldList.Get(side), m.styles)
}

// TODO: maybe make general requestInput
func (m basicsModel) requestInput(key, value string) tea.Cmd {
	min := 0
	max := 0
	var options []string
	isNumeric := false

	// TODO: maybe use SliceFROmMapValue somewhere before to not do this everytime again
	// TODO: change to constants where its reasonable
	switch key {
	case "Name":
		options = utils.SliceFromMapValue(m.data.Species, func(s models.Specie) string {
			return s.Name
		})
	case "Gender":
		options = []string{"Female", "Male", "Genderless"}
	case "Type1", "Type2":
		options = utils.SliceFromMapValue(m.data.Types, func(t models.Type) string {
			return t.Name
		})
	case "Level":
		isNumeric = true
		min = 1
		max = 100
	}

	if isNumeric {
		return requestNumberInput(key, value, min, max)
	} else {
		return requestSelectInput(key, value, options)
	}
}