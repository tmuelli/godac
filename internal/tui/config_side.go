package tui

import (
	"strings"

	"github.com/tmuelli/godac/internal/models"
)

type side[T any] struct {
	left  T
	right T
}

func (s *side[T]) GetPtr(side Side) *T {
	if side == LeftSide {
		return &s.left
	}
	return &s.right
}

func (s side[T]) Get(side Side) T {
	if side == LeftSide {
		return s.left
	}
	return s.right
}

func (s *side[T]) Set(side Side, v T) {
	if side == LeftSide {
		s.left = v
	} else {
		s.right = v
	}
}

func (m configModel) Name(side Side) string {
	basics := m.tabContent[basicsTab].(fieldListModel)
	return basics.fields.Get(side).values[NAME]
}

func (m configModel) Gender(side Side) *string {
	basics := m.tabContent[basicsTab].(fieldListModel)
	gender := basics.fields.Get(side).values[GENDER]

	var result string
	switch gender {
	case DISPLAY_MALE_NAME:
		result = "M"
	case DISPLAY_FEMALE_NAME:
		result = "F"
	case DISPLAY_GENDERLESS_NAME:
		result = "N"
	}

	return &result
}

func (m configModel) Level(side Side) int {
	basics := m.tabContent[basicsTab].(fieldListModel)

	// level
	level := atoiOrZero(basics.fields.Get(side).values[LEVEL])
	return level
}

// TODO: setting types / using types from input

func (m configModel) Moves(side Side) []string {
	moves := m.tabContent[movesTab].(movesModel)
	settings := m.MoveSettings(side)
	return moves.Names(side, settings)
}

func (m configModel) MoveSettings(side Side) []moveRowSettings {
	moves := m.tabContent[movesTab].(movesModel)
	return moves.MoveRowSettings(side)
}

func (m configModel) Hits(side Side, move string) *int {
	moves := m.tabContent[movesTab].(movesModel)
	data := moves.tab.Get(side)

	var hits int

	for _, row := range data.rows {
		if row[MOVE_COL] == move {
			hits = hitCount(hitOptions, row[HITS_COL])
		}
	}

	return &hits
}

func (m configModel) MetronomeHits(side Side, move string) *int {
	moves := m.tabContent[movesTab].(movesModel)
	data := moves.tab.Get(side)

	var metrHits int

	if data.columns[METR_HITS_COL].visible {
		for _, row := range data.rows {
			if row[MOVE_COL] == move {
				metrHits = hitCount(metronomeHitOptions, row[METR_HITS_COL])
			}
		}
	}

	return &metrHits
}

func (m configModel) HealthPercent(side Side) int {
	misc := m.tabContent[miscTab].(fieldListModel)
	// remove percent sign from string value
	healthString := misc.fields.Get(side).values[HEALTH_PERCENT]
	health := atoiOrZero(healthString[:len(healthString)-1])
	return health
}

func (m configModel) Item(side Side) *string {
	misc := m.tabContent[miscTab].(fieldListModel)
	item := misc.fields.Get(side).values[ITEM]

	if item == "" {
		return nil
	}

	return &item
}

func (m configModel) Status(side Side) *string {
	misc := m.tabContent[miscTab].(fieldListModel)
	status := misc.fields.Get(side).values[STATUS]
	status = strings.ToLower(status)
	return &status
}

func (m configModel) IvCustomStats(side Side) *models.IvStats {
	stats := m.tabContent[statsTab].(statsModel)
	values := stats.columnValues(side, IV_COL)
	s := models.IvStats{
		BaseStats: models.BaseStats{
			Hp:  values[DISPLAY_HP],
			Atk: values[DISPLAY_ATK],
			Def: values[DISPLAY_DEF],
			Spa: values[DISPLAY_SPA],
			Spd: values[DISPLAY_SPD],
			Spe: values[DISPLAY_SPE],
		},
	}

	return &s
}

func (m configModel) EvCustomStats(side Side) *models.EvStats {
	stats := m.tabContent[statsTab].(statsModel)
	values := stats.columnValues(side, EV_COL)
	s := models.EvStats{
		BaseStats: models.BaseStats{
			Hp:  values[DISPLAY_HP],
			Atk: values[DISPLAY_ATK],
			Def: values[DISPLAY_DEF],
			Spa: values[DISPLAY_SPA],
			Spd: values[DISPLAY_SPD],
			Spe: values[DISPLAY_SPE],
		},
	}

	return &s
}

func (m configModel) Boosts(side Side) models.UserInputBoosts {
	stats := m.tabContent[statsTab].(statsModel)
	values := stats.columnValues(side, BOOST_COL)
	s := models.UserInputBoosts{
		Atk: values[DISPLAY_ATK],
		Def: values[DISPLAY_DEF],
		Spa: values[DISPLAY_SPA],
		Spd: values[DISPLAY_SPD],
		Spe: values[DISPLAY_SPE],
	}

	return s
}

func (m configModel) Nature(side Side) string {
	misc := m.tabContent[miscTab].(fieldListModel)
	return misc.fields.Get(side).values[NATURE]
}

func (m configModel) InputSide(side Side) string {
	if side == LeftSide {
		return string(models.Attacker)
	}
	return string(models.Defender)
}
