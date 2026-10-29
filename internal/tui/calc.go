package tui

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/calc"
	"github.com/tmuelli/godac/internal/models"
)

const (
	basicsTab = iota
	statsTab
	miscTab
	movesTab
)

const (
	CALC_MOVE_WIDTH = 18

	LEFT_PREFIX  = "L"
	RIGHT_PREFIX = "R"
)

type calcKeyMap struct{}

var calcKeys = calcKeyMap{}

func (k calcKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{hintMove, hintCalc}
}

func (k calcKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{}}
}

type calcContext struct {
	leftName   string
	rightName  string
	leftMoves  []string
	rightMoves []string
}

type calcResult struct {
	r        models.BattleResult
	err      error
	attacker string
	defender string
	move     string
}

type calcModel struct {
	grid   gridData
	styles *styles
	data   *models.Data
	result calcResult
}

func newCalcModel(s *styles, data *models.Data) calcModel {
	columns := []gridColumn{
		{col: table.Column{}, editable: true, visible: true},
		{col: table.Column{}, editable: true, visible: true},
	}

	rows := [][]string{{"", ""}, {"", ""}, {"", ""}, {"", ""}}

	return calcModel{
		grid: gridData{
			columns: columns,
			rows:    rows,
			row:     0,
			col:     0,
		},
		styles: s,
		data:   data,
	}
}

func (m calcModel) Init() tea.Cmd {
	return nil
}

func (m calcModel) Update(msg tea.Msg) (calcModel, tea.Cmd) {
	switch msg := msg.(type) {
	case userInputMsg:
		attacker, defender := msg.input.atkInput.Name, msg.input.defInput.Name

		m.result = calcResult{attacker: attacker, defender: defender, move: msg.move}
		m.Calculate(msg.input)
		return m, nil
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, globalKeys.Up):
			m.grid.MoveRow(-1)
		case key.Matches(msg, globalKeys.Down):
			m.grid.MoveRow(1)
		case key.Matches(msg, globalKeys.Left):
			m.grid.MoveCol(-1)
		case key.Matches(msg, globalKeys.Right):
			m.grid.MoveCol(1)
		case key.Matches(msg, globalKeys.Confirm):
			usedSide := Side(m.grid.col)
			usedMoveIndex := m.grid.row
			return m, requestCalc(usedMoveIndex, usedSide)
		}
	}

	return m, nil
}

func (m calcModel) View(ctx calcContext) string {
	return m.Body(ctx)
}

func (m calcModel) Body(ctx calcContext) string {
	left := m.sidePanel(ctx.leftName, ctx.leftMoves, LeftSide)
	right := m.sidePanel(ctx.rightName, ctx.rightMoves, RightSide)
	panels := lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)

	parts := []string{panels, ""}

	// add hint to setup pokemon via config view
	if ctx.leftName == "" && ctx.rightName == "" {
		parts = append(parts, m.styles.dimText.Render(CALC_SETUP_HINT))
	}
	parts = append(parts, m.ResultView(lipgloss.Width(panels)))

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m calcModel) sidePanel(name string, slots []string, side Side) string {
	s := m.styles
	active := m.grid.col == int(side)

	prefix := LEFT_PREFIX
	if side == RightSide {
		prefix = RIGHT_PREFIX
	}

	title := name
	if title == "" {
		title = DISPLAY_MOVE_NOT_SET
	}
	title += " [" + prefix + "]"

	rows := make([]string, len(slots))
	for i, move := range slots {
		marker := UNSELECTED_MARKER
		nameStyle := lipgloss.NewStyle().Width(CALC_MOVE_WIDTH)
		text := move

		if move == "" {
			text = NOT_SET
			nameStyle = s.dimText.Width(CALC_MOVE_WIDTH)
		}

		if active && i == m.grid.row {
			marker = s.selected.Render(SELECTED_MARKER)
			if move != "" {
				nameStyle = s.selected.Width(CALC_MOVE_WIDTH)
			}
		}

		rows[i] = marker + nameStyle.Render(text)
	}

	return s.panel(title, lipgloss.JoinVertical(lipgloss.Left, rows...), active)
}

func formatHitDisplayMsg(hitCount int) string {
	return fmt.Sprintf("%d.%s:", hitCount, DISPLAY_HIT)
}

func renderDmgRow(row []int) string {
	var b strings.Builder
	for i, dmg := range row {
		if i > 0 {
			b.WriteString(", ")
		}

		b.Write(strconv.AppendInt(nil, int64(dmg), 10))
	}
	return b.String()
}

func (m calcModel) dmgRollLines(inner int) []string {
	matrix := m.result.r.DmgMatrix
	isMultihit := len(matrix) > 1

	lines := make([]string, 0, len(matrix))
	for i, row := range matrix {
		text := renderDmgRow(row)
		if isMultihit {
			text = formatHitDisplayMsg(i+1) + " " + text
		}
		lines = append(lines, m.styles.dimText.Width(inner).Render(text))
	}

	return lines
}

func (m calcModel) ResultView(width int) string {
	s := m.styles
	r := m.result
	inner := width - HORIZONTAL_PADDING

	box := func(content string) string {
		return s.window.Padding(0, 1).BorderForeground(s.border).Width(width).Render(content)
	}

	if r.err != nil {
		return box(s.errText.Width(inner).Render("Error: " + r.err.Error()))
	}

	if len(r.r.DmgMatrix) == 0 {
		return s.dimText.Render(CHOOSE_MOVE)
	}

	heading := lipgloss.NewStyle().Bold(true).
		Render(fmt.Sprintf("%s (%s) vs. %s", r.attacker, r.move, r.defender))

	dmgRange := lipgloss.NewStyle().Bold(true).
		Render(fmt.Sprintf("%d - %d", r.r.Range.Min, r.r.Range.Max))

	percRange := s.selected.
		Render(fmt.Sprintf("(%.2f%% - %.2f%%)", r.r.Range.MinPerc, r.r.Range.MaxPerc))

	lines := []string{heading, "", dmgRange + " " + percRange, "", s.dimText.Render(DAMAGE_AMOUNTS)}
	lines = append(lines, m.dmgRollLines(inner)...)

	return box(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func (m *calcModel) Calculate(input *UserInput) {
	// create new calculation
	calculation, err := calc.NewCalc(m.data, &input.atkInput, &input.defInput, &input.fieldInput)
	if err != nil {
		// log.Fatal("Error while initializing calculation data:", err)
		m.result.err = err
		return
	}

	err = calculation.Calc()
	if err != nil {
		// log.Fatal("Error while damage calculation: ", err)
		m.result.err = err
		return
	}

	m.result.r = models.BattleResult{
		Range: models.BattleResultRange{
			Min:     calculation.Result.Range.Min,
			MinPerc: calculation.Result.Range.MinPerc,
			Max:     calculation.Result.Range.Max,
			MaxPerc: calculation.Result.Range.MaxPerc,
		},
		DmgMatrix: calculation.Result.DmgMatrix,
	}
}
