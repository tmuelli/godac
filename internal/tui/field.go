package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Setting int
type SideFocus int

const (
	GeneralSetting Setting = iota
	SidesSetting

	Left			SideFocus = iota
	Right

	CHECKED = "[x]"
	UNCHECKED = "[ ]"
)

type generalGroup struct {
	label		string
	settings	[]string
	index		int
	focused		bool

	styles		*styles
}

func (g generalGroup) Update(msg tea.Msg) (generalGroup, tea.Cmd) {
	if !g.focused {
		return g, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left":
			g.index = max(g.index - 1, 0)
		case "right":
			g.index = min(g.index + 1, len(g.settings) - 1)
		}
	}
	
	return g, nil
}

func (g generalGroup) View() string {
	var opts []string

	for i, s := range g.settings {
		style := g.styles.inactiveTab
		if i == g.index {
			style = g.styles.activeTab
		}
		opts = append(opts, style.Render(s))
	}

	labelStyle := lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.Border{Left: " ", Right: " "})
	if g.focused {
		labelStyle = labelStyle.Bold(true)
	}

	rowContent := append([]string{labelStyle.Render(g.label + ": ")}, opts...)
	row := lipgloss.JoinHorizontal(lipgloss.Top, rowContent...)
	return row
}

// TODO: use later for get easier value get
func (g generalGroup) Value() string {
	return g.settings[g.index]
}

type sidesGroup struct {
	label			string
	left			bool
	right			bool
	focused			SideFocus
	groupFocused	bool
	styles			*styles
}

func (g sidesGroup) Update(msg tea.Msg) (sidesGroup, tea.Cmd) {
	if !g.groupFocused {
		return g, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left":
			g.focused = Left
		case "right":
			g.focused = Right
		case "enter":
			if g.focused == Left {
				g.left = !g.left
			} else {
				g.right = !g.right
			}
		}
	}
	
	return g, nil
}

func (g sidesGroup) View() string {
	isGroupFocused := g.groupFocused

	// left
	leftCol := UNCHECKED
	if g.left {
		leftCol = CHECKED
	}

	leftStyle := lipgloss.NewStyle().Bold(isGroupFocused)
	if g.focused == Left {
		leftStyle = leftStyle.Foreground(lipgloss.Color("#EE82EE"))
	} else {
		leftStyle.Foreground(lipgloss.NoColor{})
	}

	// label
	lableStyle := lipgloss.NewStyle().Bold(isGroupFocused)

	// right
	rightCol := UNCHECKED
	if g.right {
		rightCol = CHECKED
	}

	rightStyle := lipgloss.NewStyle().Bold(isGroupFocused)
	if g.focused == Right {
		rightStyle = rightStyle.Foreground(lipgloss.Color("#EE82EE"))
	} else {
		rightStyle.Foreground(lipgloss.NoColor{})
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftCol),
		" ",
		lableStyle.Render(g.label),
		" ",
		rightStyle.Render(rightCol),
	)
}

type fieldModel struct {
	general		[]generalGroup
	generalRow	int
	sides		[]sidesGroup
	sidesRow	int
	setting		Setting
	styles		*styles
}

func newFieldStyles(bgIsDark bool) *styles {
	inactiveTabBorder := lipgloss.Border{
		Left: " ",
		Right: " ",
	}
	
	activeTabBorder := lipgloss.Border{
		Left: "[",
		Right: "]",
	}
	
	s := new(styles)
	s.doc = lipgloss.NewStyle()
	s.inactiveTab = lipgloss.NewStyle().Border(inactiveTabBorder, true).Padding(0, 1)
	s.activeTab = lipgloss.NewStyle().Border(activeTabBorder, true).Padding(0, 1)
	s.window = lipgloss.NewStyle()

	return s
}

func newGeneralGroup(s *styles, label string, defaultIndex int, settings []string) generalGroup {
	return generalGroup{
		label: label,
		settings: settings,
		index: defaultIndex,
		styles: s,
	}
}

func newSidesGroup(s *styles, label string, defaultFocus SideFocus, defaultLeft bool, defaultRight bool) sidesGroup {
	return sidesGroup{
		label: label,
		left: defaultLeft,
		right: defaultRight,
		focused: defaultFocus,
		styles: s,
	}
}

func newFieldModel(s *styles) fieldModel {
	fieldStyles := newFieldStyles(true)

	general := []generalGroup{
		newGeneralGroup(fieldStyles, "BattleMode", 0, []string{"Singles", "Doubles"}),
		newGeneralGroup(fieldStyles, "Level", 0, []string{"Level 100", "Level 50", "Level 5"}),
		newGeneralGroup(fieldStyles, "Weather", 0, []string{"None", "Sun", "Rain", "Sand", "Hail"}),
	}
	general[0].focused = true

	defaultSideFocus := Left
	sides := []sidesGroup{
		newSidesGroup(fieldStyles, "Stealth Rock", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Reflect", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Light Screen", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Protect", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Leech Seed", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Foresight", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Helping Hand", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Tailwind", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Flower Gift", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Power Trick", defaultSideFocus, false, false),
		newSidesGroup(fieldStyles, "Switching Out", defaultSideFocus, false, false),
	}
	sides[0].groupFocused = true

	return fieldModel{
		general: general,
		generalRow: 0,
		sides: sides,
		sidesRow: 0,
		setting: GeneralSetting,
		styles: fieldStyles,
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

func (m fieldModel) FormatGeneralSettings() string {
	var title string
	var lines []string

	titleStyle := lipgloss.NewStyle()
	if m.setting == GeneralSetting {
		titleStyle = titleStyle.Bold(true).Foreground(lipgloss.Color("#FFFF00"))
	} else {
		titleStyle.Foreground(lipgloss.NoColor{})
	}

	title = "General"

	for _, g := range m.general {
		lines = append(lines, g.View())
	}
	
	return titleStyle.Render(title) + ":\n" + lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m fieldModel) FormatSideSettings() string {
	var title string
	var lines []string

	titleStyle := lipgloss.NewStyle()
	if m.setting == SidesSetting {
		titleStyle = titleStyle.Bold(true).Foreground(lipgloss.Color("#FFFF00"))
	} else {
		titleStyle.Foreground(lipgloss.NoColor{})
	}

	title = "Sides"

	for _, s := range m.sides {
		lines = append(lines, s.View())
	}

	return "\n" + titleStyle.Render(title) + ":\n" + lipgloss.JoinVertical(lipgloss.Left, lines...)
}