package tui

import (
	"fmt"
	"strconv"

	tea "charm.land/bubbletea/v2"
	"github.com/tmuelli/godac/internal/models"
)

const (
	NAME   = "Name"
	GENDER = "Gender"
	TYPE1  = "Type1"
	TYPE2  = "Type2"
	LEVEL  = "Level"
)

var basicsSpec = []fieldSpec{
	{
		key:     NAME,
		options: func(d *models.Data) []string { return d.SpecieNames() },
		onConfirm: func(data *list, value string, d *models.Data) tea.Cmd {
			specieData, ok := d.Species[value]
			if !ok {
				// should not happen => specie is validated through editConfirmMsg => if still unknown specie here, dont do anything => internal error
				return nil
			}

			fillBasicDefaults(data, &specieData, d)
			return requestAutoFillForSpecie(&specieData)
		},
	},
	// TODO: constants
	{key: GENDER, options: func(d *models.Data) []string { return []string{"Male", "Female"} }},
	{key: TYPE1, options: func(d *models.Data) []string { return d.TypeNames() }},
	{key: TYPE2, options: func(d *models.Data) []string { return d.TypeNames() }},
	{key: LEVEL, isNumeric: true, min: 1, max: 100, defaultValue: fmt.Sprintf("%d", models.Level100)},
}

func newBasicsModel(s *styles, data *models.Data) fieldListModel {
	m := newFieldListModel(s, data, basicsSpec)

	m.customUpdate = func(msg tea.Msg, data *list, d *models.Data) (tea.Cmd, bool) {
		if req, ok := msg.(setLevelRequestMsg); ok {
			data.SetValue(LEVEL, strconv.Itoa(req.level))
			return nil, true
		}
		return nil, false
	}

	return m
}

func fillBasicDefaults(data *list, specie *models.Specie, calcData *models.Data) {
	if specie == nil {
		// internal error
		return
	}

	for k := range data.values {
		if k == NAME {
			// specie name was already set in update from editConfirmMsg
			continue
		}

		var value string
		switch k {
		case GENDER:
			value = specie.Gender

			// if no value is given => default to female
			if value == "" {
				value = DEFAULT_GENDER
			}
		case TYPE1:
			value = calcData.TypeName(specie.Types[0])
		case TYPE2:
			if len(specie.Types) > 1 {
				value = calcData.TypeName(specie.Types[1])
			}
		case LEVEL:
			if data.values[LEVEL] != "" {
				// keep a level which was already set by the user (field was present before)
				continue
			}

			value = fmt.Sprintf("%d", models.Level100)
		}

		data.SetValue(k, value)
	}
}
