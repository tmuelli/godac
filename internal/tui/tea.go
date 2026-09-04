package tui

import (
	tea "charm.land/bubbletea/v2"
)

type View int

const (
	CalcView	View = iota
	ConfigView
	FieldView
)

type model struct {
	current		View
	calc		calcModel
	config		configModel
	field		fieldModel
}

func NewModel() model {
	return model{
		current:	CalcView,
		calc:		newCalcModel(),
		config:		newConfigModel(),
		field:		newFieldModel(),
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
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		// quit terminal repl
		case "q", "esc", "ctrl+c":
			return m, tea.Quit

		case "c":
			m.current = CalcView
			return m, nil

		case "r", "l":
			m.current = ConfigView
			return m, nil

		case "f":
			m.current = FieldView
			return m, nil
		}
	}

	var cmd tea.Cmd
	switch m.current {
	case CalcView:
		m.calc, cmd = m.calc.Update(msg)

	case ConfigView:
		m.config, cmd = m.config.Update(msg)

	case FieldView:
		m.field, cmd = m.field.Update(msg)
	}

	return m, cmd
}

func (m model) View() tea.View {
	var content string
	switch m.current {
	case CalcView:
		content = m.calc.View()

	case ConfigView:
		content = m.config.View()

	case FieldView:
		content = m.field.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
