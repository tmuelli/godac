package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
	"github.com/tmuelli/godac/internal/utils"
	"golang.org/x/exp/maps"
)

type miscModel struct {
	fieldList	side[list]
	styles		*styles
	data		*models.Data
}

func newMiscModel(s *styles, data *models.Data) miscModel {
	template := map[string]string{
		"Nature": "",
		"Ability": "",
		"Item": "",
		"Status": "",
		"Health": "",
	}

	options := []string{"Nature", "Ability", "Item", "Status", "Health"}

	return miscModel{
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

func (m miscModel) Init() tea.Cmd {
	return nil
}

func (m miscModel) Update(msg tea.Msg, side Side) (tabContent, tea.Cmd) {
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

func (m miscModel) View(side Side) string {
	return renderList(m.fieldList.Get(side), m.styles)
}

// TODO: maybe make general requestInput
func (m miscModel) requestInput(key, value string) tea.Cmd {
	min := 0
	max := 0
	var options []string
	isNumeric := false

	// TODO: maybe use SliceFROmMapValue somewhere before to not do this everytime again
	// TODO: change to constants where its reasonable
	switch key {
	case "Nature":
		options = utils.SliceFromMapValue(m.data.Natures, func(n models.Nature) string {
			return n.Name
		})
	case "Ability":
		options = utils.SliceFromMapValue(m.data.Abilities, func(a models.Ability) string {
			return a.Name
		})
	case "Item":
		options = utils.SliceFromMapValue(m.data.Items, func(i models.Item) string {
			return i.Name
		})
	case "Status":
		options = []string{"Healthy", "Poisoned", "Badly Poisoned", "Burned", "Paralyzed", "Asleep", "Frozen"}
	case "Health":
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