package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
)

var fieldLevels = map[string]int{
	DISPLAY_LEVEL_100: int(models.Level100),
	DISPLAY_LEVEL_50:  int(models.Level50),
	DISPLAY_LEVEL_5:   int(models.Level5),
}

type Setting int
type GeneralRow int

const (
	GeneralSetting Setting = iota
	SidesSetting

	DISPLAY_GENERAL_GROUP_LEVEL   = "Level"
	DISPLAY_GENERAL_GROUP_MODE    = "BattleMode"
	DISPLAY_GENERAL_GROUP_WEATHER = "Weather"
	DISPLAY_GENERAL_GROUP_GRAVITY = "Gravity"
)

const (
	FieldStealthRock = iota
	FieldSpikes
	FieldReflect
	FieldLightScreen
	FieldProtect
	FieldLeechSeed
	FieldForesight
	FieldHelpingHand
	FieldTailwind
	FieldFlowerGift
	FieldPowerTrick
	FieldSwichtingOut

	DISPLAY_STEALTH_ROCK  = "Stealth Rock"
	DISPLAY_SPIKES        = "Spikes"
	DISPLAY_REFLECT       = "Reflect"
	DISPLAY_LIGHT_SCREEN  = "Light Screen"
	DISPLAY_PROTECT       = "Protect"
	DISPLAY_LEECH_SEED    = "Leech Seed"
	DISPLAY_FORESIGHT     = "Foresight"
	DISPLAY_HELPING_HAND  = "Helping Hand"
	DISPLAY_TAILWIND      = "Tailwind"
	DISPLAY_FLOWER_GIFT   = "Flower Gift"
	DISPLAY_POWER_TRICK   = "Power Trick"
	DISPLAY_SWITCHING_OUT = "Switching Out"
)

const (
	FieldBattleMode GeneralRow = iota
	FieldLevel
	FieldWeather
	FieldGravity
)

type fieldKeyMap struct {
	General key.Binding
	Side    key.Binding
}

var fieldKeys = fieldKeyMap{
	General: key.NewBinding(key.WithKeys("g"), key.WithHelp("general", "general settings")),
	Side:    key.NewBinding(key.WithKeys("s"), key.WithHelp("side", "field side settings")),
}

func (k fieldKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k fieldKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.General, k.Side},
	}
}

type fieldModel struct {
	general    []generalGroup
	generalRow int
	sides      []sidesGroup
	sidesRow   int
	setting    Setting
	styles     *styles
	data       *models.Data
}

func (m fieldModel) SideGroupLevel(label string, side Side) int {
	level := 0

	for _, s := range m.sides {
		if s.label == label {
			level = s.Level(side)
		}
	}

	return level
}

func (m fieldModel) SideGroupValue(label string, side Side) bool {
	value := false

	for _, s := range m.sides {
		if s.label == label {
			value = s.Value(side)
		}
	}

	return value
}

func newFieldModel(s *styles, data *models.Data) fieldModel {
	levelGroup := newGeneralGroup(s, DISPLAY_GENERAL_GROUP_LEVEL, 0, []string{DISPLAY_LEVEL_100, DISPLAY_LEVEL_50, DISPLAY_LEVEL_5})
	levelGroup.onChange = func(value string) tea.Cmd {
		return requestSetLevel(fieldLevels[value])
	}

	general := []generalGroup{
		newGeneralGroup(s, DISPLAY_GENERAL_GROUP_MODE, 0, []string{string(models.Singles), string(models.Doubles)}),
		levelGroup,
		newGeneralGroup(s, DISPLAY_GENERAL_GROUP_WEATHER, 0, []string{
			string(models.Clear),
			string(models.Sun),
			string(models.Rain),
			string(models.Sand),
			string(models.Hail),
		}),
		newGeneralGroup(s, DISPLAY_GENERAL_GROUP_GRAVITY, 0, []string{DISPLAY_OFF, DISPLAY_ON}),
	}
	general[FieldBattleMode].focused = true

	defaultSide := LeftSide
	sides := []sidesGroup{
		newSidesGroup(s, DISPLAY_STEALTH_ROCK, defaultSide, false, false),
		newSpikesGroup(s, defaultSide),
		newSidesGroup(s, DISPLAY_REFLECT, defaultSide, false, false),
		newSidesGroup(s, DISPLAY_LIGHT_SCREEN, defaultSide, false, false),
		newSidesGroup(s, DISPLAY_PROTECT, defaultSide, false, false),
		newSidesGroup(s, DISPLAY_LEECH_SEED, defaultSide, false, false),
		newSidesGroup(s, DISPLAY_FORESIGHT, defaultSide, false, false),
		newSidesGroup(s, DISPLAY_HELPING_HAND, defaultSide, false, false),
		newSidesGroup(s, DISPLAY_TAILWIND, defaultSide, false, false),
		newSidesGroup(s, DISPLAY_FLOWER_GIFT, defaultSide, false, false),
		newSidesGroup(s, DISPLAY_POWER_TRICK, defaultSide, false, false),
		newSidesGroup(s, DISPLAY_SWITCHING_OUT, defaultSide, false, false),
	}
	sides[FieldStealthRock].groupFocused = true

	return fieldModel{
		general:    general,
		generalRow: 0,
		sides:      sides,
		sidesRow:   0,
		setting:    GeneralSetting,
		styles:     s,
		data:       data,
	}
}

func (m fieldModel) Init() tea.Cmd {
	return nil
}

func (m fieldModel) Update(msg tea.Msg) (fieldModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch m.setting {
		case GeneralSetting:
			switch {
			case key.Matches(msg, fieldKeys.Side):
				m.setting = SidesSetting
				return m, nil
			case key.Matches(msg, globalKeys.Up):
				m.general[m.generalRow].focused = false
				m.generalRow = max(m.generalRow-1, 0)
				m.general[m.generalRow].focused = true
				return m, nil
			case key.Matches(msg, globalKeys.Down):
				m.general[m.generalRow].focused = false
				m.generalRow = min(m.generalRow+1, len(m.general)-1)
				m.general[m.generalRow].focused = true
				return m, nil
			}
		case SidesSetting:
			switch {
			case key.Matches(msg, fieldKeys.General):
				m.setting = GeneralSetting
				return m, nil
			case key.Matches(msg, globalKeys.Up):
				m.sides[m.sidesRow].groupFocused = false
				m.sidesRow = max(m.sidesRow-1, 0)
				m.sides[m.sidesRow].groupFocused = true
				return m, nil
			case key.Matches(msg, globalKeys.Down):
				m.sides[m.sidesRow].groupFocused = false
				m.sidesRow = min(m.sidesRow+1, len(m.sides)-1)
				m.sides[m.sidesRow].groupFocused = true
				return m, nil
			}
		}
	}

	switch m.setting {
	case GeneralSetting:
		var cmd tea.Cmd
		m.general[m.generalRow], cmd = m.general[m.generalRow].Update(msg)
		return m, cmd

	case SidesSetting:
		var cmd tea.Cmd
		m.sides[m.sidesRow], cmd = m.sides[m.sidesRow].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m fieldModel) View() string {
	content := m.Title()
	content += "\n\n"
	content += m.Body()
	return content
}

func (m fieldModel) Title() string {
	return FIELD_TITLE
}

func (m fieldModel) Body() string {
	content := m.FormatGeneralSettings()
	content += m.FormatSideSettings()
	return content
}

func (m fieldModel) SideInput(side Side) *models.UserFieldSideInput {
	input := models.UserFieldSideInput{
		StealthRocks: m.SideGroupValue(DISPLAY_STEALTH_ROCK, side),
		Spikes:       m.SideGroupLevel(DISPLAY_SPIKES, side),
		Reflect:      m.SideGroupValue(DISPLAY_REFLECT, side),
		LightScreen:  m.SideGroupValue(DISPLAY_LIGHT_SCREEN, side),
		Protect:      m.SideGroupValue(DISPLAY_PROTECT, side),
		LeechSeed:    m.SideGroupValue(DISPLAY_LEECH_SEED, side),
		Foresight:    m.SideGroupValue(DISPLAY_FORESIGHT, side),
		HelpingHand:  m.SideGroupValue(DISPLAY_HELPING_HAND, side),
		Tailwind:     m.SideGroupValue(DISPLAY_TAILWIND, side),
		FlowerGift:   m.SideGroupValue(DISPLAY_FLOWER_GIFT, side),
		PowerTrick:   m.SideGroupValue(DISPLAY_POWER_TRICK, side),
		SwitchingOut: m.SideGroupValue(DISPLAY_SWITCHING_OUT, side),
	}

	return &input
}

func (m fieldModel) Input() *models.UserFieldInput {
	// type
	fieldType := m.general[FieldBattleMode].Value()

	// level
	level := fieldLevels[m.general[FieldLevel].Value()]

	// terrain
	terrain := m.general[FieldWeather].Value()

	// gravity
	var gravity bool
	gravitySetting := m.general[FieldGravity].Value()
	if gravitySetting == DISPLAY_ON {
		gravity = true
	}

	fieldInput := models.UserFieldInput{
		Type:      fieldType,
		Level:     level,
		Terrain:   terrain,
		IsGravity: gravity,
		Left:      *m.SideInput(LeftSide),
		Right:     *m.SideInput(RightSide),
	}

	return &fieldInput
}
