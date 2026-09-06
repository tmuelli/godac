package tui

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

type calcModel struct {
	tab		table.Model
}

func newCalcModel(s *styles) calcModel {
	columns := []table.Column{
		{Title: "Pokemon1 [R]", Width: 15},
		{Title: "", Width: 5},
		{Title: "Pokemon2 [L]", Width: 15},
	}

	rows := []table.Row{
		{"Move 1 [R1]", "", "Move 1 [L1]"},
		{"Move 2 [R2]", "", "Move 2 [L2]"},
		{"Move 3 [R3]", "", "Move 3 [L3]"},
		{"Move 4 [R4]", "", "Move 4 [L4]"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
		table.WithWidth(40),
	)

	return calcModel{
		tab: t,
	}
}

func (m calcModel) Init() tea.Cmd {
	return nil
}

func (m calcModel) Update(msg tea.Msg) (calcModel, tea.Cmd) {
	return m, nil
}

func (m calcModel) View() string {
	content := m.Title()
	content += "\n\n"
	content += m.Body()
	return content
}

func (m calcModel) Title() string {
	return "Calculation view - Press [R] or [L] for config view or [F] for field view"
}

func (m calcModel) Body() string {
	body := m.tab.View()
	body += "\n\nChoose one move\n"
	body += "\nPokemon1 (Move x) vs. Pokemon 2: MinDmg - MaxDmg (MinPerc - MaxPerc)\n"
	body += "\nPossible dmg amounts: (x, ...)\n"
	return body
}
