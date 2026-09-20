package tui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type helpKeyMap struct {
	global		globalKeyMap
	view		help.KeyMap
}

func (h helpKeyMap) ShortHelp() []key.Binding {
	return append(h.global.ShortHelp(), h.view.ShortHelp()...)
}

func (h helpKeyMap) FullHelp() [][]key.Binding {
	return append(h.global.FullHelp(), h.view.FullHelp()...)
}