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
	usedMoveIndex	int
	usedSide		Side
}

type calcRequestInfo struct {
	usedMove	string
	usedSide	Side
}

func requestCalc(moveIndex int, side Side) tea.Cmd {
	return func() tea.Msg {
		return calcRequestMsg{
			usedMoveIndex:	moveIndex,
			usedSide:		side,
		}
	}
}

type autoFillForSpecieRequestMsg struct {
	specie		*models.Specie
}

func requestAutoFillForSpecie(specie *models.Specie) tea.Cmd {
	return func() tea.Msg {
		return autoFillForSpecieRequestMsg{
			specie:		specie,
		}
	}
}

type metronomeColVisibleRequestMsg struct {
	visible		bool
}

func requestMetronomeColVisible(visible bool) tea.Cmd {
	return func() tea.Msg {
		return metronomeColVisibleRequestMsg{
			visible:	visible,
		}
	}
}