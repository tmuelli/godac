package calc

import (
	"fmt"
	"slices"

	"github.com/tmuelli/godac/internal/models"
)

func metronomeMultiplier() float64 {
	// TODO: to be implemented => use calc context for same move count
	return 0.0
}

func (d *CalcData) validItemCondition(cond models.ItemConditionFlag, ctx *models.CalcContext) bool {
	if (cond & models.ItemConSuperEffective) != 0 {
		// TODO: use ctx here
	}

	if (cond & models.ItemConNotTransformed) != 0 {
		// only applies to ditto specie
		// TODO: for now assume ditto is always transformed when calculation is done
		// TODO: implement transformation info later
		if d.specie.ID == "ditto" {
			return false
		}
	}

	if (cond & models.ItemConPhysical) != 0 {
		if d.usedMove.Category != string(Physical) {
			return false
		}
	}

	if (cond & models.ItemConSpecial) != 0 {
		if d.usedMove.Category != string(Special) {
			return false
		}
	}

	if (cond & models.ItemConMetronome) != 0 {
		// Metronome is handled in special function as it needs
		// information about consecutive move use
		// see: metronomeModifier() and finalItemModifier()
		return true
	}

	if (cond & models.ItemConLifeOrb) != 0 {
		// TODO: to be implemented life orb
		// if recoil damage is considered later add it here
		return true
	}

	if (cond & models.ItemConTypeReduction) != 0 {
		// TODO: use ctx here
	}

	return true
}

func (d *CalcData) baseItemModifier(kind StatKind) float64 {
	if d.item == nil {
		// no item => nothing to do
		return 1.0
	}

	// verify item is effective for stat
	if len(d.item.Stats) == 0 {
		return 1.0
	}

	if !slices.Contains(d.item.Stats, string(kind)) {
		return 1.0
	}

	// check specie restrictions
	if d.item.Restrictions != nil {
		if len(d.item.Restrictions.Species) != 0 {
			if !slices.Contains(d.item.Restrictions.Species, d.specie.ID) {
				return 1.0
			}
		}
	}

	// check item conditions
	if d.item.Condition != nil {
		condFlag := d.item.Condition.ToFlag()
		if !d.validItemCondition(condFlag, nil) {
			return 1.0
		}
	}

	return d.item.Multiplier
}

// checks if calc data has item, if no ids are passed and item != nil it counts as true
func (d *CalcData) hasItem(ids ...string) bool {
	if d.item == nil {
		return false
	}

	if len(ids) == 0 {
		return true
	}

	return slices.Contains(ids, d.item.Id)
}

func (d *CalcData) moveItemModifier() (float64, error) {
	if d.item == nil {
		// no item => nothing to do
		return 1.0, nil
	}

	// verify item is effective for type
	if len(d.item.Types) != 0 {
		if !slices.Contains(d.item.Types, d.usedMove.Type) {
			return 1.0, nil
		}
	}

	// check item conditions
	if d.item.Condition != nil {
		condFlag := d.item.Condition.ToFlag()
		if !d.validItemCondition(condFlag, nil) {
			return 1.0, nil
		}
	}

	// check specie restrictions
	if d.item.Restrictions != nil {
		if len(d.item.Restrictions.Species) != 0 {
			if !slices.Contains(d.item.Restrictions.Species, d.specie.ID) {
				return 1.0, nil
			}
		}
	}

	return d.item.Multiplier, nil
}

func (d *CalcData) berryModifier(options models.ItemModifierOptions) (float64, error) {
	if options.MoveType == "" {
		return 1.0, fmt.Errorf("unknown move type for berry item mod calc")
	}

	if len(d.item.Types) == 0 {
		return 1.0, nil
	}

	// check matching berry damage reduction type
	if !slices.Contains(d.item.Types, options.MoveType) {
		return 1.0, nil
	}

	// check typeReduction flag
	if d.item.Condition == nil {
		return 1.0, nil
	}

	cond := d.item.Condition.ToFlag()
	if cond & models.ItemConTypeReduction == 0 {
		return 1.0, nil
	}

	// all type reducing berries can only reduce super effective moves (except chilan berry)
	if (cond & models.ItemConSuperEffective == 0) && options.TypeMod <= 1 {
		return 1.0, nil
	}

	return d.item.Multiplier, nil
}

func (d *CalcData) finalItemModifier(options models.ItemModifierOptions) (float64, error) {
	if d.item == nil {
		// no item => nothing to do
		return 1.0, nil
	}

	switch options.SubStage {
	case models.DamagePreLoop:
		// life orb and metronome only
		if d.item.IsBerry {
			return 1.0, nil
		}

		if d.item.Condition != nil {
			cond := d.item.Condition.ToFlag()

			if cond & models.ItemConMetronome != 0 {
				return metronomeMultiplier(), nil
			}

			if cond & models.ItemConLifeOrb != 0 {
				return d.item.Multiplier, nil
			}
		}

		return 1.0, nil

	case models.DamageExpertBelt:
		// expert belt => only on attacker side
		if d.dataType != Attacker || d.item.IsBerry {
			return 1.0, nil
		}

		if d.item.Condition != nil {
			cond := d.item.Condition.ToFlag()
			if cond & models.ItemConSuperEffective != 0 && options.TypeMod > 1 {
				return d.item.Multiplier, nil
			}
		}

		return 1.0, nil

	case models.DamageBerry:
		// Berry only — defender side
		if !d.item.IsBerry || (d.dataType != Defender) {
			return 1.0, nil
		}

		return d.berryModifier(options)

	default:
		return 1.0, fmt.Errorf("unknown item damage sub stage")
	}
}

func (d *CalcData) itemModifier(kind *StatKind, stage models.ItemEffect, options *models.ItemModifierOptions) (float64, error) {
	if d.item == nil {
		return 1.0, nil
	}
	
	if d.item.Effect != stage {
		return 1.0, nil
	}

	switch (d.item.Effect) {
	case models.StatEffect:
		if kind == nil {
			return 1.0, fmt.Errorf("unknown kind in base item mod calc")
		}

		return d.baseItemModifier(*kind), nil

	case models.MoveEffect:
		if kind == nil {
			return 1.0, fmt.Errorf("unknown kind in move item mod calc")
		}

		return d.moveItemModifier()

	case models.DamageEffect:
		if options == nil {
			return 1.0, fmt.Errorf("unknown options in final damage item mod calc")
		}

		return d.finalItemModifier(*options)

	default:
		return 0.0, fmt.Errorf("invalid effect type %v", d.item.Effect)
	}
}