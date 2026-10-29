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
	// label width of columns only for visible columns
	labelWidth := 0
	for _, o := range l.options {
		if !l.hidden[o] {
			labelWidth = max(labelWidth, lipgloss.Width(o))
		}
	}
	labelWidth += COLUMN_GAP

	var rows []string
	for i, o := range l.options {
		if l.hidden[o] {
			continue
		}

		cursor := i == l.activeIndex
		marker := UNSELECTED_MARKER
		labelStyle := lipgloss.NewStyle()
		valueStyle := s.valueStyle(l.values[o])

		if cursor {
			marker = s.selected.Render(SELECTED_MARKER)
			labelStyle = s.selected
			valueStyle = s.selected
		}

		value := l.values[o]
		if value == "" {
			value = EMPTY_VALUE
		}

		rows = append(rows, marker+labelStyle.Width(labelWidth).Render(o)+valueStyle.Render(value))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
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
