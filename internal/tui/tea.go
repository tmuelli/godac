package tui

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	view		View
	calcTab		table.Model
}

type View string

const (
	CalcView	View = "calc"
	ConfigView	View = "config"
	FieldView	View = "field"
)

func (v View) Title() string {
	switch v {
	case CalcView:
		return "Calculation view - Press [R] or [L] for config view or [F] for field view"
	case ConfigView:
		return "Config view - Press [C] for calculation view or [F] for field view"
	case FieldView:
		return "Field view - Press [C] for calculation view or [R] or [L] for config view"
	}
	return "Error: Unknown view"
}

func (v View) Body(m model) string {
	switch v {
	case CalcView:
		return CalcViewBody(m)
	case ConfigView:
		return ConfigViewBody(m)
	case FieldView:
		return FieldViewBody(m)
	}
	return "Error: Unknown view"
}

func CalcViewBody(m model) string {
	body := m.calcTab.View()
	body += "\n\nChoose one move\n"
	body += "\nPokemon1 (Move x) vs. Pokemon 2: MinDmg - MaxDmg (MinPerc - MaxPerc)\n"
	body += "\nPossible dmg amounts: (x, ...)\n"
	return body
}

func ConfigViewBody(m model) string {
	return ""
}

func FieldViewBody(m model) string {
	return ""
}

func NewModel() model {
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

	return model{
		view: CalcView,
		calcTab: t,
	}
}

func Run() error {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// quit terminal repl
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return updateState(msg, m)
}

func (m model) View() tea.View {
	var s string

	s = fmt.Sprintf("\n%s\n", viewState(m))

	return tea.NewView(s)
}

func updateState(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "c", "r", "l", "f":
			changeView(&m.view, msg.String())
		}
	}

	return m, nil
}

func viewState(m model) string {
	content := m.view.Title()
	content += "\n\n"
	content += m.view.Body(m)
	return content
}

func changeView(current *View, newViewShortcut string) {
	switch newViewShortcut {
	case "c":
		if *current == CalcView {
			return
		}

		*current = CalcView
	case "r", "l":
		if *current == ConfigView {
			return
		}

		*current = ConfigView
	case "f":
		if *current == FieldView {
			return
		}

		*current = FieldView
	}
}