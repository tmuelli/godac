package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type tabContent interface {
	Init()					tea.Cmd
	Update(msg tea.Msg)		(tabContent, tea.Cmd)
	View()					string
}

type configModel struct {
	tabs		[]string
	tabContent	[]tabContent
	activeTab	int
	styles		*styles
}

func newConfigModel(s *styles) configModel {
	tabs := []string{"Basics", "Stats", "Misc.", "Moves"}
	tabContent := []tabContent{
		newBasicsModel(s),
		newStatsModel(s),
		newMiscModel(s),
		newMovesModel(s),
	}

	return configModel{
		tabs:		tabs,
		tabContent:	tabContent,
		activeTab:	0,
		styles:		s,
	}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (configModel, tea.Cmd) {
	switch msg := msg.(type) {
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
	m.tabContent[m.activeTab], cmd = m.tabContent[m.activeTab].Update(msg)
	return m, cmd
}

func (m configModel) View() string {
	content := m.Title()
	content += "\n\n"
	content += m.Body()
	return content
}

func (m configModel) Title() string {
	return "Config view - Press [C] for calculation view or [F] for field view"
}

func (m configModel) Body() string {
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

	content := m.tabContent[m.activeTab].View()
	doc.WriteString(s.window.Width((lipgloss.Width(content))).Render(content))
	return s.doc.Render(doc.String())
}
