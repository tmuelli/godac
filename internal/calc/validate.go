package calc

import (
	"fmt"

	"github.com/tmuelli/godac/internal/models"
)

func (s *StatusCondition) IsValid() bool {
	condMap := map[string]StatusCondition {
		"none":				NoneStatus,
		"healthy":			Healthy,
		"poisoned":			Poisoned,
		"badly poisoned":	BadlyPoisoned,
		"burned":			Burned,
		"paralyzed":		Paralyzed,
		"asleep":			Asleep,
		"frozen":			Frozen,
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
		Spikes: spikes,
		StealthRocks: input.StealthRocks,
		Reflect: input.Reflect,
		LightScreen: input.LightScreen,
		Protect: input.Protect,
		LeechSeed: input.LightScreen,
		Foresight: input.Foresight,
		HelpingHand: input.HelpingHand,
		Tailwind: input.Tailwind,
		FlowerGift: input.FlowerGift,
		PowerTrick: input.PowerTrick,
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
		Type: ft,
		Level: fl,
		Terrain: t,
		IsGravity: input.IsGravity,
		Left: *left,
		Right: *right,
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
	
	// specie
	var specie *models.Specie

	if s, ok := data.Species[input.Name]; ok {
		specie = &s
	}

	if specie == nil {
		return nil, fmt.Errorf("invalid input: Unknown specie %s", input.Name)
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
	for i, m := range input.Moves {
		if m == input.AttackMove {
			move = &input.Moves[i]
			break
		}
	}

	if move == nil {
		return nil, fmt.Errorf("invalid attack Move: Move is %s is not in moveset of attacker", input.AttackMove)
	}

	var usedMove *models.Move

	if m, ok := data.Moves[*move]; ok {
		usedMove = &m
	}

	if usedMove == nil {
		return nil, fmt.Errorf("invalid move used: Unknown move")
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
	switch (IvOption(input.IvSetting)) {
	case IvOptimal:
		ivSettings = &models.IvStats{
			Hp: 31,
			Atk: 31,
			Def: 31,
			Spa: 31,
			Spd: 31,
			Spe: 31,
		}
	case IvWorstCase:
		ivSettings = &models.IvStats{
			Hp: 0,
			Atk: 0,
			Def: 0,
			Spa: 0,
			Spd: 0,
			Spe: 0,
		}
	case IvCustom:
		return nil, fmt.Errorf("invalid input: Custom iv setting not supported yet")
	default:
		return nil, fmt.Errorf("invalid input: Unknown iv setting")
	}

	// ev settings
	var evSettings *models.EvStats
	switch (EvOption(input.EvSetting)) {
	case EvOptimal:
		evSettings = &models.EvStats{
			Hp: 31,
			Atk: 31,
			Def: 31,
			Spa: 31,
			Spd: 31,
			Spe: 31,
		}
	case EvWorstCase:
		evSettings = &models.EvStats{
			Hp: 0,
			Atk: 0,
			Def: 0,
			Spa: 0,
			Spd: 0,
			Spe: 0,
		}
	case EvCustom:
		return nil, fmt.Errorf("invalid input: Custom ev setting not supported yet")
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

	// item
	var item *models.Item

	if input.Item != nil {
		if i, ok := data.Items[*input.Item]; ok {
			item = &i
		}

		if item == nil {
			return nil, fmt.Errorf("invalid input: Unknown item")
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

	// successful validation of user input
	calcData := CalcData{
		specie: specie,
		specieTypes: specieTypes,
		level: input.Level,
		health: input.HealthPercent,
		nature: nature,
		ability: ability,
		item: item,
		usedMove: usedMove,
		ivs: ivSettings,
		evs: evSettings,
		status: statusCon,
		boosts: &models.Boosts{},
	}
	return &calcData, nil
}