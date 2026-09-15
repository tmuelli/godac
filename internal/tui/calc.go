package tui

import (
	"fmt"

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

type calcContext struct {
	leftName	string
	rightName	string
	leftMoves	[]string
	rightMoves	[]string
}

type calcResult struct {
	r		models.BattleResult
	err		error
}

type calcModel struct {
	grid		gridData
	styles		*styles
	data		*models.Data
	result		calcResult
}

func newCalcModel(s *styles, data *models.Data) calcModel {
	columns := []table.Column{
		{Title: "Pokemon1 [L]", Width: 15},
		{Title: "Pokemon2 [R]", Width: 15},
	}

	rows := [][]string{
		{"Move 1 [L1]", "Move 1 [R1]"},
		{"Move 2 [L2]", "Move 2 [R2]"},
		{"Move 3 [L3]", "Move 3 [R3]"},
		{"Move 4 [L4]", "Move 4 [R4]"},
	}

	return calcModel{
		grid:	gridData{
			columns:	columns,
			rows:		rows,
			row:		0,
			col:		0,
		},
		styles:		s,
		data:		data,
	}
}

func (m calcModel) Init() tea.Cmd {
	return nil
}

func (m calcModel) Update(msg tea.Msg) (calcModel, tea.Cmd) {
	switch msg := msg.(type) {
	case userInputMsg:
		m.Calculate(msg.input)
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
			case "up":
				m.grid.MoveRow(-1)
			case "down":
				m.grid.MoveRow(1)
			case "left":
				m.grid.MoveCol(-1)
			case "right":
				m.grid.MoveCol(1)
			case "enter":
				usedSide := Side(m.grid.col)
				usedMove := m.grid.rows[m.grid.row][m.grid.col]
				return m, requestCalc(usedMove, usedSide)
		}
	}

	return m, nil
}

func (m calcModel) View(ctx calcContext) string {
	content := m.Title()
	content += "\n\n"
	content += m.Body(ctx)
	return content
}

func (m calcModel) Title() string {
	return "Calculation view - Press [R] or [L] for config view or [F] for field view"
}

func (m calcModel) Body(ctx calcContext) string {
	m.fillFromContext(ctx)

	body := renderGrid(m.grid, m.styles)
	body += "\n\nChoose one move\n"
	body += m.ResultView()
	return body
}

func (m calcModel) ResultView() string {
	content := "\nPokemon1 (Move x) vs. Pokemon 2:\n"
	content += fmt.Sprintf("%d - %d (%f - %f)\n",
		m.result.r.Range.Min,
		m.result.r.Range.Max,
		m.result.r.Range.MinPerc,
		m.result.r.Range.MaxPerc,
	)
	content += "\nPossible dmg amounts:\n"
	// TODO: format damage matrix
	return content
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
		m.grid.columns[side].Title = fmt.Sprintf("%s [%s]", sideName, sidePrefix)
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
			Min: calculation.Result.Range.Min,
			MinPerc: 0,
			Max: calculation.Result.Range.Max,
			MaxPerc: 0,
		},
		DmgMatrix: calculation.Result.DmgMatrix,
	}
}