package tui

import (
	"slices"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
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
	configSide	Side
	config		configModel
	field		fieldModel

	styles		*styles
	data		*models.Data

	editing				bool
	editMode			EditMode
	editID				string
	editOptions			[]string
	editMin, editMax	int
	input				textinput.Model
	inputValid			bool

	calcRequestInfo		calcRequestInfo
}

func NewModel(data *models.Data) model {
	s := newStyles(true)

	return model{
		current:	CalcView,
		calc:		newCalcModel(s, data),
		config:		newConfigModel(s, data),
		field:		newFieldModel(s, data),
		data:		data,
		styles:		s,
		inputValid:	true,
	}
}

func Run(data *models.Data) error {
	p := tea.NewProgram(NewModel(data))
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
	case calcRequestMsg:
		moveSettings := m.config.MoveSettings(msg.usedSide)
		moves := m.config.Moves(msg.usedSide)

		if !moveSettings[msg.usedMoveIndex].isSet {
			return m, nil
		}

		m.calcRequestInfo.usedMove = moves[msg.usedMoveIndex]
		m.calcRequestInfo.usedSide = msg.usedSide

		userInput := m.Input()
		inputMsg := userInputMsg{
			input:	userInput,
		}

		var cmd tea.Cmd
		m.calc, cmd = m.calc.Update(inputMsg)
		return m, cmd

	case editRequestMsg:
		m.editing = true
		m.editID = msg.fieldID
		m.editMode = msg.mode

		m.input = textinput.New()

		switch m.editMode {
		case SelectEditMode:
			m.input.SetSuggestions(msg.options)
			m.input.ShowSuggestions = true
			m.editOptions = msg.options

		case NumberEditMode:
			m.editMin = msg.min
			m.editMax = msg.max
		}
		
		// m.input.SetValue(msg.initial)
		// m.input.Placeholder = msg.initial
		m.input.Focus()
		return m, nil

	case tea.KeyPressMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				value := m.input.Value()
				m.inputValid = true

				switch m.editMode {
				case SelectEditMode:
					if !slices.Contains(m.editOptions, value) {
						m.inputValid = false
						return m, nil
					}

				case NumberEditMode:
					n, err := strconv.Atoi(value)
					if err != nil || n < m.editMin || n >m.editMax {
						m.inputValid = false
						return m, nil
					}
				}

				// input is valid => editing can now be exited
				m.editing = false

				var cmd tea.Cmd
				editMsg := editConfirmMsg{fieldID: m.editID, value: value}
				switch m.current {
				case CalcView:
					m.calc, cmd = m.calc.Update(editMsg)

				case ConfigView:
					m.config, cmd = m.config.Update(editMsg, m.configSide)

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

		case "r":
			m.current = ConfigView
			m.configSide = RightSide
			return m, nil

		case "l":
			m.current = ConfigView
			m.configSide = LeftSide
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
		m.config, cmd = m.config.Update(msg, m.configSide)

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
		ctx := calcContext{
			leftName:			m.config.Name(LeftSide),
			rightName:			m.config.Name(RightSide),
			leftMoves:			m.config.Moves(LeftSide),
			rightMoves:			m.config.Moves(RightSide),
		}
		content = m.calc.View(ctx)

	case ConfigView:
		content = m.config.View(m.configSide)

	case FieldView:
		content = m.field.View()
	}

	doc.WriteString(content)

	if m.editing {
		doc.WriteString("\n\n")
		doc.WriteString(m.input.View())

		if !m.inputValid {
			// TODO: better error message and better visualization
			doc.WriteString("\nInvalid input. Please try again.\n")
		}
	}

	v := tea.NewView(
		// s.doc.Render(doc.String()),
		doc.String(),
	)
	v.AltScreen = true
	return v
}

func (m model) Input() *UserInput {
	input := UserInput{
		leftInput:	*m.config.Input(LeftSide, &m.calcRequestInfo),
		rightInput:	*m.config.Input(RightSide, &m.calcRequestInfo),
		fieldInput:	*m.field.Input(),
	}

	return &input
}
