package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
)

type Setting int

const (
	GeneralSetting Setting = iota
	SidesSetting
)

type fieldModel struct {
	general		[]generalGroup
	generalRow	int
	sides		[]sidesGroup
	sidesRow	int
	setting		Setting
	styles		*styles
	data		*models.Data
}

// TODO: implement value getter for spikes
func (m fieldModel) SideGroupValue(label string, side SideFocus) bool {
	value := false
	
	for _, s := range m.sides {
		if s.label == label {
			value = s.Value(side)
		}
	}

	return value
}

func newFieldModel(s *styles, data *models.Data) fieldModel {
	general := []generalGroup{
		newGeneralGroup(s, "BattleMode", 0, []string{"Singles", "Doubles"}),
		newGeneralGroup(s, "Level", 0, []string{"Level 100", "Level 50", "Level 5"}),
		newGeneralGroup(s, "Weather", 0, []string{"None", "Sun", "Rain", "Sand", "Hail"}),
	}
	general[0].focused = true

	// TODO: add spikes settings
	defaultSideFocus := Left
	sides := []sidesGroup{
		newSidesGroup(s, "Stealth Rock", defaultSideFocus, false, false),
		newSidesGroup(s, "Reflect", defaultSideFocus, false, false),
		newSidesGroup(s, "Light Screen", defaultSideFocus, false, false),
		newSidesGroup(s, "Protect", defaultSideFocus, false, false),
		newSidesGroup(s, "Leech Seed", defaultSideFocus, false, false),
		newSidesGroup(s, "Foresight", defaultSideFocus, false, false),
		newSidesGroup(s, "Helping Hand", defaultSideFocus, false, false),
		newSidesGroup(s, "Tailwind", defaultSideFocus, false, false),
		newSidesGroup(s, "Flower Gift", defaultSideFocus, false, false),
		newSidesGroup(s, "Power Trick", defaultSideFocus, false, false),
		newSidesGroup(s, "Switching Out", defaultSideFocus, false, false),
	}
	sides[0].groupFocused = true

	return fieldModel{
		general: general,
		generalRow: 0,
		sides: sides,
		sidesRow: 0,
		setting: GeneralSetting,
		styles: s,
		data:	data,
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
			switch msg.String() {
			case "s":
				m.setting = SidesSetting
				return m, nil
			case "up":
				m.general[m.generalRow].focused = false
				m.generalRow = max(m.generalRow - 1, 0)
				m.general[m.generalRow].focused = true
				return m, nil
			case "down":
				m.general[m.generalRow].focused = false
				m.generalRow = min(m.generalRow + 1, len(m.general) - 1)
				m.general[m.generalRow].focused = true
				return m, nil
			}
		case SidesSetting:
			switch msg.String() {
			case "g":
				m.setting = GeneralSetting
				return m, nil
			case "up":
				m.sides[m.sidesRow].groupFocused = false
				m.sidesRow = max(m.sidesRow - 1, 0)
				m.sides[m.sidesRow].groupFocused = true
				return m, nil
			case "down":
				m.sides[m.sidesRow].groupFocused = false
				m.sidesRow = min(m.sidesRow + 1, len(m.sides) - 1)
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
	return "Field view - Press [C] for calculation view or [R] or [L] for config view"
}

func (m fieldModel) Body() string {
	content := m.FormatGeneralSettings()
	content += m.FormatSideSettings()
	return content
}

func (m fieldModel) SideInput(side SideFocus) *models.UserFieldSideInput {
	input := models.UserFieldSideInput{
		StealthRocks:	m.SideGroupValue("Stealth Rock", side),
		// TODO: implement value getter for spikes
		Spikes:			0,
		Reflect:		m.SideGroupValue("Reflect", side),
		LightScreen:	m.SideGroupValue("Light Screen", side),
		Protect:		m.SideGroupValue("Protect", side),
		LeechSeed:		m.SideGroupValue("Leech Seed", side),
		Foresight:		m.SideGroupValue("Foresight", side),
		HelpingHand: 	m.SideGroupValue("Helping Hand", side),
		Tailwind:		m.SideGroupValue("Tailwind", side),
		FlowerGift:		m.SideGroupValue("Flower Gift", side),
		PowerTrick:		m.SideGroupValue("Power Trick", side),
		SwitchingOut:	m.SideGroupValue("Switching Out", side),
	}
	
	return &input
}

func (m fieldModel) Input() *models.UserFieldInput {
	// type
	fieldType := m.general[0].Value()

	// level
	var level int
	switch m.general[1].Value() {
	case "Level 100":
		level = 100
	case "Level 50":
		level = 50
	case "Level 5":
		level = 5
	}

	// terrain
	terrain := m.general[2].Value()

	fieldInput := models.UserFieldInput{
		Type:		fieldType,
		Level:		level,
		Terrain:	terrain,
		IsGravity:	false,	// TODO implement later in field view
		Left:		*m.SideInput(Left),
		Right:		*m.SideInput(Right),
	}
	
	return &fieldInput
}