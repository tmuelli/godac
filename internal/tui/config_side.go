package tui

import (
	"strconv"

	"github.com/tmuelli/godac/internal/models"
)

type side[T any] struct {
	left	T
	right	T
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
	basics := m.tabContent[0].(basicsModel)
	return basics.values.Get(side)["Name"]
}

func (m configModel) Gender(side Side) *string {
	basics := m.tabContent[0].(basicsModel)
	gender := basics.values.Get(side)["Gender"]
	return &gender
}

func (m configModel) Level(side Side) int {
	basics := m.tabContent[0].(basicsModel)

	// level
	level, err := strconv.Atoi(basics.values.Get(side)["Level"])
	// TODO: error handling
	if err != nil {
		
	}
	return level
}

// TODO: setting types / using types from input

func (m configModel) Moves(side Side) []string {
	moves := m.tabContent[3].(movesModel)
	return moves.Names(side)
}

func (m configModel) Hits(side Side, move string) *int {
	moves := m.tabContent[3].(movesModel)
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
	misc := m.tabContent[2].(miscModel)

	health, err := strconv.Atoi(misc.values.Get(side)["Health"])
	// TODO: error handling
	if err != nil {
		
	}
	return health
}

func (m configModel) IvCustomStats(side Side) *models.IvStats {
	stats := m.tabContent[1].(statsModel)
	data := stats.tab.Get(side)

	s := models.IvStats{}

	for _, row := range data.rows {
		// TODO: put into constant to generalize with ev logic
		v := row[2]

		value, err := strconv.Atoi(v)
		// TODO: error handling
		if err != nil {

		}

		switch row[0] {
		case "Hp":
			s.Hp = value
		case "Atk":
			s.Atk = value
		case "Def":
			s.Def = value
		case "Spa":
			s.Spa = value
		case "Spd":
			s.Spd = value
		case "Spe":
			s.Spe = value
		}
	}

	return &s
}

func (m configModel) EvCustomStats(side Side) *models.EvStats {
	stats := m.tabContent[1].(statsModel)
	data := stats.tab.Get(side)

	s := models.EvStats{}

	for _, row := range data.rows {
		// TODO: put into constant to generalize with ev logic
		v := row[3]

		value, err := strconv.Atoi(v)
		// TODO: error handling
		if err != nil {

		}

		switch row[0] {
		case "Hp":
			s.Hp = value
		case "Atk":
			s.Atk = value
		case "Def":
			s.Def = value
		case "Spa":
			s.Spa = value
		case "Spd":
			s.Spd = value
		case "Spe":
			s.Spe = value
		}
	}

	return &s
}

func (m configModel) Item(side Side) *string {
	misc := m.tabContent[2].(miscModel)
	item := misc.values.Get(side)["Item"]
	return &item
}

func (m configModel) Status(side Side) *string {
	misc := m.tabContent[2].(miscModel)
	status := misc.values.Get(side)["Status"]
	return &status
}

func (m configModel) Boosts(side Side) models.UserInputBoosts {
	stats := m.tabContent[1].(statsModel)
	data := stats.tab.Get(side)
	boosts := models.UserInputBoosts{}
	
	for _, row := range data.rows{
		// TODO: use constant
		v := row[5]

		value, err := strconv.Atoi(v)
		// TODO: error handling
		if err != nil {

		}

		switch row[0] {
		case "Hp":
			// TODO: remove from view / selectable disable
		case "Atk":
			boosts.Atk = value
		case "Def":
			boosts.Def = value
		case "Spa":
			boosts.Spa = value
		case "Spd":
			boosts.Spd = value
		case "Spe":
			boosts.Spe = value
		}
	}

	return boosts
}