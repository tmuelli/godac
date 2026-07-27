package calc

import (
	"fmt"
	"slices"
	"strings"

	"github.com/tmuelli/godac/internal/models"
)

func (c *Calculation) isWeatherActive() bool {
	if (c.attacker.ability.Id == "airlock" || c.attacker.ability.Id == "cloudnine") ||
	   (c.defender.ability.Id == "airlock" || c.defender.ability.Id == "cloudnine") {
		return true
	}

	return false
}

func (c *Calculation) handleForecast(side *CalcData) error {
	if side.specie.ID == "castform" && side.ability.Id == "forecast" {
		side.specieTypes = nil
		new := "Normal"
		switch (c.field.Terrain) {
		case models.Sun:
			new = "Fire"
		case models.Rain:
			new = "Water"
		case models.Hail:
			new = "Ice"
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

	if source.ability.Id == "intimidate" && !slices.Contains(blockAbilities, target.ability.Id) {
		if target.ability.Id == "simple" {
			target.boosts.Atk.Set(-2)
		} else {
			target.boosts.Atk.Set(-1)
		}
	}
	
	return nil
}

func (c *Calculation) handleDownload(source *CalcData, target *CalcData) error {
	if source.ability.Id == "download" {
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
		// TODO: maybe summarize with forecast function with exception for sand case => check if thats worth it
		switch c.field.Terrain {
		case models.Sun:
			c.attacker.usedMove.Type = "Fire"
		case models.Rain:
			c.attacker.usedMove.Type = "Water"
		case models.Hail:
			c.attacker.usedMove.Type = "Ice"
		case models.Sand:
			c.attacker.usedMove.Type = "Rock"
		default:
			c.attacker.usedMove.Type = "Normal"
		}
	case "judgement":
		if c.attacker.item != nil && strings.HasSuffix(c.attacker.item.Id, "plate") {
			c.attacker.usedMove.Type = c.attacker.item.GetPlateBoostType()
		}
	case "naturalgift":
		if c.attacker.item != nil && c.attacker.item.IsBerry {
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
	
	// check forecast
	err := c.handleForecast(c.attacker)
	if err != nil {
		return err
	}

	err = c.handleForecast(c.defender)
	if err != nil {
		return err
	}

	// TODO: checkItem implementation => ability Klutz = except iron ball in gen 4 and not ev items

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

	// check if ability of defender is active
	// TODO: change IsActive receiver to not reflect Active member logic => more like is effective or switch logic => is ability overwridden / deactivated
	if !c.defender.ability.IsActive(c.attacker.ability.Id) {
		c.defender.ability.Active = false
	}

	// check critical hit
	preventCritAbilities := []string{"battlearmor", "shellarmor"}
	if c.attacker.usedMove.Critical && c.defender.ability.Active && slices.Contains(preventCritAbilities, c.defender.ability.Id) {
		c.attacker.usedMove.Critical = false
	}

	// handle special move setups
	if c.attacker.usedMove.NeedSpecialMoveSetup() {
		c.handleSpecialMoveSetup()
	}

	// check normalize
	if c.attacker.ability.Active && c.attacker.ability.Id == "normalize" && c.attacker.usedMove.ID != "struggle" {
		c.attacker.usedMove.Type = "Normal"
	}

	return nil
}