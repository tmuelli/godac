package calc

import (
	"fmt"
	"os"
	"testing"

	"github.com/tmuelli/godac/internal/models"
)

type testCase struct {
    desc        string // human readable, e.g. "Gengar Shadow Ball vs Blissey"

    // attacker
    atkSpecies  string
    atkLevel    int
    atkItem     string
    atkNature   string
	atkIVs      models.IvStats
	atkEVs      models.EvStats
    atkBoosts   models.UserInputBoosts
    atkStatus   string
	atkHealth	int

    // move
    atkMove        string
	atkMoves		[]string
    atkMoveCrit    bool
	defMove			string
	defMoves		[]string

    // defender
    defSpecies  string
    defLevel    int
    defItem     string
    defNature   string
	defIVs      models.IvStats
	defEVs      models.EvStats
	defBoosts   models.UserInputBoosts
	defHealth	int

    // field
	fieldType   string
	fieldLevel  int
    weather     string
    reflect     bool
    lightScreen bool

    // expected result — min and max of the 16-roll spread
    wantMin     int
    wantMax     int
}

// defaultLevel so you don't repeat 100 everywhere
const defaultLevel = 100

func runTest(t *testing.T, data *models.Data, tc testCase) {
    t.Helper() // makes failures point to the call site, not this function

    atkLevel := tc.atkLevel
    if atkLevel == 0 {
        atkLevel = defaultLevel
    }
    defLevel := tc.defLevel
    if defLevel == 0 {
        defLevel = defaultLevel
    }

	atkEvSetting := "worstCase"
	atkIvSetting := "optimal"

	defEvSetting := "worstCase"
	defIvSetting := "optimal"

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
        Name:			tc.atkSpecies,
        Level:			atkLevel,
        Item:			atkItem,
        Nature:			tc.atkNature,
        EvSetting:		atkEvSetting,
        IvSetting:		atkIvSetting,
        Boosts:			tc.atkBoosts,
        AttackMove:		tc.atkMove,
		Moves:			tc.atkMoves,
		HealthPercent:  atkHp,
    }

    defInput := &models.UserInput{
        Name:			tc.defSpecies,
        Level:			defLevel,
        Item:			defItem,
        Nature:			tc.defNature,
        EvSetting:		defEvSetting,
        IvSetting:		defIvSetting,
		Boosts:			tc.defBoosts,
		HealthPercent:  defHp,
		AttackMove:		tc.defMove,
		Moves:			tc.defMoves,
    }

	terrain := tc.weather
	if terrain == "" {
		terrain = "Clear"
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
		Type: fieldType,
		Level: fieldLevel,
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

func TestGen4BasicDamage(t *testing.T) {
    data := loadTestData(t) // see step 5

    tests := []testCase{
        {
            desc:       "Poochyena vs Poochyena",
            atkSpecies: "Poochyena",
            atkNature:  "Hardy",
            atkMove:       "Bite",
			atkMoves:		[]string{"Bite"},
            defSpecies: "Poochyena",
            defNature:  "Hardy",
            defMove:       "Bite",
			defMoves:		[]string{"Bite"},
            wantMin:    45,
            wantMax:    53,
        },
        // add more cases here
    }

    for _, tc := range tests {
        t.Run(tc.desc, func(t *testing.T) {
            runTest(t, data, tc)
        })
    }
}
