package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
)

type tabContent interface {
	Init()									tea.Cmd
	Update(msg tea.Msg, side Side)	(tabContent, tea.Cmd)
	View(side Side)					string
}

type configModel struct {
	tabs		[]string
	tabContent	[]tabContent
	activeTab	int
	styles		*styles
	data		*models.Data
}

func newConfigModel(s *styles, data *models.Data) configModel {
	tabs := []string{"Basics", "Stats", "Misc.", "Moves"}
	tabContent := []tabContent{
		newBasicsModel(s, data),
		newStatsModel(s, data),
		newMiscModel(s, data),
		newMovesModel(s, data),
	}

	return configModel{
		tabs:		tabs,
		tabContent:	tabContent,
		activeTab:	0,
		styles:		s,
		data:		data,
	}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg, side Side) (configModel, tea.Cmd) {
	switch msg := msg.(type) {
	case autoFillForSpecieRequestMsg:
		var cmds []tea.Cmd
		var cmd tea.Cmd

		m.tabContent[statsTab], cmd = m.tabContent[statsTab].Update(msg, side)
		cmds = append(cmds, cmd)

		m.tabContent[miscTab], cmd = m.tabContent[miscTab].Update(msg, side)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil

		case "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.tabContent[m.activeTab], cmd = m.tabContent[m.activeTab].Update(msg, side)
	return m, cmd
}

func (m configModel) View(side Side) string {
	content := m.Title()
	content += "\n\n"
	content += m.Body(side)
	return content
}

func (m configModel) Title() string {
	return "Config view - Press [C] for calculation view or [F] for field view"
}

func (m configModel) Body(side Side) string {
	if m.styles == nil {
		return ""
	}
	
	doc := strings.Builder{}
	s := m.styles

	var renderedTabs []string

	for i, t := range m.tabs {
		style := m.styles.tab
		if i == m.activeTab {
			style = style.BorderForeground(m.styles.highlight).Bold(true)
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	content := m.tabContent[m.activeTab].View(side)
	doc.WriteString(s.window.Width((lipgloss.Width(content))).Render(content))
	return s.doc.Render(doc.String())
}

func (m configModel) Input(side Side, reqInfo *calcRequestInfo) *models.UserInput {
	attackMove := ""
	if side == Side(reqInfo.usedSide) {
		attackMove = reqInfo.usedMove
	}

	return &models.UserInput{
		Name:				m.Name(side),
		Level:				m.Level(side),
		Gender:				m.Gender(side),
		AttackMove:			attackMove,
		Moves:				m.Moves(side),
		HealthPercent:		m.HealthPercent(side),
		// TODO: implement shortcut iv settings (optimal, worstcase) input option
		IvSetting:			"custom",
		IvCustom:			m.IvCustomStats(side),
		// TODO: implement shortcut ev settings (optimal, worstcase) input option
		EvSetting:			"custom",
		EvCustom:			m.EvCustomStats(side),
		Item:				m.Item(side),
		Status:				m.Status(side),
		Hits:				m.Hits(side, attackMove),
		Boosts:				m.Boosts(side),
	}
}
