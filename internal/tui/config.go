package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type configModel struct {
	tabs		[]string
	tabContent	[]string
	activeTab	int
	styles		*styles
}

func newConfigModel(s *styles) configModel {
	tabs := []string{"Basics", "Stats", "Misc.", "Moves"}
	tabContent := []string{"Basics here", "Stats here", "Misc. here", "Moves here"}

	return configModel{
		tabs:		tabs,
		tabContent:	tabContent,
		activeTab:	0,
		styles:		newStyles(true),
	}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (configModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "right", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil

		case "left", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	return m, nil
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
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = s.activeTab
		} else {
			style = s.inactiveTab
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(s.window.Width((lipgloss.Width(row))).Render(m.tabContent[m.activeTab]))
	return s.doc.Render(doc.String())
}
