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

func (g sidesGroup) View(active bool) string {
	s := g.styles
	isCursorRow := active && g.groupFocused

	marker := UNSELECTED_MARKER
	label := lipgloss.NewStyle().Width(SIDES_LABEL_WIDTH)

	if isCursorRow {
		marker = s.selected.Render(SELECTED_MARKER)
		label = label.Bold(true)
	}

	cell := func(side Side) string {
		val := g.values.Get(side)

		style := s.dimText

		if val > 0 {
			style = s.okText
		}

		if isCursorRow && g.focused == side {
			style = s.selected
		}

		return style.Width(SIDE_COL_WIDTH).Render(g.cell(val))
	}

	return marker + label.Render(g.label) + cell(LeftSide) + cell(RightSide)
}

func (g sidesGroup) Level(side Side) int {
	return g.values.Get(side)
}

func (g sidesGroup) Value(side Side) bool {
	return g.Level(side) > 0
}
