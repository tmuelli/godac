package tui

import tea "charm.land/bubbletea/v2"

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

func requestTextInput(fieldID, initial string) tea.Cmd {
	return func() tea.Msg {
		return editRequestMsg{fieldID: fieldID, initial: initial}
	}
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