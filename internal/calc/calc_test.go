package calc

import (
	"fmt"
	"os"
	"testing"

	"github.com/tmuelli/godac/internal/models"
)

type testCase struct {
	desc string // human readable, e.g. "Gengar Shadow Ball vs Blissey"

	// attacker
	atkSpecies   string
	atkLevel     int
	atkItem      string
	atkAbility   string // empty = first ability of the specie
	atkAbilityOn *bool  // nil = default state of the ability
	atkNature    string
	atkTypes     []string // type ids, empty = types of the specie
	atkIVs       models.IvStats
	atkEVs       models.EvStats
	atkBoosts    models.UserInputBoosts
	atkStatus    string
	atkHealth    int

	// move
	atkMove      string
	atkMoves     []string
	atkMoveCrit  bool
	moveOverride *models.UserMoveOverride
	hits         int // multihit move hit count, 0 = default
	defMove      string
	defMoves     []string

	// defender
	defSpecies   string
	defLevel     int
	defItem      string
	defAbility   string
	defAbilityOn *bool
	defNature    string
	defTypes     []string
	defIVs       models.IvStats
	defEVs       models.EvStats
	defBoosts    models.UserInputBoosts
	defHealth    int

	// field
	fieldType   string
	fieldLevel  int
	weather     string
	reflect     bool
	lightScreen bool

	// expected result — min and max of the 16-roll spread
	wantMin int
	wantMax int
}

// defaultLevel so you don't repeat 100 everywhere
const defaultLevel = 100

// evSetting maps test EVs to the input setting: all zero means no EVs,
// anything else is passed as custom values
func evSetting(evs models.EvStats) (string, *models.EvStats) {
	if evs == (models.EvStats{}) {
		return "worstCase", nil
	}
	return "custom", &evs
}

// ivSetting maps test IVs to the input setting: all zero means default 31 IVs,
// anything else is passed as custom values
func ivSetting(ivs models.IvStats) (string, *models.IvStats) {
	if ivs == (models.IvStats{}) {
		return "optimal", nil
	}
	return "custom", &ivs
}

// ev builds an EV spread, in order hp, atk, def, spa, spd, spe
func ev(hp, atk, def, spa, spd, spe int) models.EvStats {
	return models.EvStats{BaseStats: models.BaseStats{Hp: hp, Atk: atk, Def: def, Spa: spa, Spd: spd, Spe: spe}}
}

// calcFor builds the inputs from the test case and runs the calculation
func calcFor(t *testing.T, data *models.Data, tc testCase) *Calculation {
	t.Helper() // makes failures point to the call site, not this function

	atkLevel := tc.atkLevel
	if atkLevel == 0 {
		atkLevel = defaultLevel
	}
	defLevel := tc.defLevel
	if defLevel == 0 {
		defLevel = defaultLevel
	}

	atkEvSetting, atkEvCustom := evSetting(tc.atkEVs)
	atkIvSetting, atkIvCustom := ivSetting(tc.atkIVs)

	defEvSetting, defEvCustom := evSetting(tc.defEVs)
	defIvSetting, defIvCustom := ivSetting(tc.defIVs)

	var atkStatus *string
	if tc.atkStatus != "" {
		atkStatus = &tc.atkStatus
	}

	var hits *int
	if tc.hits != 0 {
		hits = &tc.hits
	}

	atkHp := tc.atkHealth
	if atkHp == 0 {
		atkHp = 100
	}

	defHp := tc.defHealth
	if defHp == 0 {
		defHp = 100
	}

	var atkItem *string
	if tc.atkItem != "" {
		atkItem = &tc.atkItem
	}

	var defItem *string
	if tc.defItem != "" {
		defItem = &tc.defItem
	}

	atkInput := &models.UserInput{
		Name:          tc.atkSpecies,
		Level:         atkLevel,
		Ability:       tc.atkAbility,
		AbilityOn:     tc.atkAbilityOn,
		Item:          atkItem,
		Nature:        tc.atkNature,
		Types:         tc.atkTypes,
		MoveOverride:  tc.moveOverride,
		EvSetting:     atkEvSetting,
		EvCustom:      atkEvCustom,
		IvSetting:     atkIvSetting,
		IvCustom:      atkIvCustom,
		Status:        atkStatus,
		Hits:          hits,
		Boosts:        tc.atkBoosts,
		AttackMove:    tc.atkMove,
		Moves:         tc.atkMoves,
		HealthPercent: atkHp,
		InputSide:     string(models.Attacker),
		MetronomeHits: nil,
	}

	defInput := &models.UserInput{
		Name:          tc.defSpecies,
		Level:         defLevel,
		Ability:       tc.defAbility,
		AbilityOn:     tc.defAbilityOn,
		Item:          defItem,
		Nature:        tc.defNature,
		Types:         tc.defTypes,
		EvSetting:     defEvSetting,
		EvCustom:      defEvCustom,
		IvSetting:     defIvSetting,
		IvCustom:      defIvCustom,
		Boosts:        tc.defBoosts,
		HealthPercent: defHp,
		AttackMove:    tc.defMove,
		Moves:         tc.defMoves,
		InputSide:     string(models.Defender),
		MetronomeHits: nil,
	}

	terrain := tc.weather
	if terrain == "" {
		terrain = string(models.Clear)
	}

	fieldType := tc.fieldType
	if fieldType == "" {
		fieldType = "Singles"
	}

	fieldLevel := tc.fieldLevel
	if fieldLevel == 0 {
		fieldLevel = 100
	}

	fieldInput := &models.UserFieldInput{
		Type:    fieldType,
		Level:   fieldLevel,
		Terrain: terrain,
		Right: models.UserFieldSideInput{
			Reflect:     tc.reflect,
			LightScreen: tc.lightScreen,
		},
	}

	calc, err := NewCalc(data, atkInput, defInput, fieldInput)
	if err != nil {
		t.Fatalf("[%s] NewCalc error: %v", tc.desc, err)
	}

	err = calc.Calc()
	if err != nil {
		t.Fatalf("[%s] Calc error: %v", tc.desc, err)
	}

	return calc
}

func runTest(t *testing.T, data *models.Data, tc testCase) {
	t.Helper()

	calc := calcFor(t, data, tc)

	if calc.Result.Range.Min != tc.wantMin || calc.Result.Range.Max != tc.wantMax {
		t.Errorf("[%s]\n  got  %d-%d\n  want %d-%d",
			tc.desc, calc.Result.Range.Min, calc.Result.Range.Max, tc.wantMin, tc.wantMax)
	}
}

var testData *models.Data

func TestMain(m *testing.M) {
	var err error
	testData, err = models.NewData()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load test data: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// helper used by test functions
func loadTestData(t *testing.T) *models.Data {
	t.Helper()
	if testData == nil {
		t.Fatal("test data not loaded")
	}
	return testData
}

func runTests(t *testing.T, tests []testCase) {
	t.Helper()
	data := loadTestData(t)

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			runTest(t, data, tc)
		})
	}
}

// Expected values in the tests below come from the Smogon damage calc (gen 4),
// run with the same inputs. Defaults match: level 100, 31 IVs, no EVs unless given.
// The attacker/defender ability is always the first ability of the specie.

func TestGen4BasicDamage(t *testing.T) {
	runTests(t, []testCase{
		{
			desc:       "Poochyena vs Poochyena",
			atkSpecies: "Poochyena", atkNature: "Hardy", atkMove: "Bite", atkMoves: []string{"Bite"},
			defSpecies: "Poochyena", defNature: "Hardy", defMove: "Bite", defMoves: []string{"Bite"},
			wantMin: 45, wantMax: 53,
		},
		{
			desc:       "Gengar Choice Specs +1 Focus Blast vs Chansey (smogon calc test)",
			atkSpecies: "Gengar", atkItem: "Choice Specs", atkNature: "Timid", atkEVs: ev(0, 0, 0, 252, 0, 0),
			atkBoosts: models.UserInputBoosts{Spa: 1}, atkMove: "Focus Blast", atkMoves: []string{"Focus Blast"},
			defSpecies: "Chansey", defItem: "Leftovers", defNature: "Calm", defEVs: ev(252, 0, 0, 0, 252, 0),
			wantMin: 408, wantMax: 482,
		},
		{
			desc:       "Garchomp Earthquake vs Metagross (STAB, super effective)",
			atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Earthquake", atkMoves: []string{"Earthquake"},
			defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 218, wantMax: 258,
		},
		{
			desc:       "Garchomp Outrage vs Scizor (STAB, resisted)",
			atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Outrage", atkMoves: []string{"Outrage"},
			defSpecies: "Scizor", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 78, wantMax: 92,
		},
		{
			desc:       "Nature: neutral Hardy instead of Adamant",
			atkSpecies: "Garchomp", atkNature: "Hardy", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Earthquake", atkMoves: []string{"Earthquake"},
			defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 198, wantMax: 234,
		},
		{
			desc:       "Level 50 vs level 50",
			atkSpecies: "Garchomp", atkLevel: 50, atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Earthquake", atkMoves: []string{"Earthquake"},
			defSpecies: "Metagross", defLevel: 50, defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 116, wantMax: 138,
		},
	})
}

func TestGen4Immunities(t *testing.T) {
	garchomp := testCase{
		atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0),
		atkMove: "Earthquake", atkMoves: []string{"Earthquake"},
	}
	gengar := testCase{
		atkSpecies: "Gengar", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0),
		atkMove: "Shadow Ball", atkMoves: []string{"Shadow Ball"},
	}
	zard := testCase{
		atkSpecies: "Charizard", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0),
		atkMove: "Flamethrower", atkMoves: []string{"Flamethrower"},
	}
	starmie := testCase{
		atkSpecies: "Starmie", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0),
		atkMove: "Surf", atkMoves: []string{"Surf"},
	}

	with := func(base testCase, desc, defSpecies, defNature string, defEVs models.EvStats) testCase {
		base.desc, base.defSpecies, base.defNature, base.defEVs = desc, defSpecies, defNature, defEVs
		return base
	}

	runTests(t, []testCase{
		with(garchomp, "Ground vs Levitate (Gengar)", "Gengar", "Hardy", models.EvStats{}),
		with(gengar, "Ghost vs Normal type (Chansey)", "Chansey", "Calm", ev(252, 0, 0, 0, 252, 0)),
		with(gengar, "Ghost vs Normal type (Snorlax, also Immunity)", "Snorlax", "Careful", ev(252, 0, 0, 0, 252, 0)),
		with(zard, "Fire vs Flash Fire (Heatran)", "Heatran", "Hardy", models.EvStats{}),
		with(starmie, "Water vs Water Absorb (Vaporeon)", "Vaporeon", "Hardy", models.EvStats{}),
	})
}

func TestGen4Abilities(t *testing.T) {
	runTests(t, []testCase{
		{
			desc:       "Pure Power (Medicham) Brick Break vs Snorlax",
			atkSpecies: "Medicham", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Brick Break", atkMoves: []string{"Brick Break"},
			defSpecies: "Snorlax", defNature: "Careful", defEVs: ev(252, 0, 0, 0, 252, 0),
			wantMin: 468, wantMax: 552,
		},
		{
			desc:       "Thick Fat (Dewgong) halves Flamethrower",
			atkSpecies: "Charizard", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Flamethrower", atkMoves: []string{"Flamethrower"},
			defSpecies: "Dewgong", defNature: "Hardy",
			wantMin: 78, wantMax: 92,
		},
	})
}

func TestGen4Items(t *testing.T) {
	garchomp := testCase{
		atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0),
	}
	metagross := testCase{
		defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
	}

	choiceBand := garchomp
	choiceBand.desc, choiceBand.atkItem = "Choice Band Earthquake vs Metagross", "Choice Band"
	choiceBand.atkMove, choiceBand.atkMoves = "Earthquake", []string{"Earthquake"}
	choiceBand.defSpecies, choiceBand.defNature, choiceBand.defEVs = metagross.defSpecies, metagross.defNature, metagross.defEVs
	choiceBand.wantMin, choiceBand.wantMax = 324, 384

	runTests(t, []testCase{
		choiceBand,
		{
			desc:       "Life Orb Shadow Ball vs Starmie",
			atkSpecies: "Gengar", atkItem: "Life Orb", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Shadow Ball", atkMoves: []string{"Shadow Ball"},
			defSpecies: "Starmie", defNature: "Bold", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 428, wantMax: 506,
		},
		{
			desc:       "Type boosting item (Charcoal) Flamethrower vs Metagross",
			atkSpecies: "Charizard", atkItem: "Charcoal", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Flamethrower", atkMoves: []string{"Flamethrower"},
			defSpecies: "Metagross", defNature: "Hardy",
			wantMin: 396, wantMax: 468,
		},
		{
			desc:       "Resist berry (Occa Berry) halves super effective Flamethrower",
			atkSpecies: "Charizard", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Flamethrower", atkMoves: []string{"Flamethrower"},
			defSpecies: "Metagross", defNature: "Hardy", defItem: "Occa Berry",
			wantMin: 165, wantMax: 195,
		},
	})
}

func TestGen4Weather(t *testing.T) {
	zard := testCase{
		atkSpecies: "Charizard", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Flamethrower", atkMoves: []string{"Flamethrower"},
		defSpecies: "Metagross", defNature: "Hardy",
	}
	starmie := testCase{
		atkSpecies: "Starmie", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Surf", atkMoves: []string{"Surf"},
		defSpecies: "Charizard", defNature: "Bold", defEVs: ev(252, 0, 0, 0, 252, 0),
	}

	cases := []struct {
		desc     string
		tc       testCase
		weather  string
		min, max int
	}{
		{"Fire move, no weather", zard, "", 330, 390},
		{"Fire move in sun", zard, "Sun", 492, 582},
		{"Fire move in rain", zard, "Rain", 168, 198},
		{"Water move in rain", starmie, "Rain", 372, 440},
		{"Water move in sun", starmie, "Sun", 126, 150},
	}

	tests := make([]testCase, 0, len(cases))
	for _, c := range cases {
		c.tc.desc, c.tc.weather, c.tc.wantMin, c.tc.wantMax = c.desc, c.weather, c.min, c.max
		tests = append(tests, c.tc)
	}
	runTests(t, tests)
}

func TestGen4FieldAndStatus(t *testing.T) {
	garchomp := testCase{
		atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Earthquake", atkMoves: []string{"Earthquake"},
		defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
	}

	reflect := garchomp
	reflect.desc, reflect.reflect, reflect.wantMin, reflect.wantMax = "Reflect halves physical damage", true, 110, 132

	burned := garchomp
	burned.desc, burned.atkStatus, burned.wantMin, burned.wantMax = "Burn halves physical damage", "burned", 110, 132

	doubles := garchomp
	doubles.desc, doubles.fieldType, doubles.wantMin, doubles.wantMax = "Doubles spread move reduction", "Doubles", 164, 194

	boosted := garchomp
	boosted.desc, boosted.atkBoosts, boosted.wantMin, boosted.wantMax = "+2 attack", models.UserInputBoosts{Atk: 2}, 432, 510

	lowered := garchomp
	lowered.desc, lowered.defBoosts, lowered.wantMin, lowered.wantMax = "-1 defense on defender", models.UserInputBoosts{Def: -1}, 324, 384

	lightScreen := testCase{
		desc:       "Light Screen halves special damage",
		atkSpecies: "Gengar", atkItem: "Choice Specs", atkNature: "Timid", atkEVs: ev(0, 0, 0, 252, 0, 0),
		atkBoosts: models.UserInputBoosts{Spa: 1}, atkMove: "Focus Blast", atkMoves: []string{"Focus Blast"},
		defSpecies: "Chansey", defItem: "Leftovers", defNature: "Calm", defEVs: ev(252, 0, 0, 0, 252, 0),
		lightScreen: true,
		wantMin:     204, wantMax: 242,
	}

	runTests(t, []testCase{reflect, burned, doubles, boosted, lowered, lightScreen})
}

func TestGen4FixedAndVariableDamage(t *testing.T) {
	// Garchomp / Machamp / Heatran / etc. with setup common per case
	runTests(t, []testCase{
		{
			desc:       "Seismic Toss deals level damage",
			atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Seismic Toss", atkMoves: []string{"Seismic Toss"},
			defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 100, wantMax: 100,
		},
		{
			desc:       "Dragon Rage deals 40",
			atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Dragon Rage", atkMoves: []string{"Dragon Rage"},
			defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 40, wantMax: 40,
		},
		{
			desc:       "Low Kick vs heavy target (Snorlax)",
			atkSpecies: "Machamp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Low Kick", atkMoves: []string{"Low Kick"},
			defSpecies: "Snorlax", defNature: "Careful", defEVs: ev(252, 0, 0, 0, 252, 0),
			wantMin: 612, wantMax: 722,
		},
		{
			desc:       "Low Kick vs lighter target (Metagross)",
			atkSpecies: "Machamp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Low Kick", atkMoves: []string{"Low Kick"},
			defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 128, wantMax: 152,
		},
		{
			desc:       "Grass Knot vs Metagross",
			atkSpecies: "Venusaur", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Grass Knot", atkMoves: []string{"Grass Knot"},
			defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 98, wantMax: 116,
		},
		{
			desc:       "Gyro Ball slower attacker",
			atkSpecies: "Metagross", atkNature: "Brave", atkEVs: ev(0, 252, 0, 0, 0, 0), atkMove: "Gyro Ball", atkMoves: []string{"Gyro Ball"},
			defSpecies: "Garchomp", defNature: "Adamant", defEVs: ev(0, 252, 0, 0, 0, 0),
			wantMin: 72, wantMax: 85,
		},
		{
			desc:       "Eruption at full health",
			atkSpecies: "Typhlosion", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Eruption", atkMoves: []string{"Eruption"},
			defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 522, wantMax: 614,
		},
		{
			// 50% of 297 HP is 148.5, the tui rounds to 149 HP => 75 BP (smogon calc with 149 current HP)
			desc:       "Eruption at 50% health",
			atkSpecies: "Typhlosion", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Eruption", atkMoves: []string{"Eruption"},
			atkHealth:  50,
			defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 260, wantMax: 308,
		},
	})
}

func TestGen4TripleKickHits(t *testing.T) {
	data := loadTestData(t)

	// the move data only allows 3 hits for triple kick, the smogon calc does not
	// escalate the base power per hit so the expectation here is the cartridge rule:
	// base power of 10 / 20 / 30 for hit 1 / 2 / 3
	const hits = 3

	tc := testCase{
		desc:       "Triple Kick with 3 hits",
		atkSpecies: "Hitmontop", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0),
		atkMove: "Triple Kick", atkMoves: []string{"Triple Kick"}, hits: hits,
		defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
	}

	calc := calcFor(t, data, tc)

	if got := len(calc.Result.DmgMatrix); got != hits {
		t.Fatalf("got %d damage rows, want %d", got, hits)
	}

	for i, row := range calc.Result.DmgMatrix {
		if len(row) != 16 {
			t.Fatalf("hit %d has %d rolls, want 16", i+1, len(row))
		}

		if i > 0 && row[15] <= calc.Result.DmgMatrix[i-1][15] {
			t.Errorf("hit %d max %d should be higher than hit %d max %d", i+1, row[15], i, calc.Result.DmgMatrix[i-1][15])
		}
	}

	// reported range is the first hit
	if first := calc.Result.DmgMatrix[0]; calc.Result.Range.Min != first[0] || calc.Result.Range.Max != first[15] {
		t.Errorf("range %d-%d should match first hit %d-%d", calc.Result.Range.Min, calc.Result.Range.Max, first[0], first[15])
	}
}

// Expected values for the ability cases come from the smogon calc (gen 4), except intimidate:
// smogon calc ignores intimidate in gen 4, so the expectation is the same attack with -1 attack.
func TestGen4AbilityInput(t *testing.T) {
	on, off := true, false

	// target for the toggle ability cases
	metagross := testCase{defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0)}

	flashFire := func(desc string, abilityOn *bool, min, max int) testCase {
		tc := metagross
		tc.desc = desc
		tc.atkSpecies, tc.atkNature, tc.atkEVs = "Heatran", "Modest", ev(0, 0, 0, 252, 0, 0)
		tc.atkMove, tc.atkMoves = "Eruption", []string{"Eruption"}
		tc.atkAbilityOn = abilityOn
		tc.wantMin, tc.wantMax = min, max
		return tc
	}

	plus := func(desc string, abilityOn *bool, min, max int) testCase {
		return testCase{
			desc:       desc,
			atkSpecies: "Plusle", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkAbilityOn: abilityOn,
			atkMove: "Thunderbolt", atkMoves: []string{"Thunderbolt"},
			defSpecies: "Starmie", defNature: "Bold", defEVs: ev(252, 0, 0, 0, 252, 0),
			wantMin: min, wantMax: max,
		}
	}

	slowStart := func(desc string, abilityOn *bool, min, max int) testCase {
		tc := metagross
		tc.desc = desc
		tc.atkSpecies, tc.atkNature, tc.atkEVs = "Regigigas", "Adamant", ev(0, 252, 0, 0, 0, 0)
		tc.atkMove, tc.atkMoves = "Facade", []string{"Facade"}
		tc.atkAbilityOn = abilityOn
		tc.wantMin, tc.wantMax = min, max
		return tc
	}

	// unburden doubles the speed of the attacker, which halves the base power of gyro ball
	unburden := func(desc string, abilityOn *bool, min, max int) testCase {
		return testCase{
			desc:       desc,
			atkSpecies: "Metagross", atkNature: "Brave", atkEVs: ev(0, 252, 0, 0, 0, 0), atkAbility: "Unburden", atkAbilityOn: abilityOn,
			atkMove: "Gyro Ball", atkMoves: []string{"Gyro Ball"},
			defSpecies: "Garchomp", defNature: "Adamant", defEVs: ev(0, 252, 0, 0, 0, 0),
			wantMin: min, wantMax: max,
		}
	}

	// the defender intimidates the attacker, which lowers the attack of the attacker by one stage
	intimidate := func(desc string, abilityOn *bool, min, max int) testCase {
		return testCase{
			desc:       desc,
			atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0),
			atkMove: "Earthquake", atkMoves: []string{"Earthquake"},
			defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0), defAbility: "Intimidate", defAbilityOn: abilityOn,
			wantMin: min, wantMax: max,
		}
	}

	runTests(t, []testCase{
		{
			desc:       "Technician boosts moves up to 60 base power",
			atkSpecies: "Scizor", atkAbility: "Technician", atkNature: "Hardy",
			atkMove: "Bug Bite", atkMoves: []string{"Bug Bite"},
			defSpecies: "Chansey", defNature: "Hardy",
			wantMin: 621, wantMax: 732,
		},
		{
			desc:       "Specie default ability is used without ability input",
			atkSpecies: "Scizor", atkNature: "Hardy",
			atkMove: "Bug Bite", atkMoves: []string{"Bug Bite"},
			defSpecies: "Chansey", defNature: "Hardy",
			wantMin: 415, wantMax: 489,
		},
		{
			desc:       "Pinsir Hyper Cutter Earthquake vs Levitate",
			atkSpecies: "Pinsir", atkItem: "Choice Band", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0),
			atkMove: "Earthquake", atkMoves: []string{"Earthquake"},
			defSpecies: "Gengar", defItem: "Choice Specs", defNature: "Timid", defEVs: ev(0, 0, 0, 252, 0, 0), defBoosts: models.UserInputBoosts{Spa: 1},
			wantMin: 0, wantMax: 0,
		},
		{
			desc:       "Mold Breaker ignores Levitate (smogon calc test)",
			atkSpecies: "Pinsir", atkItem: "Choice Band", atkAbility: "Mold Breaker", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0),
			atkMove: "Earthquake", atkMoves: []string{"Earthquake"},
			defSpecies: "Gengar", defItem: "Choice Specs", defNature: "Timid", defEVs: ev(0, 0, 0, 252, 0, 0), defBoosts: models.UserInputBoosts{Spa: 1},
			wantMin: 528, wantMax: 622,
		},
		flashFire("Flash Fire off by default", nil, 588, 692),
		flashFire("Flash Fire off", &off, 588, 692),
		flashFire("Flash Fire on", &on, 878, 1034),
		plus("Plus off by default", nil, 224, 266),
		plus("Plus on", &on, 338, 398),
		slowStart("Slow Start on by default", nil, 22, 27),
		slowStart("Slow Start off", &off, 44, 52),
		unburden("Unburden off by default", nil, 72, 85),
		unburden("Unburden on", &on, 36, 43),
		intimidate("Intimidate on by default", nil, 144, 170),
		intimidate("Intimidate on", &on, 144, 170),
		intimidate("Intimidate off", &off, 218, 258),
	})
}

// needs the type boost of plates in items.json (move effect with the plate type)
func TestGen4PlateWithTechnician(t *testing.T) {
	runTests(t, []testCase{
		{
			desc:       "Technician boost applies before the item boost (smogon calc test)",
			atkSpecies: "Scizor", atkItem: "Insect Plate", atkAbility: "Technician", atkNature: "Hardy",
			atkMove: "Bug Bite", atkMoves: []string{"Bug Bite"},
			defSpecies: "Chansey", defNature: "Hardy",
			wantMin: 745, wantMax: 877,
		},
	})
}

func TestAbilityInputValidation(t *testing.T) {
	data := loadTestData(t)

	input := func(ability string) *models.UserInput {
		return &models.UserInput{
			Name: "Metagross", Level: 100, Nature: "Hardy", Ability: ability,
			EvSetting: "worstCase", IvSetting: "optimal", HealthPercent: 100,
			InputSide: string(models.Defender),
		}
	}

	if _, err := validateInput(data, input("Technician")); err != nil {
		t.Errorf("known ability should be valid: %v", err)
	}

	if _, err := validateInput(data, input("")); err != nil {
		t.Errorf("empty ability should fall back to the specie ability: %v", err)
	}

	if _, err := validateInput(data, input("Not An Ability")); err == nil {
		t.Errorf("unknown ability should be invalid")
	}
}

// Expected values for the override cases come from the smogon calc (gen 4) with overrides set
// on the pokemon / move. Removing a type can not be compared: the smogon overrides are merged
// by index and can not drop the second type.
func TestGen4TypeOverrides(t *testing.T) {
	zard := testCase{
		atkSpecies: "Charizard", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0),
		atkMove: "Flamethrower", atkMoves: []string{"Flamethrower"},
		defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
	}

	bothTypes := zard
	bothTypes.desc, bothTypes.defTypes = "Both defender types replaced (Water/Fire)", []string{"water", "fire"}
	bothTypes.wantMin, bothTypes.wantMax = 41, 48

	garchomp := testCase{
		atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0),
		atkMove: "Flamethrower", atkMoves: []string{"Flamethrower"},
		defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
	}

	noStab := garchomp
	noStab.desc, noStab.wantMin, noStab.wantMax = "Fire move without Fire type has no STAB", 112, 134

	stab := garchomp
	stab.desc, stab.atkTypes = "Attacker with Fire as second type gets STAB", []string{"dragon", "fire"}
	stab.wantMin, stab.wantMax = 168, 200

	runTests(t, []testCase{
		bothTypes,
		noStab,
		stab,
		{
			desc:       "Defender with one type gets a second type",
			atkSpecies: "Starmie", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0), atkMove: "Surf", atkMoves: []string{"Surf"},
			defSpecies: "Chansey", defNature: "Hardy", defTypes: []string{"normal", "fire"},
			wantMin: 272, wantMax: 324,
		},
	})
}

func TestGen4MoveOverrides(t *testing.T) {
	bp, crit := 200, true

	earthquake := testCase{
		atkSpecies: "Garchomp", atkNature: "Adamant", atkEVs: ev(0, 252, 0, 0, 0, 0),
		atkMove: "Earthquake", atkMoves: []string{"Earthquake"},
		defSpecies: "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
	}

	override := func(desc string, o models.UserMoveOverride, min, max int) testCase {
		tc := earthquake
		tc.desc, tc.moveOverride, tc.wantMin, tc.wantMax = desc, &o, min, max
		return tc
	}

	runTests(t, []testCase{
		override("Base power 200 instead of 100", models.UserMoveOverride{BasePower: &bp}, 432, 510),
		override("Move type Fire instead of Ground (no STAB, super effective)", models.UserMoveOverride{Type: "fire"}, 146, 172),
		override("Critical hit", models.UserMoveOverride{Critical: &crit}, 438, 516),
		{
			desc:       "Shadow Ball as physical move uses Atk and Def",
			atkSpecies: "Gengar", atkNature: "Modest", atkEVs: ev(0, 0, 0, 252, 0, 0),
			atkMove: "Shadow Ball", atkMoves: []string{"Shadow Ball"},
			moveOverride: &models.UserMoveOverride{Category: "Physical"},
			defSpecies:   "Metagross", defNature: "Impish", defEVs: ev(252, 0, 252, 0, 0, 0),
			wantMin: 33, wantMax: 40,
		},
	})
}

func TestGetTypesFromIDs(t *testing.T) {
	data := loadTestData(t)

	for _, ids := range [][]string{{"fire"}, {"fire", "flying"}} {
		if _, err := getTypesFromIDs(data, ids); err != nil {
			t.Errorf("types %v should be valid: %v", ids, err)
		}
	}

	for _, ids := range [][]string{nil, {"nope"}, {"fire", "fire"}, {"fire", "water", "grass"}} {
		if _, err := getTypesFromIDs(data, ids); err == nil {
			t.Errorf("types %v should be invalid", ids)
		}
	}
}

func TestMoveOverrideValidation(t *testing.T) {
	data := loadTestData(t)

	tooHigh := models.MOVE_BASE_POWER_MAX + 1
	valid := 100

	input := func(o *models.UserMoveOverride) *models.UserInput {
		return &models.UserInput{
			Name: "Garchomp", Level: 100, Nature: "Hardy",
			AttackMove: "Earthquake", Moves: []string{"Earthquake"}, MoveOverride: o,
			EvSetting: "worstCase", IvSetting: "optimal", HealthPercent: 100,
			InputSide: string(models.Attacker),
		}
	}

	for desc, o := range map[string]*models.UserMoveOverride{
		"no override":  nil,
		"empty":        {},
		"base power":   {BasePower: &valid},
		"type":         {Type: "fire"},
		"category":     {Category: "Special"},
		"status moves": {Category: "Status"},
	} {
		if _, err := validateInput(data, input(o)); err != nil {
			t.Errorf("%s should be valid: %v", desc, err)
		}
	}

	for desc, o := range map[string]*models.UserMoveOverride{
		"base power out of range": {BasePower: &tooHigh},
		"unknown type":            {Type: "nope"},
		"unknown category":        {Category: "Mystic"},
	} {
		if _, err := validateInput(data, input(o)); err == nil {
			t.Errorf("%s should be invalid", desc)
		}
	}
}
