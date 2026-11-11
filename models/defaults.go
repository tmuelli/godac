package models

import "slices"

// DefaultUserInput returns the starting input for a species.
func (d *Data) DefaultUserInput(s *Specie) UserInput {
	if s == nil {
		return UserInput{}
	}

	abilityName := s.Abilities["0"]

	var abilityOn *bool
	if ability, ok := d.Abilities[abilityName]; ok {
		if defaultOn, hasToggle := TogglableAbilities[ability.Id]; hasToggle {
			abilityOn = &defaultOn
		}
	}

	status := "healthy"

	g := s.Gender
	if g == "" {
		g = "F"
	}

	baseStats := s.BaseStats

	return UserInput{
		Name:              s.Name,
		Level:             int(Level100),
		Nature:            DEFAULT_NATURE,
		HealthPercent:     100,
		Ability:           abilityName,
		AbilityOn:         abilityOn,
		Status:            &status,
		Gender:            &g,
		Types:             slices.Clone(s.Types),
		BaseStatsOverride: &baseStats,
		IvSetting:         "custom",
		EvSetting:         "custom",
		IvCustom: &IvStats{
			BaseStats: BaseStats{
				Hp:  IV_MAX,
				Atk: IV_MAX,
				Def: IV_MAX,
				Spa: IV_MAX,
				Spd: IV_MAX,
				Spe: IV_MAX,
			},
		},
		EvCustom: &EvStats{},
	}
}

// DefaultMoveInput returns the starting values for a move.
func (d *Data) DefaultMoveInput(m *Move) MoveDefaults {
	if m == nil {
		return MoveDefaults{}
	}

	basePower := m.BasePower
	critical := m.Critical

	hits := 1
	if len(m.Multihit) > 0 {
		hits = slices.Min(m.Multihit)
	}

	return MoveDefaults{
		Override: &UserMoveOverride{
			BasePower: &basePower,
			Type:      m.Type,
			Category:  m.Category,
			Critical:  &critical,
		},
		Hits:          &hits,
		MetronomeHits: new(0),
	}
}
