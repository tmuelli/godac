package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
)

type viewTab struct {
	label string
	view  View
}

const (
	HORIZONTAL_PADDING            = 4
	HEADER_SEPERATOR              = " | "
	HEADER_ATTACKER_SHORTCUT_INFO = " (l)"
	HEADER_DEFENDER_SHORTCUT_INFO = " (r)"
)

var viewTabs = []viewTab{
	{"Calc", CalcView},
	{"Config", ConfigView},
	{"Field", FieldView},
}

func (m model) header() string {
	s := m.styles

	w := DEFAULT_WIDTH
	if m.width > 0 {
		// subtract horizontal padding
		w = max(m.width-HORIZONTAL_PADDING, DEFAULT_MIN_WIDTH)
	}

	var labels []string
	for _, t := range viewTabs {
		if t.view == m.current {
			labels = append(labels, s.title.Underline(true).Render(t.label))
		} else {
			labels = append(labels, s.dimText.Render(t.label))
		}
	}

	sep := lipgloss.NewStyle().Foreground(s.border).Render(HEADER_SEPERATOR)

	left := lipgloss.NewStyle().Bold(true).Render(APP_NAME) + "   " + strings.Join(labels, sep)

	right := ""
	// config view additionally displays which side is currently selected
	if m.current == ConfigView {
		if m.configSide == LeftSide {
			right = s.dimText.Render(string(models.Attacker) + HEADER_ATTACKER_SHORTCUT_INFO)
		} else {
			right = s.dimText.Render(string(models.Defender) + HEADER_DEFENDER_SHORTCUT_INFO)
		}
	}

	gap := max(1, w-lipgloss.Width(left)-lipgloss.Width(right))
	line := left + strings.Repeat(" ", gap) + right
	rule := lipgloss.NewStyle().Foreground(s.border).Render(strings.Repeat("─", w))

	return line + "\n" + rule
}
