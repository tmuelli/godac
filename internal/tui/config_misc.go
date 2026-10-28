package tui

import (
	"fmt"
	"image/color"
	"math"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
)

const (
	NATURE         = "Nature"
	ABILITY        = "Ability"
	ABILITY_ON     = "Ability On"
	ITEM           = "Item"
	STATUS         = "Status"
	HEALTH         = "Health"
	HEALTH_PERCENT = "Health Percent"
	METRONOME      = "Metronome"

	// colors
	HP_BAR_WIDTH    = 30
	HP_GREEN_ABOVE  = 0.5
	HP_YELLOW_ABOVE = 0.2

	COLOR_HP_BAR_GREEN  = "#00dd60"
	COLOR_HP_BAR_YELLOW = "#f5d538"
	COLOR_HP_BAR_RED    = "#ee4928"
	COLOR_HP_BAR_EMPTY  = "#444444"
)

func hpBarColor(total, _ float64) color.Color {
	switch {
	case total > HP_GREEN_ABOVE:
		return lipgloss.Color(COLOR_HP_BAR_GREEN)
	case total > HP_YELLOW_ABOVE:
		return lipgloss.Color(COLOR_HP_BAR_YELLOW)
	default:
		return lipgloss.Color(COLOR_HP_BAR_RED)
	}
}

type miscModel struct {
	fieldListModel
	maxHP side[int]
	hpBar progress.Model
}

var miscSpecs = []fieldSpec{
	{key: NATURE, options: func(d *models.Data) []string { return d.NatureNames() }},
	{
		key:     ABILITY,
		options: func(d *models.Data) []string { return d.AbilityNames() },
		onConfirm: func(data *list, value string, d *models.Data) tea.Cmd {
			setAbilityOn(data, d, value)
			return nil
		},
	},
	{key: ABILITY_ON, isToggle: true},
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
	bar := progress.New(
		progress.WithColorFunc(hpBarColor),
		progress.WithWidth(HP_BAR_WIDTH),
		progress.WithoutPercentage(),
		progress.WithFillCharacters('█', '░'),
	)
	bar.EmptyColor = lipgloss.Color(COLOR_HP_BAR_EMPTY)

	model := newFieldListModel(s, data, miscSpecs)
	model.customUpdate = func(msg tea.Msg, data *list, d *models.Data) (tea.Cmd, bool) {
		if customMsg, ok := msg.(autoFillForSpecieRequestMsg); ok {
			fillMiscDefaults(data, customMsg.specie, d)
			return nil, true
		}
		return nil, false
	}

	return miscModel{fieldListModel: model, hpBar: bar}
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

func (m miscModel) View(side Side) string {
	listView := m.fieldListModel.View(side)
	perc := healthPercent(m.fields.Get(side))
	return listView + "\n" + m.hpBar.ViewAs(float64(perc)/100)
}

func (m *miscModel) setMaxHP(side Side, maxHP int) {
	m.maxHP.Set(side, maxHP)

	list := m.fields.GetPtr(side)
	if list.values[HEALTH_PERCENT] == "" {
		// not filled
		return
	}

	list.SetValue(HEALTH, fmt.Sprintf("%d", percentToHealth(healthPercent(*list), maxHP)))
}

func fillMiscDefaults(data *list, specie *models.Specie, d *models.Data) {
	if specie == nil {
		// internal error
		return
	}

	for k := range data.values {
		var value string
		switch k {
		case ABILITY:
			value = specie.Abilities["0"]
		case ABILITY_ON:
			setAbilityOn(data, d, specie.Abilities["0"])
			continue
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

// default display value for toggle abilities => empty for abilities without toggle
func abilityOnValue(d *models.Data, abilityName string) string {
	ability, ok := d.Abilities[abilityName]
	if !ok {
		return ""
	}

	defaultOn, hasToggle := models.TogglableAbilities[ability.Id]
	if !hasToggle {
		return ""
	}

	if defaultOn {
		return DISPLAY_YES
	}

	return DISPLAY_NO
}

// reset toggle row to abilities default
func setAbilityOn(data *list, d *models.Data, abilityName string) {
	value := abilityOnValue(d, abilityName)
	data.SetValue(ABILITY_ON, value)
	data.SetHidden(ABILITY_ON, value == "")
}
