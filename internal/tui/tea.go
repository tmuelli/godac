package tui

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
)

type View int

const (
	CalcView View = iota
	ConfigView
	FieldView
)

const (
	SHORT_HELP_SEPERATOR = " · "
)

type model struct {
	current    View
	calc       calcModel
	configSide Side
	config     configModel
	field      fieldModel

	styles *styles
	data   *models.Data

	editing          bool
	editMode         EditMode
	editID           string
	editOptions      []string
	editMin, editMax int
	input            textinput.Model
	inputValid       bool

	calcRequestInfo calcRequestInfo

	help help.Model

	width  int
	height int
}

func NewModel(data *models.Data) model {
	s := newStyles(true)

	m := model{
		current:    CalcView,
		calc:       newCalcModel(s, data),
		config:     newConfigModel(s, data),
		field:      newFieldModel(s, data),
		data:       data,
		styles:     s,
		inputValid: true,
		help:       help.New(),
	}
	m.help.Styles = s.help
	m.help.ShortSeparator = SHORT_HELP_SEPERATOR

	return m
}

func Run(data *models.Data) error {
	p := tea.NewProgram(NewModel(data))
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func (m model) Init() tea.Cmd {
	return tea.RequestBackgroundColor
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(max(msg.Width-HORIZONTAL_PADDING, DEFAULT_MIN_WIDTH))
		return m, nil

	case tea.BackgroundColorMsg:
		// detected background color => set styles accordingly
		m.styles.apply(msg.IsDark())
		m.help.Styles = m.styles.help
		return m, nil

	case setLevelRequestMsg:
		var cmd tea.Cmd
		m.config, cmd = m.config.Update(msg, m.configSide)
		return m, cmd

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
			input: userInput,
			move:  m.calcRequestInfo.usedMove,
		}

		var cmd tea.Cmd
		m.calc, cmd = m.calc.Update(inputMsg)
		return m, cmd

	case editRequestMsg:
		m.editing = true
		m.editID = msg.fieldID
		m.editMode = msg.mode
		m.inputValid = true

		m.input = textinput.New()
		m.input.SetStyles(m.styles.input)
		m.input.Prompt = msg.fieldID + INPUT_MARKER
		m.input.Placeholder = msg.initial

		switch m.editMode {
		case SelectEditMode:
			m.input.SetSuggestions(msg.options)
			m.input.ShowSuggestions = true
			m.editOptions = msg.options

		case NumberEditMode:
			m.editMin = msg.min
			m.editMax = msg.max
		}

		m.input.Focus()
		return m, nil

	case tea.KeyPressMsg:
		if m.editing {
			switch {
			case key.Matches(msg, globalKeys.Confirm):
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
					if err != nil || n < m.editMin || n > m.editMax {
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
			case key.Matches(msg, globalKeys.Cancel):
				m.editing = false
				return m, nil
			}

			// reset errors (if any) on new keystroke
			m.inputValid = true

			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch {
		// quit terminal repl
		case key.Matches(msg, globalKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, globalKeys.Calc):
			m.current = CalcView
			return m, nil

		case key.Matches(msg, globalKeys.Config):
			if msg.String() == "r" {
				m.current = ConfigView
				m.configSide = RightSide
				return m, nil
			} else {
				m.current = ConfigView
				m.configSide = LeftSide
				return m, nil
			}

		case key.Matches(msg, globalKeys.Field):
			m.current = FieldView
			return m, nil

		case key.Matches(msg, globalKeys.Help):
			m.help.ShowAll = !m.help.ShowAll
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
	doc.WriteString(m.header())
	doc.WriteString("\n\n")

	switch m.current {
	case CalcView:
		ctx := calcContext{
			leftName:   m.config.Name(LeftSide),
			rightName:  m.config.Name(RightSide),
			leftMoves:  m.config.Moves(LeftSide),
			rightMoves: m.config.Moves(RightSide),
		}
		content = m.calc.View(ctx)

	case ConfigView:
		content = m.config.View(m.configSide)

	case FieldView:
		content = m.field.View()
	}

	doc.WriteString(content)

	// footer for help menu
	footer := m.help.View(m.KeyMap())
	bottom := m.editBar() + "\n" + footer
	top := doc.String()

	// pin the footer to the bottom
	pad := max(0, m.height-2-lipgloss.Height(top)-lipgloss.Height(bottom))
	page := top + strings.Repeat("\n", pad) + bottom

	v := tea.NewView(m.styles.doc.Render(page))
	v.AltScreen = true
	return v
}

func (m model) Input() *UserInput {
	input := UserInput{
		leftInput:  *m.config.Input(LeftSide, &m.calcRequestInfo),
		rightInput: *m.config.Input(RightSide, &m.calcRequestInfo),
		fieldInput: *m.field.Input(),
	}

	return &input
}

func (m model) KeyMap() help.KeyMap {
	if m.editing {
		return editKeyMap{
			selecting: m.editMode == SelectEditMode,
		}
	}

	switch m.current {
	case ConfigView:
		return helpKeyMap{global: globalKeys, view: configKeys}
	case FieldView:
		return helpKeyMap{global: globalKeys, view: fieldKeys}
	default:
		return helpKeyMap{global: globalKeys, view: calcKeys}
	}
}

func (m model) editHint() string {
	s := m.styles

	if !m.inputValid {
		switch m.editMode {
		case NumberEditMode:
			return s.errText.Render(fmt.Sprintf(ERR_INVALID_NUMBER_INPUT, m.editMin, m.editMax))
		default:
			return s.errText.Render(fmt.Sprintf(ERR_INVALID_INPUT, m.input.Value()))
		}
	}

	if m.editMode == NumberEditMode {
		return s.dimText.Render(fmt.Sprintf("%d-%d", m.editMin, m.editMax))
	}

	matches := m.input.MatchedSuggestions()
	switch {
	case m.input.Value() == "":
		return s.dimText.Render(fmt.Sprintf(EDIT_TYPE_HINT, len(m.editOptions)))
	case len(matches) == 0:
		return s.dimText.Render(EDIT_TYPE_HINT_NO_OPTIONS)
	}

	showAvailOpts := matches[:min(len(matches), DEFAULT_EDIT_HINTS)]
	hint := strings.Join(showAvailOpts, HINT_SEPERATOR)
	if extendedOpts := len(matches) - len(showAvailOpts); extendedOpts > 0 {
		hint += fmt.Sprintf(EDIT_TYPE_HINTS_EXTENDED, extendedOpts)
	}

	return s.dimText.Render(hint)
}

func (m model) editBar() string {
	s := m.styles

	if !m.editing {
		return lipgloss.NewStyle().Height(DEFAULT_EDIT_BAR_HEIGHT).Render("")
	}

	w := DEFAULT_WIDTH
	if m.width > 0 {
		w = min(max(m.width-HORIZONTAL_PADDING, DEFAULT_MIN_WIDTH), DEFAULT_WIDTH)
	}

	border := s.accent
	if !m.inputValid {
		border = s.err
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(border).
		Padding(0, 1).
		Width(w).
		Render(m.input.View())

	return box + "\n" + m.editHint()
}
