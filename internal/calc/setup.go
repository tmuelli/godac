package calc

import (
	"fmt"
	"strings"

	"github.com/tmuelli/godac/internal/models"
)

func (c *Calculation) isWeatherActive() bool {
	if (c.attacker.hasAbility("airlock", "cloudnine")) ||
		(c.defender.hasAbility("airlock", "cloudnine")) {
		return true
	}

	return false
}

func (c *Calculation) handleForecast(side *CalcData) error {
	if side.specie.ID == "castform" && side.hasAbility("forecast") {
		side.specieTypes = nil
		var new string

		switch c.field.Terrain {
		case models.Sun, models.Rain, models.Hail, models.Clear:
			new = getTypeStringFromWeather(c.field.Terrain)
		default:
			new = "Normal"
		}

		newType, ok := c.data.Types[new]
		if !ok {
			return fmt.Errorf("invalid forecast type")
		}

		side.specieTypes = append(side.specieTypes, &newType)
	}

	return nil
}

func (c *Calculation) handleItem(source *CalcData, target *CalcData) error {
	blockAbilities := []string{"clearbody", "whitesmoke", "hypercutter"}

	if source.hasAbility("intimidate") && !target.hasAbility(blockAbilities...) {
		if target.hasAbility("simple") {
			target.boosts.Atk.Set(-2)
		} else {
			target.boosts.Atk.Set(-1)
		}
	}

	return nil
}

func (c *Calculation) handleDownload(source *CalcData, target *CalcData) error {
	if source.hasAbility("download") {
		def := target.specie.BaseStats.Def
		spd := target.specie.BaseStats.Spd

		if spd <= def {
			source.boosts.Spa.Set(1)
		} else {
			source.boosts.Atk.Set(1)
		}
	}

	return nil
}

func (c *Calculation) handleSpecialMoveSetup() {
	switch c.attacker.usedMove.ID {
	case "weatherball":
		c.attacker.usedMove.Type = getTypeStringFromWeather(c.field.Terrain)
	case "judgement":
		if c.attacker.hasItem() && strings.HasSuffix(c.attacker.item.Id, "plate") {
			c.attacker.usedMove.Type = c.attacker.item.GetPlateBoostType()
		}
	case "naturalgift":
		if c.attacker.hasItem() && c.attacker.item.IsBerry {
			c.attacker.usedMove.Type = c.attacker.item.NaturalGift.Type
			c.attacker.usedMove.BasePower = c.attacker.item.NaturalGift.BasePower
		}
	case "brickbreak":
		c.field.Right.LightScreen = false
		c.field.Right.Reflect = false
	}
}

func (c *Calculation) setup() error {
	// check if weather is active => e.g. abilities air lock or cloud nine active
	if !c.isWeatherActive() {
		c.field.Terrain = models.Clear
	}

	// check if ability of defender is active
	if c.defender.ability.IsIgnoredBy(c.attacker.ability.Id) {
		c.defender.ability.Active = false
	}

	// check forecast
	err := c.handleForecast(c.attacker)
	if err != nil {
		return err
	}

	err = c.handleForecast(c.defender)
	if err != nil {
		return err
	}

	// check powertrick
	if c.field.Left.PowerTrick {
		c.attacker.specie.BaseStats.Atk, c.attacker.specie.BaseStats.Def = c.attacker.specie.BaseStats.Def, c.attacker.specie.BaseStats.Atk
	}

	if c.field.Right.PowerTrick {
		c.defender.specie.BaseStats.Atk, c.defender.specie.BaseStats.Def = c.defender.specie.BaseStats.Def, c.defender.specie.BaseStats.Atk
	}

	// check item
	err = c.handleItem(c.attacker, c.defender)
	if err != nil {
		return err
	}

	err = c.handleItem(c.defender, c.attacker)
	if err != nil {
		return err
	}

	// check download
	err = c.handleDownload(c.attacker, c.defender)
	if err != nil {
		return err
	}

	err = c.handleDownload(c.defender, c.attacker)
	if err != nil {
		return err
	}

	// check critical hit
	preventCritAbilities := []string{"battlearmor", "shellarmor"}
	if c.attacker.usedMove.Critical && c.defender.hasAbility(preventCritAbilities...) {
		c.attacker.usedMove.Critical = false
	}

	// handle special move setups
	if c.attacker.usedMove.NeedSpecialMoveSetup() {
		c.handleSpecialMoveSetup()
	}

	// check normalize
	if c.attacker.hasAbility("normalize") && c.attacker.usedMove.ID != "struggle" {
		c.attacker.usedMove.Type = "Normal"
	}

	return nil
}
