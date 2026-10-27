package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
)

// TODO: implement current health in hp
const (
	NATURE         = "Nature"
	ABILITIY       = "Abiltiy"
	ITEM           = "Item"
	STATUS         = "Status"
	HEALTH_PERCENT = "Health Percent"
	CURRENT_HEALTH = "Current Health"
	METRONOME      = "Metronome"
)

var miscSpecs = []fieldSpec{
	{key: NATURE, options: func(d *models.Data) []string { return d.NatureNames() }},
	{key: ABILITIY, options: func(d *models.Data) []string { return d.AbilityNames() }},
	{
		key:     ITEM,
		options: func(d *models.Data) []string { return d.ItemNames() },
		onConfirm: func(data *list, value string, d *models.Data) tea.Cmd {
			isMetronomeColVisible := value == METRONOME
			return requestMetronomeColVisible(isMetronomeColVisible)
		},
	},
	{key: STATUS, options: func(d *models.Data) []string {
		return STATUS_DISPLAY_NAMES
	}},
	{key: HEALTH_PERCENT, isNumeric: true, min: 1, max: 100, defaultValue: fmt.Sprintf("%d", models.Level100)},
}

func newMiscModel(s *styles, data *models.Data) fieldListModel {
	model := newFieldListModel(s, data, miscSpecs)
	model.customUpdate = func(msg tea.Msg, data *list, d *models.Data) (tea.Cmd, bool) {
		if customMsg, ok := msg.(autoFillForSpecieRequestMsg); ok {
			fillMiscDefaults(data, customMsg.specie)
			return nil, true
		}

		return nil, false
	}

	return model
}

func fillMiscDefaults(data *list, specie *models.Specie) {
	if specie == nil {
		// internal error
		return
	}

	for k := range data.values {
		var value string
		switch k {
		case ABILITIY:
			value = specie.Abilities["0"]
		case STATUS:
			value = DISPLAY_STATUS_HEALTHY
		case HEALTH_PERCENT:
			value = fmt.Sprintf("%d%%", 100)
		case NATURE:
			value = DISPLAY_DEFAULT_NATURE
		}

		data.SetValue(k, value)
	}
}
