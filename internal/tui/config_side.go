package tui

import (
	"github.com/tmuelli/godac/internal/models"
)

type side[T any] struct {
	left	T
	right	T
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
	// TODO: maybe iota constant for consistent tab accessing
	basics := m.tabContent[basicsTab].(basicsModel)
	return basics.fieldList.Get(side).values["Name"]
}

func (m configModel) Gender(side Side) *string {
	basics := m.tabContent[basicsTab].(basicsModel)
	gender := basics.fieldList.Get(side).values["Gender"]
	return &gender
}

func (m configModel) Level(side Side) int {
	basics := m.tabContent[basicsTab].(basicsModel)

	// level
	level := atoiOrZero(basics.fieldList.Get(side).values["Level"])
	return level
}

// TODO: setting types / using types from input

func (m configModel) Moves(side Side) []string {
	moves := m.tabContent[movesTab].(movesModel)
	return moves.Names(side)
}

func (m configModel) Hits(side Side, move string) *int {
	moves := m.tabContent[movesTab].(movesModel)
	data := moves.tab.Get(side)

	var hits int

	for _, m := range data.rows {
		// TODO: use constants here
		if m[0] == move {
			switch m[5] {
			case "Once", "1 hits":
				hits = 1
			case "Twice", "2 hits":
				hits = 2
			case "3 times", "3 hits":
				hits = 3
			case "4 times", "4 hits":
				hits = 4
			case "5 times", "5 hits":
				hits = 5
			}
		}
	}

	return &hits
}

func (m configModel) HealthPercent(side Side) int {
	misc := m.tabContent[miscTab].(miscModel)
	health := atoiOrZero(misc.fieldList.Get(side).values["Health"])
	return health
}

func (m configModel) Item(side Side) *string {
	misc := m.tabContent[miscTab].(miscModel)
	item := misc.fieldList.Get(side).values["Item"]
	return &item
}

func (m configModel) Status(side Side) *string {
	misc := m.tabContent[miscTab].(miscModel)
	status := misc.fieldList.Get(side).values["Status"]
	return &status
}

func (m configModel) IvCustomStats(side Side) *models.IvStats {
	stats := m.tabContent[statsTab].(statsModel)
	values := stats.columnValues(side, 2)
	s := models.IvStats{
		BaseStats: models.BaseStats{
			Hp:		values["Hp"],
			Atk:	values["Atk"],
			Def:	values["Def"],
			Spa:	values["Spa"],
			Spd:	values["Spd"],
			Spe:	values["Spe"],
		},
	}

	return &s
}

func (m configModel) EvCustomStats(side Side) *models.EvStats {
	stats := m.tabContent[statsTab].(statsModel)
	values := stats.columnValues(side, 3)
	s := models.EvStats{
		BaseStats: models.BaseStats{
			Hp:		values["Hp"],
			Atk:	values["Atk"],
			Def:	values["Def"],
			Spa:	values["Spa"],
			Spd:	values["Spd"],
			Spe:	values["Spe"],
		},
	}

	return &s
}

func (m configModel) Boosts(side Side) models.UserInputBoosts {
	stats := m.tabContent[statsTab].(statsModel)
	values := stats.columnValues(side, 5)
	s := models.UserInputBoosts{
		Atk:	values["Atk"],
		Def:	values["Def"],
		Spa:	values["Spa"],
		Spd:	values["Spd"],
		Spe:	values["Spe"],
	}

	return s
}