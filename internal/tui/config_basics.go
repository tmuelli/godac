package tui

import (
	tea "charm.land/bubbletea/v2"
)

type basicsModel struct {
	styles		*styles
}

func newBasicsModel(s *styles) basicsModel {
	return basicsModel{
		styles: s,
	}
}

func (m basicsModel) Init() tea.Cmd {
	return nil
}

func (m basicsModel) Update(msg tea.Msg) (tabContent, tea.Cmd) {
	return m, nil
}

func (m basicsModel) View() string {
	return "Basics here"
}