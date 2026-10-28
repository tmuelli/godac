package calc

import (
	"fmt"
	"math"

	"github.com/tmuelli/godac/internal/models"
)

func (s *StatusCondition) IsValid() bool {
	condMap := map[string]StatusCondition{
		"none":           NoneStatus,
		"healthy":        Healthy,
		"poisoned":       Poisoned,
		"badly poisoned": BadlyPoisoned,
		"burned":         Burned,
		"paralyzed":      Paralyzed,
		"asleep":         Asleep,
		"frozen":         Frozen,
	}

	con, ok := condMap[string(*s)]
	if !ok {
		return false
	}

	*s = con
	return true
}

func validateFieldSideInput(input *models.UserFieldSideInput) (*models.FieldSide, error) {
	if input == nil {
		return nil, fmt.Errorf("internal error: Invalid pointer to user field side input")
	}

	// field side spikes stack
	spikes := models.FieldSpikesStack(input.Spikes)
	if !spikes.IsValid() {
		return nil, fmt.Errorf("invalid field side input: Unknown spikes stack type %d", input.Spikes)
	}

	// successfully validated side input
	side := models.FieldSide{
		Spikes:       spikes,
		StealthRocks: input.StealthRocks,
		Reflect:      input.Reflect,
		LightScreen:  input.LightScreen,
		Protect:      input.Protect,
		LeechSeed:    input.LightScreen,
		Foresight:    input.Foresight,
		HelpingHand:  input.HelpingHand,
		Tailwind:     input.Tailwind,
		FlowerGift:   input.FlowerGift,
		PowerTrick:   input.PowerTrick,
		SwitchingOut: input.SwitchingOut,
	}
	return &side, nil
}

func validateFieldInput(input *models.UserFieldInput) (*models.Field, error) {
	if input == nil {
		return nil, fmt.Errorf("internal error: Invalid pointer to user field input")
	}

	// field type
	ft := models.FieldType(input.Type)
	if !ft.IsValid() {
		return nil, fmt.Errorf("invalid field input: Unknown field type %s", input.Type)
	}

	// field level
	fl := models.FieldLevel(input.Level)
	if !fl.IsValid() {
		return nil, fmt.Errorf("invalid field input: Invalid field level %v", input.Level)
	}

	// field terrain
	t := models.FieldTerrain(input.Terrain)
	if !t.IsValid() {
		return nil, fmt.Errorf("invalid field input: Unknown field terrain %s", input.Terrain)
	}

	// left side
	left, err := validateFieldSideInput(&input.Left)
	if err != nil {
		return nil, fmt.Errorf("invalid left field side input: %v", err)
	}

	// right side
	right, err := validateFieldSideInput(&input.Right)
	if err != nil {
		return nil, fmt.Errorf("invalid right field side input: %v", err)
	}

	// successfully validated field input
	field := models.Field{
		Type:      ft,
		Level:     fl,
		Terrain:   t,
		IsGravity: input.IsGravity,
		Left:      *left,
		Right:     *right,
	}

	return &field, nil
}

func validateInput(data *models.Data, input *models.UserInput) (*CalcData, error) {
	if data == nil {
		return nil, fmt.Errorf("internal error: Invalid base data")
	}

	if input == nil {
		return nil, fmt.Errorf("internal error: Invalid pointer to user input")
	}

	// input side
	var calcType models.CalcDataType
	switch input.InputSide {
	case string(models.Attacker):
		calcType = models.Attacker
	case string(models.Defender):
		calcType = models.Defender
	default:
		return nil, fmt.Errorf("invalid input: Unknown input side %s", input.InputSide)
	}

	// specie
	var specie *models.Specie

	if s, ok := data.Species[input.Name]; ok {
		specie = &s
	}

	if specie == nil {
		return nil, fmt.Errorf("invalid input: Unknown specie %s", input.Name)
	}

	// specie gender
	var gender Gender
	if input.Gender != nil {
		if !gender.isValid(input.Gender) {
			return nil, fmt.Errorf("invalid input: Unknown gender %s", *input.Gender)
		}

		if !gender.isValidForSpecie(specie) {
			return nil, fmt.Errorf("invalid input: gender %s is not valid for specie %s", *input.Gender, specie.ID)
		}
	}

	// nature
	var nature *models.Nature

	if n, ok := data.Natures[input.Nature]; ok {
		nature = &n
	}

	if nature == nil {
		return nil, fmt.Errorf("invalid input: Unknown nature %s", input.Nature)
	}

	// moves
	var move *string
	var usedMove *models.Move

	if calcType == models.Attacker {
		for i, m := range input.Moves {
			if m == input.AttackMove {
				move = &input.Moves[i]
				break
			}
		}

		if move == nil {
			return nil, fmt.Errorf("invalid attack Move: Move is %s is not in moveset of attacker", input.AttackMove)
		}

		if m, ok := data.Moves[*move]; ok {
			usedMove = &m
		}

		if usedMove == nil {
			return nil, fmt.Errorf("invalid move used: Unknown move")
		}

		// check for multiple hits or multihit move
		usedMove.Hits = DEFAULT_HIT_COUNT

		if input.Hits != nil {
			if *input.Hits == 0 {
				return nil, fmt.Errorf("invalid hits value: Must hit at least once")
			}

			if len(usedMove.Multihit) != 0 {
				if !usedMove.ValidMultihit(*input.Hits) {
					return nil, fmt.Errorf("invalid hits value: Must be value in multihit range")
				}
			}

			usedMove.Hits = *input.Hits
		}
	}

	// level
	if input.Level < 1 || input.Level > 100 {
		return nil, fmt.Errorf("invalid input: Invalid level %d", input.Level)
	}

	// health
	if input.HealthPercent < 1 || input.HealthPercent > 100 {
		return nil, fmt.Errorf("invalid input: Invalid health %d", input.HealthPercent)
	}

	// iv settings
	var ivSettings *models.IvStats
	switch IvOption(input.IvSetting) {
	case IvOptimal:
		ivSettings = &models.IvStats{
			BaseStats: models.BaseStats{
				Hp:  31,
				Atk: 31,
				Def: 31,
				Spa: 31,
				Spd: 31,
				Spe: 31,
			},
		}
	case IvWorstCase:
		ivSettings = &models.IvStats{
			BaseStats: models.BaseStats{
				Hp:  0,
				Atk: 0,
				Def: 0,
				Spa: 0,
				Spd: 0,
				Spe: 0,
			},
		}
	case IvCustom:
		if input.IvCustom == nil {
			return nil, fmt.Errorf("invalid input: Custom iv settings only supported when custom values are given")
		}

		ivSettings = input.IvCustom
		if !ivSettings.Valid() {
			return nil, fmt.Errorf("invalid input: Invalid custom iv values")
		}
	default:
		return nil, fmt.Errorf("invalid input: Unknown iv setting")
	}

	// ev settings
	var evSettings *models.EvStats
	switch EvOption(input.EvSetting) {
	case EvOptimal:
		evSettings = &models.EvStats{
			BaseStats: models.BaseStats{
				Hp:  31,
				Atk: 31,
				Def: 31,
				Spa: 31,
				Spd: 31,
				Spe: 31,
			},
		}
	case EvWorstCase:
		evSettings = &models.EvStats{
			BaseStats: models.BaseStats{
				Hp:  0,
				Atk: 0,
				Def: 0,
				Spa: 0,
				Spd: 0,
				Spe: 0,
			},
		}
	case EvCustom:
		if input.EvCustom == nil {
			return nil, fmt.Errorf("invalid input: Custom ev settings only supported when custom values are given")
		}

		evSettings = input.EvCustom
		if !evSettings.Valid() {
			return nil, fmt.Errorf("invalid input: Invalid custom ev values")
		}
	default:
		return nil, fmt.Errorf("invalid input: Unknown ev setting")
	}

	// specie type
	var specieTypes []*models.Type
	specieTypes, err := getTypesFromSpecie(data, specie)
	if err != nil {
		return nil, fmt.Errorf("internal error: Couldnt get type to specie")
	}

	// specie ability
	var ability *models.Ability
	ability, err = getAbilityFromSpecie(data, specie)
	if err != nil {
		return nil, fmt.Errorf("internal error: Couldnt get ability to specie")
	}
	// initial abilities are always active
	ability.Active = true

	// item
	var item *models.Item
	var metronomeHits int

	if input.Item != nil {
		if i, ok := data.Items[*input.Item]; ok {
			item = &i
		}

		if item == nil {
			return nil, fmt.Errorf("invalid input: Unknown item")
		}

		// metronome hits
		if item.Id == "metronome" {
			if input.MetronomeHits == nil {
				return nil, fmt.Errorf("invalid input: metronome hits invalid")
			}

			if *input.MetronomeHits < 0 || *input.MetronomeHits > 10 {
				return nil, fmt.Errorf("invalid input: metronome hits value is out of range")
			}

			metronomeHits = *input.MetronomeHits
		}
	}

	// status condition
	statusCon := NoneStatus

	if input.Status != nil {
		statusCon = StatusCondition(*input.Status)
		if !statusCon.IsValid() {
			return nil, fmt.Errorf("invalid input: Unknown status condition")
		}
	}

	// boosts
	boosts := models.NewBoosts(&input.Boosts)

	// successful validation of user input
	calcData := CalcData{
		specie:        specie,
		specieTypes:   specieTypes,
		level:         input.Level,
		percentHealth: input.HealthPercent,
		currentHealth: 0,
		maxHealth:     0,
		nature:        nature,
		ability:       ability,
		item:          item,
		usedMove:      usedMove,
		ivs:           ivSettings,
		evs:           evSettings,
		status:        statusCon,
		boosts:        &boosts,
		gender:        &gender,
		dataType:      calcType,
		metrHits:      metronomeHits,
	}

	// absolute health from percentage
	calcData.maxHealth = calcData.effectiveHpStat()
	calcData.currentHealth = max(1, int(math.Round(float64(calcData.percentHealth)*float64(calcData.maxHealth)/100)))

	return &calcData, nil
}
