package models

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

var testData *Data

func TestMain(m *testing.M) {
	var err error
	testData, err = NewData()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load test data: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func ptr[T any](v T) *T { return &v }

// only what differs between species; everything else comes from wantFor
type defaultsCase struct {
	species   string
	ability   string
	gender    string
	abilityOn *bool
	types     []string
}

// wantFor builds the full expected input: the shared baseline plus the row's values
func wantFor(tc defaultsCase) UserInput {
	baseStats := testData.Species[tc.species].BaseStats

	return UserInput{
		Name:              tc.species,
		Level:             int(Level100),
		Nature:            "Hardy",
		HealthPercent:     100,
		IvSetting:         "custom",
		IvCustom:          &IvStats{BaseStats: BaseStats{Hp: 31, Atk: 31, Def: 31, Spa: 31, Spd: 31, Spe: 31}},
		EvSetting:         "custom",
		EvCustom:          &EvStats{},
		Ability:           tc.ability,
		AbilityOn:         tc.abilityOn,
		Status:            ptr("healthy"),
		Gender:            ptr(tc.gender),
		Types:             tc.types,
		BaseStatsOverride: &baseStats,
	}
}

func jsonOf(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b)
}

func TestDefaultUserInput(t *testing.T) {
	tests := []defaultsCase{
		{species: "Poochyena", ability: "Run Away", gender: "F", types: []string{"dark"}},
		{species: "Gyarados", ability: "Intimidate", gender: "F", abilityOn: ptr(true), types: []string{"water", "flying"}},
		{species: "Regigigas", ability: "Slow Start", gender: "N", abilityOn: ptr(true), types: []string{"normal"}},
		{species: "Flareon", ability: "Flash Fire", gender: "F", abilityOn: ptr(false), types: []string{"fire"}},
		{species: "Magnemite", ability: "Magnet Pull", gender: "N", types: []string{"electric", "steel"}},
		{species: "Nidoking", ability: "Poison Point", gender: "M", types: []string{"poison", "ground"}},
		{species: "Chansey", ability: "Natural Cure", gender: "F", types: []string{"normal"}},
		{species: "Mr. Mime", ability: "Soundproof", gender: "F", types: []string{"psychic"}},
	}

	for _, tc := range tests {
		t.Run(tc.species, func(t *testing.T) {
			specie, ok := testData.Species[tc.species]
			if !ok {
				t.Fatalf("%s missing from species data", tc.species)
			}

			got := testData.DefaultUserInput(&specie)
			want := wantFor(tc)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("\n got: %s\nwant: %s", jsonOf(t, got), jsonOf(t, want))
			}
		})
	}
}

func TestDefaultUserInput_NilSpecie(t *testing.T) {
	got := testData.DefaultUserInput(nil)
	want := UserInput{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n got: %s\nwant: %s", jsonOf(t, got), jsonOf(t, want))
	}
}

func TestDefaultUserInput_NoSharedState(t *testing.T) {
	tc := defaultsCase{species: "Poochyena", ability: "Run Away", gender: "F", types: []string{"dark"}}

	specie, ok := testData.Species[tc.species]
	if !ok {
		t.Fatalf("%s missing from species data", tc.species)
	}

	firstInput := testData.DefaultUserInput(&specie)
	if firstInput.IvCustom == nil || firstInput.Gender == nil || firstInput.BaseStatsOverride == nil {
		t.Fatal("IvCustom and Gender must be set for this test")
	}

	if len(firstInput.Types) != 1 || (len(firstInput.Types) > 0 && firstInput.Types[0] != "dark") {
		t.Fatal("Incorrect types")
	}

	firstInput.IvCustom.Atk = 0
	*firstInput.Gender = "M"
	firstInput.Types[0] = "water"
	firstInput.BaseStatsOverride.Atk = 1

	secondInput := testData.DefaultUserInput(&specie)

	if len(secondInput.Types) == 0 {
		t.Fatal("Invalid types")
	}

	want := wantFor(tc)

	if !reflect.DeepEqual(secondInput, want) {
		t.Errorf("\n got: %s\nwant: %s", jsonOf(t, secondInput), jsonOf(t, want))
	}
}

func TestDefaultUserInput_AllSpecies(t *testing.T) {
	for name, specie := range testData.Species {
		t.Run(name, func(t *testing.T) {
			got := testData.DefaultUserInput(&specie)

			if got.Name != name {
				t.Errorf("Name = %q, want %q", got.Name, name)
			}

			if _, ok := testData.Abilities[got.Ability]; !ok {
				t.Errorf("Unknown ability: %q", got.Ability)
			}

			if _, ok := testData.Natures[got.Nature]; !ok {
				t.Errorf("Unknown nature: %q", got.Nature)
			}

			if got.Gender == nil {
				t.Error("gender is nil")
			} else {
				if *got.Gender != "M" && *got.Gender != "F" && *got.Gender != "N" {
					t.Errorf("invalid gender %q", *got.Gender)
				}
				if specie.Gender != "" && *got.Gender != specie.Gender {
					t.Errorf("gender = %q, but species is fixed to %q", *got.Gender, specie.Gender)
				}
			}

			if !reflect.DeepEqual(got.Types, specie.Types) {
				t.Errorf("\n got: %s\nwant: %s", jsonOf(t, got.Types), jsonOf(t, specie.Types))
			}

			if got.BaseStatsOverride == nil {
				t.Error("BaseStatsOverride is nil")
			} else if !reflect.DeepEqual(*got.BaseStatsOverride, specie.BaseStats) {
				t.Errorf("\n got: %s\nwant: %s", jsonOf(t, got.BaseStatsOverride), jsonOf(t, specie.BaseStats))
			}
		})
	}
}
