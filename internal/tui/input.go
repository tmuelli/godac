package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
)

type EditMode int

const (
	NumberEditMode		EditMode = iota
	SelectEditMode
)

type editRequestMsg struct {
	fieldID		string
	initial		string

	mode		EditMode
	min, max	int
	options		[]string
}

type editConfirmMsg struct {
	fieldID		string
	value		string
}

func requestNumberInput(fieldID, initial string, min, max int) tea.Cmd {
	return func() tea.Msg {
		return editRequestMsg{
			fieldID:	fieldID,
			initial:	initial,
			min:		min,
			max:		max,
			mode:		NumberEditMode,
		}
	}
}

func requestSelectInput(fieldID, initial string, options []string) tea.Cmd {
	return func() tea.Msg {
		return editRequestMsg{
			fieldID:	fieldID,
			initial:	initial,
			options:	options,
			mode:		SelectEditMode,
		}
	}
}

type UserInput struct {
	fieldInput	models.UserFieldInput
	leftInput	models.UserInput
	rightInput	models.UserInput
}

type userInputMsg struct {
	input		*UserInput
}

type calcRequestMsg struct {
	usedMove	string
	usedSide	int
}

type calcRequestInfo struct {
	usedMove	string
	usedSide	int
}

func requestCalc(move string, side int) tea.Cmd {
	return func() tea.Msg {
		return calcRequestMsg{
			usedMove:	move,
			usedSide:	side,
		}
	}
}