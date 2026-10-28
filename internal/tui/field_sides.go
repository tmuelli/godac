package tui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
)

const (
	CHECKED   = "[x]"
	UNCHECKED = "[ ]"
)

type sidesGroup struct {
	label        string
	values       side[int]
	maxValue     int
	focused      Side
	groupFocused bool
	styles       *styles
}

func newSidesGroup(s *styles, label string, defaultFocus Side, defaultLeft bool, defaultRight bool) sidesGroup {
	var values side[int]

	if defaultLeft {
		values.left = 1
	} else if defaultRight {
		values.right = 1
	}

	return sidesGroup{
		label:    label,
		values:   values,
		maxValue: 1,
		focused:  defaultFocus,
		styles:   s,
	}
}

func newSpikesGroup(s *styles, defaultFocus Side) sidesGroup {
	g := newSidesGroup(s, DISPLAY_SPIKES, defaultFocus, false, false)
	g.maxValue = int(models.Spikes3)
	return g
}

func (g sidesGroup) Update(msg tea.Msg) (sidesGroup, tea.Cmd) {
	if !g.groupFocused {
		return g, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, globalKeys.Left):
			g.focused = LeftSide
		case key.Matches(msg, globalKeys.Right):
			g.focused = RightSide
		case key.Matches(msg, globalKeys.Confirm):
			v := g.values.GetPtr(g.focused)
			// cycle through valid options (0 to 1: 0 -> 1 -> 0 | 0 to 3: 0 -> 1 -> 2 -> 3 -> 0)
			*v = (*v + 1) % (g.maxValue + 1)
		}
	}

	return g, nil
}

func (g sidesGroup) cell(v int) string {
	// boolean toggle
	if g.maxValue == 1 {
		if v == 1 {
			return CHECKED
		}

		return UNCHECKED
	}

	// spikes specific case
	if v == 0 {
		return UNCHECKED
	}

	return fmt.Sprintf("[%d]", v)
}

func (g sidesGroup) View() string {
	isGroupFocused := g.groupFocused

	// left
	leftCol := g.cell(g.values.left)

	leftStyle := lipgloss.NewStyle().Bold(isGroupFocused)
	if g.focused == LeftSide {
		leftStyle = leftStyle.Foreground(g.styles.highlight)
	}

	// label
	lableStyle := lipgloss.NewStyle().Bold(isGroupFocused)

	// right
	rightCol := g.cell(g.values.right)

	rightStyle := lipgloss.NewStyle().Bold(isGroupFocused)
	if g.focused == RightSide {
		rightStyle = rightStyle.Foreground(g.styles.highlight)
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

func (g sidesGroup) Level(side Side) int {
	return g.values.Get(side)
}

func (g sidesGroup) Value(side Side) bool {
	return g.Level(side) > 0
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

	title = FIELD_SIDES_TITLE

	for _, s := range m.sides {
		lines = append(lines, s.View())
	}

	return "\n" + titleStyle.Render(title) + ":\n" + lipgloss.JoinVertical(lipgloss.Left, lines...)
}
