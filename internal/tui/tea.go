package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
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

	styles		*styles

	editing		bool
	input		textinput.Model
	editID		string
}

func NewModel() model {
	s := newStyles(true)

	return model{
		current:	ConfigView,
		calc:		newCalcModel(s),
		config:		newConfigModel(s),
		field:		newFieldModel(s),
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
	case editRequestMsg:
		m.editing = true
		m.editID = msg.fieldID
		m.input = textinput.New()
		// m.input.SetValue(msg.initial)
		m.input.Placeholder = msg.initial
		m.input.Focus()
		return m, nil

	case tea.KeyPressMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				value := m.input.Value()
				m.editing = false

				var cmd tea.Cmd
				editMsg := editConfirmMsg{fieldID: m.editID, value: value}
				switch m.current {
				case CalcView:
					m.calc, cmd = m.calc.Update(editMsg)

				case ConfigView:
					m.config, cmd = m.config.Update(editMsg)

				case FieldView:
					m.field, cmd = m.field.Update(editMsg)
				}
				return m, cmd
			case "esc":
				m.editing = false
				return m, nil
			}

			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

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

	doc := strings.Builder{}
	// s := m.styles

	switch m.current {
	case CalcView:
		content = m.calc.View()

	case ConfigView:
		content = m.config.View()

	case FieldView:
		content = m.field.View()
	}

	doc.WriteString(content)

	if m.editing {
		doc.WriteString("\n\n")
		doc.WriteString(m.input.View())
	}

	v := tea.NewView(
		// s.doc.Render(doc.String()),
		doc.String(),
	)
	v.AltScreen = true
	return v
}

type editRequestMsg struct {
	fieldID		string
	initial		string
}

type editConfirmMsg struct {
	fieldID		string
	value		string
}

func requestTextInput(fieldID, initial string) tea.Cmd {
	return func() tea.Msg {
		return editRequestMsg{fieldID: fieldID, initial: initial}
	}
}
