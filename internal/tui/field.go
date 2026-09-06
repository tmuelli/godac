package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

type fieldModel struct {
	general		[]generalGroup
	generalRow	int
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

func newFieldModel(s *styles) fieldModel {
	fieldStyles := newFieldStyles(true)

	general := []generalGroup{
		newGeneralGroup(fieldStyles, "BattleMode", 0, []string{"Singles", "Doubles"}),
		newGeneralGroup(fieldStyles, "Level", 0, []string{"Level 100", "Level 50", "Level 5"}),
		newGeneralGroup(fieldStyles, "Weather", 0, []string{"None", "Sun", "Rain", "Sand", "Hail"}),
	}
	general[0].focused = true

	return fieldModel{
		general: general,
		generalRow: 0,
	}
}

func (m fieldModel) Init() tea.Cmd {
	return nil
}

func (m fieldModel) Update(msg tea.Msg) (fieldModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
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
	}
	
	var cmd tea.Cmd
	m.general[m.generalRow], cmd = m.general[m.generalRow].Update(msg)
	return m, cmd
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
	content := "General:\n"
	content += m.FormatGeneralSettings()
	return content
}

func (m fieldModel) FormatGeneralSettings() string {
	var lines []string

	for _, g := range m.general {
		lines = append(lines, g.View())
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}