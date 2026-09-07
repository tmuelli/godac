package tui

import (
	tea "charm.land/bubbletea/v2"
)

type miscModel struct {
	styles		*styles
}

func newMiscModel(s *styles) miscModel {
	return miscModel{
		styles: s,
	}
}

func (m miscModel) Init() tea.Cmd {
	return nil
}

func (m miscModel) Update(msg tea.Msg) (tabContent, tea.Cmd) {
	return m, nil
}

func (m miscModel) View() string {
	return "Misc here"
}