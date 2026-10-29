package tui

import "charm.land/bubbles/v2/key"

type globalKeyMap struct {
	Calc    key.Binding
	Config  key.Binding
	Field   key.Binding
	Help    key.Binding
	Quit    key.Binding
	Confirm key.Binding
	Cancel  key.Binding
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
}

var globalKeys = globalKeyMap{
	// help menu
	Calc:   key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "calc view")),
	Config: key.NewBinding(key.WithKeys("r", "l"), key.WithHelp("r/l", "config view")),
	Field:  key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "field view")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "quit")),

	// edit mode
	Confirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
	Cancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel editing")),

	// navigation
	Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("up", "navigate up")),
	Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("down", "navigate down")),
	Left:  key.NewBinding(key.WithKeys("left"), key.WithHelp("left", "navigate left")),
	Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("right", "navigate right")),
}

// display-only key bindings for help menu in footer => not used anywhere (e.g. key.Matches)
var (
	hintMove     = key.NewBinding(key.WithKeys("up", "down", "left", "right"), key.WithHelp("↑↓←→", "move"))
	hintMoveRows = key.NewBinding(key.WithKeys("up", "down"), key.WithHelp("↑↓", "move"))
	hintRows     = key.NewBinding(key.WithKeys("up", "down"), key.WithHelp("↑↓", "row"))
	hintChange   = key.NewBinding(key.WithKeys("left", "right"), key.WithHelp("←→", "change"))
	hintEdit     = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "edit"))
	hintCalc     = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "calculate"))
	hintViews    = key.NewBinding(key.WithKeys("c", "r", "l", "f"), key.WithHelp("c/r/l/f", "views"))
	hintConfirm  = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm"))
	hintCancel   = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
	hintType     = key.NewBinding(key.WithHelp("type", "to filter"))
)

func (k globalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{hintViews, k.Help, k.Quit}
}

func (k globalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Calc, k.Config, k.Field},
		{k.Up, k.Down, k.Left, k.Right},
		{k.Confirm, k.Cancel},
		{k.Help, k.Quit},
	}
}
