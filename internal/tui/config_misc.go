package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
	"github.com/tmuelli/godac/internal/utils"
	"golang.org/x/exp/maps"
)

// TODO: implement current health in hp
const (
	NATURE = "Nature"
	ABILTIY = "Abiltiy"
	ITEM = "Item"
	STATUS = "Status"
	HEALTH_PERCENT = "Health Percent"
	CURRENT_HEALTH = "Current Health"
)

var KEYS = []string{
	NATURE,
	ABILTIY,
	ITEM,
	STATUS,
	HEALTH_PERCENT,
	CURRENT_HEALTH,
}

type miscModel struct {
	fieldList	side[list]
	styles		*styles
	data		*models.Data
}

func newMiscModel(s *styles, data *models.Data) miscModel {
	template := make(map[string]string, len(KEYS))
	for _, k := range KEYS {
		template[k] = ""
	}

	options := KEYS

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
	case autoFillForSpecieRequestMsg:
		m.fillMiscDefaults(data, msg.specie)
		return m, nil

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
	case NATURE:
		options = utils.SliceFromMapValue(m.data.Natures, func(n models.Nature) string {
			return n.Name
		})
	case ABILTIY:
		options = utils.SliceFromMapValue(m.data.Abilities, func(a models.Ability) string {
			return a.Name
		})
	case ITEM:
		options = utils.SliceFromMapValue(m.data.Items, func(i models.Item) string {
			return i.Name
		})
	case STATUS:
		options = []string{"Healthy", "Poisoned", "Badly Poisoned", "Burned", "Paralyzed", "Asleep", "Frozen"}
	case HEALTH_PERCENT:
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

func (m miscModel) fillMiscDefaults(data *list, specie *models.Specie) {
	if specie == nil {
		// internal error
		return
	}

	for k := range data.values {
		var value string
		switch k {
		case ABILTIY:
			value = specie.Abilities["0"]
		case STATUS:
			value = "healthy"
		case HEALTH_PERCENT:
			value = fmt.Sprintf("%d%%", 100)
		case NATURE:
			value = DEFAULT_NATURE
		}

		data.SetValue(k, value)
	}
}