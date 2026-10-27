package tui

import (
	"fmt"
	"math"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
)

const (
	NATURE         = "Nature"
	ABILITIY       = "Abiltiy"
	ITEM           = "Item"
	STATUS         = "Status"
	HEALTH         = "Health"
	HEALTH_PERCENT = "Health Percent"
	METRONOME      = "Metronome"
)

type miscModel struct {
	fieldListModel
	maxHP side[int]
}

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
	{key: HEALTH, isNumeric: true},
}

func newMiscModel(s *styles, data *models.Data) miscModel {
	model := newFieldListModel(s, data, miscSpecs)
	model.customUpdate = func(msg tea.Msg, data *list, d *models.Data) (tea.Cmd, bool) {
		if customMsg, ok := msg.(autoFillForSpecieRequestMsg); ok {
			fillMiscDefaults(data, customMsg.specie)
			return nil, true
		}
		return nil, false
	}

	return miscModel{fieldListModel: model}
}

func (m miscModel) Update(msg tea.Msg, side Side) (tabContent, tea.Cmd) {
	// confirm on health row => set max for validation
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok && key.Matches(keyMsg, globalKeys.Confirm) {
		list := m.fields.GetPtr(side)
		if list.options[list.activeIndex] == HEALTH {
			maxHP := m.maxHP.Get(side)
			if maxHP == 0 {
				// no specie => no max hp
				return m, nil
			}

			return m, requestNumberInput(HEALTH, list.values[HEALTH], 1, maxHP)
		}
	}

	// default fieldListModel Update
	base, cmd := m.fieldListModel.Update(msg, side)
	m.fieldListModel = base.(fieldListModel)

	// update health fields
	if edit, ok := msg.(editConfirmMsg); ok {
		list := m.fields.GetPtr(side)
		maxHP := m.maxHP.Get(side)

		switch edit.fieldID {
		case HEALTH_PERCENT:
			perc := atoiOrZero(edit.value)
			list.SetValue(HEALTH_PERCENT, fmt.Sprintf("%d%%", perc))
			list.SetValue(HEALTH, fmt.Sprintf("%d", percentToHealth(perc, maxHP)))
		case HEALTH:
			hp := atoiOrZero(edit.value)
			list.SetValue(HEALTH_PERCENT, fmt.Sprintf("%d%%", healthToPercent(hp, maxHP)))
		}
	}

	return m, cmd
}

func (m *miscModel) setMaxHP(side Side, maxHP int) {
	m.maxHP.Set(side, maxHP)

	list := m.fields.GetPtr(side)
	percStr := strings.TrimSuffix(list.values[HEALTH_PERCENT], "%")
	if percStr == "" {
		// not filled
		return
	}

	list.SetValue(HEALTH, fmt.Sprintf("%d", percentToHealth(atoiOrZero(percStr), maxHP)))
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

func percentToHealth(perc, maxHP int) int {
	return int(math.Round(float64(perc) * float64(maxHP) / 100))
}

func healthToPercent(hp, maxHP int) int {
	return max(1, int(math.Ceil(float64(hp)*100/float64(maxHP))))
}
