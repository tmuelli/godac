package models

import (
	"reflect"
	"slices"
	"testing"
)

type defaultsMoveCase struct {
	move      string
	basePower int
	moveType  string
	category  string
	hits      int
}

// wantMoveFor builds the full expected input: the shared baseline plus the row's values
func wantMoveFor(tc defaultsMoveCase) MoveDefaults {
	return MoveDefaults{
		Override: &UserMoveOverride{
			BasePower: new(tc.basePower),
			Type:      tc.moveType,
			Category:  tc.category,
			Critical:  new(false), // default for critical is always false
		},
		Hits:          new(tc.hits),
		MetronomeHits: new(0), // this equals "Never"
	}
}

func TestDefaultMoveInput(t *testing.T) {
	tests := []defaultsMoveCase{
		{move: "Bite", basePower: 60, moveType: "dark", category: "Physical", hits: 1},
		{move: "Fury Swipes", basePower: 18, moveType: "normal", category: "Physical", hits: 2},
		{move: "Double Kick", basePower: 30, moveType: "fighting", category: "Physical", hits: 2},
		{move: "Triple Kick", basePower: 10, moveType: "fighting", category: "Physical", hits: 3},
		{move: "Swords Dance", basePower: 0, moveType: "normal", category: "Status", hits: 1},
		{move: "Seismic Toss", basePower: 0, moveType: "fighting", category: "Physical", hits: 1},
	}

	for _, tc := range tests {
		t.Run(tc.move, func(t *testing.T) {
			move, ok := testData.Moves[tc.move]
			if !ok {
				t.Fatalf("%s missing from moves data", tc.move)
			}

			got := testData.DefaultMoveInput(&move)
			want := wantMoveFor(tc)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("\n got: %s\nwant: %s", jsonOf(t, got), jsonOf(t, want))
			}
		})
	}
}

func TestDefaultMoveInput_NilMove(t *testing.T) {
	got := testData.DefaultMoveInput(nil)
	want := MoveDefaults{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n got: %s\nwant: %s", jsonOf(t, got), jsonOf(t, want))
	}
}

func TestDefaultMoveInput_NoSharedState(t *testing.T) {
	tc := defaultsMoveCase{move: "Bite", basePower: 60, moveType: "dark", category: "Physical", hits: 1}

	move, ok := testData.Moves[tc.move]
	if !ok {
		t.Fatalf("%s missing from moves data", tc.move)
	}

	firstInput := testData.DefaultMoveInput(&move)
	if firstInput.Override == nil || firstInput.Override.BasePower == nil || firstInput.Hits == nil || firstInput.Override.Critical == nil || firstInput.MetronomeHits == nil {
		t.Fatal("Override, BasePower, Hits, Metronome Hits and Critical must be set for this test")
	}

	*firstInput.Override.BasePower = 70
	*firstInput.Override.Critical = true
	*firstInput.Hits = 2
	*firstInput.MetronomeHits = 3

	secondInput := testData.DefaultMoveInput(&move)

	want := wantMoveFor(tc)

	if !reflect.DeepEqual(secondInput, want) {
		t.Errorf("\n got: %s\nwant: %s", jsonOf(t, secondInput), jsonOf(t, want))
	}
}

func TestDefaultMoveInput_AllMoves(t *testing.T) {
	for name, move := range testData.Moves {
		t.Run(name, func(t *testing.T) {
			got := testData.DefaultMoveInput(&move)

			if got.Override == nil {
				t.Fatal("Override must be set")
			}

			if got.Override.BasePower == nil || got.Override.Critical == nil || got.Hits == nil || got.MetronomeHits == nil {
				t.Fatal("BasePower, Critical, Hits and MetronomeHits must be set")
			}

			if _, ok := testData.Types[got.Override.Type]; !ok {
				t.Errorf("Unknown move type: %q", got.Override.Type)
			}

			if got.Override.Category == "" {
				t.Error("Move category is empty")
			} else {
				cat := got.Override.Category
				if cat != "Physical" && cat != "Special" && cat != "Status" {
					t.Errorf("invalid move category %q", cat)
				}
			}

			if *got.Hits < 1 || *got.Hits > 5 {
				t.Errorf("invalid hits range %d", *got.Hits)
			}

			if len(move.Multihit) > 0 {
				if want := slices.Min(move.Multihit); *got.Hits != want {
					t.Errorf("hits = %d, want %d (lowest multihit)", *got.Hits, want)
				}
			} else if *got.Hits != 1 {
				t.Errorf("hits = %d, want 1 (not a multihit move)", *got.Hits)
			}

			if *got.Override.BasePower != move.BasePower {
				t.Errorf("Invalid base power, got %d want %d", *got.Override.BasePower, move.BasePower)
			}

			if *got.Override.BasePower < MOVE_BASE_POWER_MIN || *got.Override.BasePower > MOVE_BASE_POWER_MAX {
				t.Errorf("base power %d outside %d..%d", *got.Override.BasePower, MOVE_BASE_POWER_MIN, MOVE_BASE_POWER_MAX)
			}

			if got.Override.Type != move.Type {
				t.Errorf("Invalid move type, got %s want %s", got.Override.Type, move.Type)
			}

			if got.Override.Category != move.Category {
				t.Errorf("Invalid move category, got %s want %s", got.Override.Category, move.Category)
			}

			if *got.Override.Critical != move.Critical {
				t.Errorf("Invalid value for move critical field got %v want %v", *got.Override.Critical, move.Critical)
			}

			if *got.MetronomeHits != 0 {
				t.Errorf("Invalid value for move metronome hits %v", *got.MetronomeHits)
			}
		})
	}
}
