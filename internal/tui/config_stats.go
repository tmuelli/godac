package tui

import (
	tea "charm.land/bubbletea/v2"
)

type statsModel struct {
	styles		*styles
}

func newStatsModel(s *styles) statsModel {
	return statsModel{
		styles: s,
	}
}

func (m statsModel) Init() tea.Cmd {
	return nil
}

func (m statsModel) Update(msg tea.Msg) (tabContent, tea.Cmd) {
	return m, nil
}

func (m statsModel) View() string {
	return "Stats here"
}