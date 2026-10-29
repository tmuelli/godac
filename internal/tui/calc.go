package tui

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/calc"
	"github.com/tmuelli/godac/internal/models"
)

const (
	basicsTab = iota
	statsTab
	miscTab
	movesTab
)

type calcKeyMap struct{}

var calcKeys = calcKeyMap{}

func (k calcKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
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
	r   models.BattleResult
	err error
}

type calcModel struct {
	grid     gridData
	styles   *styles
	data     *models.Data
	result   calcResult
	usedMove string
}

func newCalcModel(s *styles, data *models.Data) calcModel {
	columns := []gridColumn{
		{col: table.Column{Title: "Pokemon1 [L]", Width: 15}, editable: true, visible: true},
		{col: table.Column{Title: "Pokemon2 [R]", Width: 15}, editable: true, visible: true},
	}

	rows := [][]string{
		{"Move 1 [L1]", "Move 1 [R1]"},
		{"Move 2 [L2]", "Move 2 [R2]"},
		{"Move 3 [L3]", "Move 3 [R3]"},
		{"Move 4 [L4]", "Move 4 [R4]"},
	}

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
		m.usedMove = msg.move
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
	m.fillFromContext(ctx)

	body := renderGrid(m.grid, m.styles)
	body += fmt.Sprintf("\n\n%s\n", CHOOSE_MOVE)
	body += m.ResultView(ctx)
	return body
}

func formatHitDisplayMsg(hitCount int) string {
	return fmt.Sprintf("%d.%s", hitCount, DISPLAY_HIT)
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

func renderDmgMatrix(matrix [][]int) string {
	result := ""
	isMultihit := false

	if len(matrix) == 0 {
		return ERR_INVALID_OUTPUT
	}

	if len(matrix) > 1 {
		isMultihit = true
	}

	// first damage row => always present
	var prefix string
	if isMultihit {
		prefix = formatHitDisplayMsg(1)
	}

	result += fmt.Sprintf("\n%s\n", DAMAGE_AMOUNTS)
	result += fmt.Sprintf("%s %s\n", prefix, renderDmgRow(matrix[0]))

	// additional lines (multihit output)
	if isMultihit {
		var hitStr strings.Builder

		for i, r := range matrix {
			if i == 0 {
				continue
			}

			fmt.Fprintf(&hitStr, "%s %s\n", formatHitDisplayMsg(i+1), renderDmgRow(r))
		}

		result += hitStr.String()
	}

	return result
}

func (m calcModel) buildBattleString(ctx calcContext) string {
	return fmt.Sprintf("\n%s (%s) vs. %s:\n", ctx.leftName, m.usedMove, ctx.rightName)
}

func (m calcModel) ResultView(ctx calcContext) string {
	if m.result.err != nil {
		return fmt.Sprintf("Error: %v\n", m.result.err)
	}

	if len(m.result.r.DmgMatrix) > 0 {
		content := m.buildBattleString(ctx)
		content += fmt.Sprintf("%d - %d (%.2f%% - %.2f%%)\n",
			m.result.r.Range.Min,
			m.result.r.Range.Max,
			m.result.r.Range.MinPerc,
			m.result.r.Range.MaxPerc,
		)
		content += renderDmgMatrix(m.result.r.DmgMatrix)
		return content
	}
	return ""
}

// FIXME: Move logic into Update rather than View => View should only render and not mutate
func (m calcModel) fillFromContext(ctx calcContext) {
	// left
	m.fillSideFromContext(ctx.leftName, ctx.leftMoves, LeftSide, "L")

	// right
	m.fillSideFromContext(ctx.rightName, ctx.rightMoves, RightSide, "R")
}

// FIXME: Move logic into Update rather than View => View should only render and not mutate
func (m calcModel) fillSideFromContext(sideName string, sideMoves []string, side Side, sidePrefix string) {
	if sideName != "" {
		m.grid.columns[side].col.Title = fmt.Sprintf("%s [%s]", sideName, sidePrefix)
	}

	for i, move := range sideMoves {
		if move == "" {
			continue
		}

		m.grid.rows[i][side] = fmt.Sprintf("%s [%s%d]", move, sidePrefix, i+1)
	}
}

func (m *calcModel) Calculate(input *UserInput) {
	// create new calculation
	calculation, err := calc.NewCalc(m.data, &input.leftInput, &input.rightInput, &input.fieldInput)
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
