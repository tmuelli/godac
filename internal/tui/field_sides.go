package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SideFocus int

const (
	Left			SideFocus = iota
	Right

	CHECKED = "[x]"
	UNCHECKED = "[ ]"
)

type sidesGroup struct {
	label			string
	left			bool
	right			bool
	focused			SideFocus
	groupFocused	bool
	styles			*styles
}

func newSidesGroup(s *styles, label string, defaultFocus SideFocus, defaultLeft bool, defaultRight bool) sidesGroup {
	return sidesGroup{
		label: label,
		left: defaultLeft,
		right: defaultRight,
		focused: defaultFocus,
		styles: s,
	}
}

func (g sidesGroup) Update(msg tea.Msg) (sidesGroup, tea.Cmd) {
	if !g.groupFocused {
		return g, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left":
			g.focused = Left
		case "right":
			g.focused = Right
		case "enter":
			if g.focused == Left {
				g.left = !g.left
			} else {
				g.right = !g.right
			}
		}
	}
	
	return g, nil
}

func (g sidesGroup) View() string {
	isGroupFocused := g.groupFocused

	// left
	leftCol := UNCHECKED
	if g.left {
		leftCol = CHECKED
	}

	leftStyle := lipgloss.NewStyle().Bold(isGroupFocused)
	if g.focused == Left {
		leftStyle = leftStyle.Foreground(g.styles.highlight)
	} else {
		leftStyle.Foreground(lipgloss.NoColor{})
	}

	// label
	lableStyle := lipgloss.NewStyle().Bold(isGroupFocused)

	// right
	rightCol := UNCHECKED
	if g.right {
		rightCol = CHECKED
	}

	rightStyle := lipgloss.NewStyle().Bold(isGroupFocused)
	if g.focused == Right {
		rightStyle = rightStyle.Foreground(g.styles.highlight)
	} else {
		rightStyle.Foreground(lipgloss.NoColor{})
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftCol),
		" ",
		lableStyle.Render(g.label),
		" ",
		rightStyle.Render(rightCol),
	)
}

// TODO: implement value getter for spikes
func (g sidesGroup) Value(side SideFocus) bool {
	if side == Left {
		return g.left
	} else {
		return g.right
	}
}

func (m fieldModel) FormatSideSettings() string {
	var title string
	var lines []string

	titleStyle := lipgloss.NewStyle()
	if m.setting == SidesSetting {
		titleStyle = titleStyle.Bold(true).Foreground(lipgloss.Color("#FFFF00"))
	} else {
		titleStyle.Foreground(lipgloss.NoColor{})
	}

	title = "Sides"

	for _, s := range m.sides {
		lines = append(lines, s.View())
	}

	return "\n" + titleStyle.Render(title) + ":\n" + lipgloss.JoinVertical(lipgloss.Left, lines...)
}

