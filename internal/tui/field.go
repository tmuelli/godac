package tui

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

type fieldModel struct {
	tab		table.Model
}

func newFieldModel() fieldModel {
	columns := []table.Column{
		{Title: "Left [R]", Width: 15},
		{Title: "Field", Width: 5},
		{Title: "Right [L]", Width: 15},
	}

	rows := []table.Row{
		{"[]", "Singles | Doubles", "[]"},
		{"[]", "Lightscreen", "[]"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
		table.WithWidth(40),
	)

	return fieldModel{
		tab: t,
	}
}

func (m fieldModel) Init() tea.Cmd {
	return nil
}

func (m fieldModel) Update(msg tea.Msg) (fieldModel, tea.Cmd) {
	return m, nil
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
	return ""
}