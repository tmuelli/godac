package tui

import (
	"slices"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type helpKeyMap struct {
	global globalKeyMap
	view   help.KeyMap
}

func (h helpKeyMap) ShortHelp() []key.Binding {
	return append(slices.Clone(h.view.ShortHelp()), h.global.ShortHelp()...)
}

func (h helpKeyMap) FullHelp() [][]key.Binding {
	return append(h.global.FullHelp(), h.view.FullHelp()...)
}

type editKeyMap struct{}

func (editKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{hintConfirm, hintCancel, hintType}
}

func (editKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{hintConfirm, hintCancel, hintType},
	}
}
