package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmuelli/godac/internal/models"
	"golang.org/x/exp/maps"
)

type list struct {
	values      map[string]string
	options     []string
	activeIndex int
}

func (l *list) SetValue(fieldID, value string) {
	l.values[fieldID] = value
}

func (l *list) SetActiveRow(rowIndex int) {
	l.activeIndex = max(0, min(l.activeIndex+rowIndex, len(l.options)-1))
}

func renderList(l list, s *styles) string {
	var lines []string
	maxLineWidth := 0

	// get max width and fill lines
	for _, o := range l.options {
		line := o + ": " + "[" + l.values[o] + "]"
		maxLineWidth = max(maxLineWidth, len(line))
		lines = append(lines, line)
	}

	var styledLines []string
	for i, line := range lines {
		style := lipgloss.NewStyle().Width(maxLineWidth + 20)
		if i == l.activeIndex {
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
			left:  list{values: maps.Clone(template), options: options, activeIndex: 0},
			right: list{values: maps.Clone(template), options: options, activeIndex: 0},
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
		switch msg.String() {
		case "up":
			data.SetActiveRow(-1)
		case "down":
			data.SetActiveRow(1)
		case "enter":
			key := data.options[data.activeIndex]
			value := data.values[key]
			return m, m.requestInput(key, value)
		}
	}

	return m, nil
}

func (m fieldListModel) View(side Side) string {
	return renderList(m.fields.Get(side), m.styles)
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
