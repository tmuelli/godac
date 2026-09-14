package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
	"github.com/tmuelli/godac/internal/utils"
	"golang.org/x/exp/maps"
)

type miscModel struct {
	values		side[map[string]string]
	options		[]string
	activeIndex	int
	styles		*styles
	data		*models.Data
}

func newMiscModel(s *styles, data *models.Data) miscModel {
	template := map[string]string{
		"Nature": "",
		"Ability": "",
		"Item": "",
		"Status": "",
		"Health": "",
	}

	return miscModel{
		values:		side[map[string]string]{
			left:	maps.Clone(template),
			right:	maps.Clone(template),
		},
		options:		maps.Keys(template),
		activeIndex:	0,
		styles:			s,
		data:			data,
	}
}

func (m miscModel) Init() tea.Cmd {
	return nil
}

func (m miscModel) Update(msg tea.Msg, side Side) (tabContent, tea.Cmd) {
	switch msg := msg.(type) {
	case editConfirmMsg:
		m.values.Get(side)[msg.fieldID] = msg.value
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.activeIndex = max(m.activeIndex - 1, 0)
		case "down":
			m.activeIndex = min(m.activeIndex + 1, len(m.options) - 1)
		case "enter":
			key := m.options[m.activeIndex]
			value := m.values.Get(side)[key]
			return m, m.requestInput(key, value)
		}
	}

	return m, nil
}

func (m miscModel) View(side Side) string {
	values := m.values.Get(side)

	var lines []string

	for i, o := range m.options {
		style := lipgloss.NewStyle()
		if i == m.activeIndex {
			style = style.Bold(true)
		}
		line := o + ": " + "[" + values[o] + "]"
		lines = append(lines, style.Render(line))
	}

	rows := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return rows
}

// TODO: maybe make general requestInput
func (m miscModel) requestInput(key, value string) tea.Cmd {
	min := 0
	max := 0
	var options []string
	isNumeric := false

	// TODO: maybe use SliceFROmMapValue somewhere before to not do this everytime again
	// TODO: change to constants where its reasonable
	switch key {
	case "Nature":
		options = utils.SliceFromMapValue(m.data.Natures, func(n models.Nature) string {
			return n.Name
		})
	case "Ability":
		options = utils.SliceFromMapValue(m.data.Abilities, func(a models.Ability) string {
			return a.Name
		})
	case "Item":
		options = utils.SliceFromMapValue(m.data.Items, func(i models.Item) string {
			return i.Name
		})
	case "Status":
		options = []string{"Healthy", "Poisoned", "Badly Poisoned", "Burned", "Paralyzed", "Asleep", "Frozen"}
	case "Health":
		isNumeric = true
		min = 1
		max = 100
	}

	if isNumeric {
		return requestNumberInput(key, value, min, max)
	} else {
		return requestSelectInput(key, value, options)
	}
}