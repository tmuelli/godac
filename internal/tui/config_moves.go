package tui

import (
	tea "charm.land/bubbletea/v2"
)

type movesModel struct {
	styles		*styles
}

func newMovesModel(s *styles) movesModel {
	return movesModel{
		styles: s,
	}
}

func (m movesModel) Init() tea.Cmd {
	return nil
}

func (m movesModel) Update(msg tea.Msg) (tabContent, tea.Cmd) {
	return m, nil
}

func (m movesModel) View() string {
	return "Moves here"
}