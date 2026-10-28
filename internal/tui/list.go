package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
	"golang.org/x/exp/maps"
)

type list struct {
	values      map[string]string
	options     []string
	activeIndex int
	hidden      map[string]bool
}

func (l *list) SetValue(fieldID, value string) {
	l.values[fieldID] = value
}

func (l *list) SetActiveRow(rowIndex int) {
	i := l.activeIndex
	for {
		i += rowIndex

		if i < 0 || i >= len(l.options) {
			return
		}

		if !l.hidden[l.options[i]] {
			l.activeIndex = i
			return
		}
	}
}

func (l *list) SetHidden(fieldID string, hidden bool) {
	l.hidden[fieldID] = hidden

	if l.hidden[l.options[l.activeIndex]] {
		l.SetActiveRow(-1)
	}
}

func renderList(l list, s *styles) string {
	var lines []string
	var optionsIndexes []int
	maxLineWidth := 0

	// get max width and fill lines => hidden rows are not rendered
	for i, o := range l.options {
		if l.hidden[o] {
			continue
		}

		line := o + ": " + "[" + l.values[o] + "]"
		maxLineWidth = max(maxLineWidth, len(line))
		lines = append(lines, line)
		optionsIndexes = append(optionsIndexes, i)
	}

	var styledLines []string
	for i, line := range lines {
		style := lipgloss.NewStyle().Width(maxLineWidth + 20)
		if optionsIndexes[i] == l.activeIndex {
			style = style.Bold(true)
		}
		styledLines = append(styledLines, style.Render(line))
	}

	rows := lipgloss.JoinVertical(lipgloss.Left, styledLines...)
	return rows
}

type fieldSpec struct {
	key          string
	isNumeric    bool
	min          int
	max          int
	options      func(*models.Data) []string
	defaultValue string
	isToggle     bool
	onConfirm    func(data *list, value string, d *models.Data) tea.Cmd
}

type fieldListModel struct {
	fields       side[list]
	specs        []fieldSpec
	styles       *styles
	data         *models.Data
	customUpdate func(msg tea.Msg, data *list, d *models.Data) (tea.Cmd, bool)
}

func newFieldListModel(s *styles, data *models.Data, specs []fieldSpec) fieldListModel {
	template := map[string]string{}
	options := make([]string, len(specs))

	for i, spec := range specs {
		template[spec.key] = ""
		options[i] = spec.key
	}

	return fieldListModel{
		fields: side[list]{
			left:  list{values: maps.Clone(template), options: options, activeIndex: 0, hidden: map[string]bool{}},
			right: list{values: maps.Clone(template), options: options, activeIndex: 0, hidden: map[string]bool{}},
		},
		specs:  specs,
		styles: s,
		data:   data,
	}
}

func (m fieldListModel) specForKey(key string) fieldSpec {
	for _, s := range m.specs {
		if s.key == key {
			return s
		}
	}

	return fieldSpec{key: key}
}

func (m fieldListModel) Init() tea.Cmd {
	return nil
}

func (m fieldListModel) Update(msg tea.Msg, side Side) (tabContent, tea.Cmd) {
	data := m.fields.GetPtr(side)

	// check for custom updates on fieldListModel
	if m.customUpdate != nil {
		if cmd, isCustom := m.customUpdate(msg, data, m.data); isCustom {
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case editConfirmMsg:
		data.SetValue(msg.fieldID, msg.value)

		spec := m.specForKey(msg.fieldID)
		if spec.onConfirm != nil {
			return m, spec.onConfirm(data, msg.value, m.data)
		}

		return m, nil
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, globalKeys.Up):
			data.SetActiveRow(-1)
		case key.Matches(msg, globalKeys.Down):
			data.SetActiveRow(1)
		case key.Matches(msg, globalKeys.Confirm):
			key := data.options[data.activeIndex]

			if spec := m.specForKey(key); spec.isToggle {
				return m, m.toggle(data, spec)
			}

			value := data.values[key]
			return m, m.requestInput(key, value)
		}
	}

	return m, nil
}

func (m fieldListModel) View(side Side) string {
	return renderList(m.fields.Get(side), m.styles)
}

func (m fieldListModel) toggle(data *list, spec fieldSpec) tea.Cmd {
	next := ""
	switch data.values[spec.key] {
	case DISPLAY_YES:
		next = DISPLAY_NO
	case DISPLAY_NO:
		next = DISPLAY_YES
	default:
		return nil
	}

	data.SetValue(spec.key, next)

	if spec.onConfirm != nil {
		return spec.onConfirm(data, next, m.data)
	}

	return nil
}

func (m fieldListModel) requestInput(key, value string) tea.Cmd {
	spec := m.specForKey(key)

	if spec.isNumeric {
		return requestNumberInput(key, value, spec.min, spec.max)
	}

	if spec.options != nil {
		return requestSelectInput(key, value, spec.options(m.data))
	}

	return requestTextInput(key, value)
}
