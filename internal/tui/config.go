package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
)

const (
	DISPLAY_TAB_NAME_BASICS = "Basics"
	DISPLAY_TAB_NAME_STATS  = "Stats"
	DISPLAY_TAB_NAME_MISC   = "Misc."
	DISPLAY_TAB_NAME_MOVES  = "Moves"
)

type configKeyMap struct {
	NextTab key.Binding
	PrevTab key.Binding
}

var configKeys = configKeyMap{
	NextTab: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
	PrevTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
}

func (k configKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k configKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextTab, k.PrevTab},
	}
}

type tabContent interface {
	Init() tea.Cmd
	Update(msg tea.Msg, side Side) (tabContent, tea.Cmd)
	View(side Side) string
}

type configModel struct {
	tabs       []string
	tabContent []tabContent
	activeTab  int
	styles     *styles
	data       *models.Data
}

func newConfigModel(s *styles, data *models.Data) configModel {
	tabs := []string{DISPLAY_TAB_NAME_BASICS, DISPLAY_TAB_NAME_STATS, DISPLAY_TAB_NAME_MISC, DISPLAY_TAB_NAME_MOVES}
	tabContent := []tabContent{
		newBasicsModel(s, data),
		newStatsModel(s, data),
		newMiscModel(s, data),
		newMovesModel(s, data),
	}

	return configModel{
		tabs:       tabs,
		tabContent: tabContent,
		activeTab:  0,
		styles:     s,
		data:       data,
	}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg, side Side) (configModel, tea.Cmd) {
	switch msg := msg.(type) {
	case hpInputsChangedMsg:
		m.refreshMaxHP(side)
		return m, nil

	case setLevelRequestMsg:
		var cmds []tea.Cmd

		for _, s := range []Side{LeftSide, RightSide} {
			var cmd tea.Cmd
			m.tabContent[basicsTab], cmd = m.tabContent[basicsTab].Update(msg, s)
			cmds = append(cmds, cmd)
		}

		m.refreshMaxHP(LeftSide)
		m.refreshMaxHP(RightSide)

		return m, tea.Batch(cmds...)

	case metronomeColVisibleRequestMsg:
		var cmd tea.Cmd
		m.tabContent[movesTab], cmd = m.tabContent[movesTab].Update(msg, side)
		return m, cmd

	case autoFillForSpecieRequestMsg:
		var cmds []tea.Cmd
		var cmd tea.Cmd

		m.tabContent[statsTab], cmd = m.tabContent[statsTab].Update(msg, side)
		cmds = append(cmds, cmd)

		m.tabContent[miscTab], cmd = m.tabContent[miscTab].Update(msg, side)
		cmds = append(cmds, cmd)

		m.refreshMaxHP(side)

		return m, tea.Batch(cmds...)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, configKeys.NextTab):
			if (m.activeTab + 1) > (len(m.tabs) - 1) {
				m.activeTab = 0
			} else {
				m.activeTab++
			}
			return m, nil

		case key.Matches(msg, configKeys.PrevTab):
			if (m.activeTab - 1) < 0 {
				m.activeTab = len(m.tabs) - 1
			} else {
				m.activeTab--
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.tabContent[m.activeTab], cmd = m.tabContent[m.activeTab].Update(msg, side)
	return m, cmd
}

func (m configModel) View(side Side) string {
	content := m.Title(side)
	content += "\n\n"
	content += m.Body(side)
	return content
}

func (m configModel) Title(side Side) string {
	var sideDesc string
	if side == LeftSide {
		sideDesc = string(models.Attacker)
	} else {
		sideDesc = string(models.Defender)
	}

	return fmt.Sprintf("%s %s\n", CONFIG_TITLE, sideDesc)
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
		Name:          m.Name(side),
		Level:         m.Level(side),
		Nature:        m.Nature(side),
		Gender:        m.Gender(side),
		AttackMove:    attackMove,
		Moves:         m.Moves(side),
		HealthPercent: m.HealthPercent(side),
		// FUTURE FEATURES: implement shortcut iv and ev settings (optimal, worstcase) input option
		IvSetting:     "custom",
		IvCustom:      m.IvCustomStats(side),
		EvSetting:     "custom",
		EvCustom:      m.EvCustomStats(side),
		Types:         m.Types(side),
		MoveOverride:  m.MoveOverride(side, attackMove),
		Ability:       m.Ability(side),
		AbilityOn:     m.AbilityOn(side),
		Item:          m.Item(side),
		Status:        m.Status(side),
		Hits:          m.Hits(side, attackMove),
		MetronomeHits: m.MetronomeHits(side, attackMove),
		Boosts:        m.Boosts(side),
		InputSide:     m.InputSide(side),
	}
}
