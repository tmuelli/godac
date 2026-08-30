package models

import (
	"math"
	"slices"
)

/****************************/
/* Types: user input/output */
/****************************/

/* User input */
type UserInput struct {
	Name          string   `json:"name"`
	Level         int      `json:"level"`
	Nature        string   `json:"nature"`
	AttackMove    string   `json:"attackMove"`
	Moves         []string `json:"moves"`
	HealthPercent int      `json:"healthPercent"`
	IvSetting     string   `json:"ivSetting"`
	EvSetting     string   `json:"evSetting"`
	Item          *string  `json:"item"`
	Status        *string  `json:"status"`
	Gender        *string  `json:"gender"`
	Hits          *int     `json:"hits"`
}

/* User field input */
type UserFieldSideInput struct {
	StealthRocks	bool	`json:"stealth_rocks"`
	Spikes			int		`json:"spikes"`
	Reflect			bool	`json:"reflect"`
	LightScreen		bool	`json:"light_screen"`
	Protect			bool	`json:"protect"`
	LeechSeed		bool	`json:"leech_seed"`
	Foresight		bool	`json:"foresight"`
	HelpingHand		bool	`json:"helping_hand"`
	Tailwind		bool	`json:"tailwind"`
	FlowerGift		bool	`json:"flower_gift"`
	PowerTrick		bool	`json:"power_trick"`
	SwitchingOut	bool	`json:"switching_out"`
}

type UserFieldInput struct {
	Type			string				`json:"type"`
	Level			int					`json:"level"`
	Terrain			string				`json:"terrain"`
	IsGravity		bool				`json:"is_gravity"`
	Left			UserFieldSideInput	`json:"left"`
	Right			UserFieldSideInput	`json:"right"`
}

/* Calculation result */
type BattleResultRange struct {
	Min		int
	MinPerc	float64
	Max		int
	MaxPerc	float64
}

type BattleResult struct {
	Range BattleResultRange
	DmgMatrix [][]int
}

/********************/
/* Types: core data */
/********************/

/* BaseStats */
type BaseStats struct {
	Hp	int	`json:"hp"`
	Atk	int	`json:"atk"`
	Def	int	`json:"def"`
	Spa	int	`json:"spa"`
	Spd	int	`json:"spd"`
	Spe	int	`json:"spe"`
}

/* IvStats */
type IvStats struct {
	Hp	int
	Atk	int
	Def	int
	Spa	int
	Spd	int
	Spe	int
}

/* EvStats */
type EvStats struct {
	Hp	int
	Atk	int
	Def	int
	Spa	int
	Spd	int
	Spe	int
}

/* Boosts */
const (
	MAX_STAT_BOOST = 6
	MIN_STAT_BOOST = -6
)

type StatBoost struct {
	value	int
}

func (sb *StatBoost) Set(mod int) {
	if mod == 0 {
		return
	}

	newValue := sb.value + mod
	if newValue < MIN_STAT_BOOST {
		sb.value = MIN_STAT_BOOST
	} else if newValue > MAX_STAT_BOOST {
		sb.value = MAX_STAT_BOOST
	} else {
		sb.value = newValue
	}
}

type Boosts struct {
	Atk StatBoost
	Def StatBoost
	Spa StatBoost
	Spd StatBoost
	Spe StatBoost
}

func (b *Boosts) Get(stat string) (int) {
	switch (stat) {
	case "atk":
		return b.Atk.value
	case "def":
		return b.Def.value
	case "spa":
		return b.Spa.value
	case "spd":
		return b.Spd.value
	case "spe":
		return b.Spe.value
	default:
		return 0
	}
}

/* Ability */
type Ability struct {
	Kind	string
	Id		string
	Name	string
	Active  bool
}

func (a Ability) Key() string {
	return a.Name
}

func (a *Ability) IsActive(atkAbility string) bool {
	ignoredAbilities := []string{"Battle Armor", "Clear Body", "Damp", "Dry Skin",
    "Filter", "Flash Fire", "Flower Gift", "Heatproof",
    "Hyper Cutter", "Immunity", "Inner Focus", "Insomnia",
    "Keen Eye", "Leaf Guard", "Levitate", "Lightning Rod",
    "Limber", "Magma Armor", "Marvel Scale", "Motor Drive",
    "Oblivious", "Own Tempo", "Sand Veil", "Shell Armor",
    "Shield Dust", "Simple", "Snow Cloak", "Solid Rock",
    "Soundproof", "Sticky Hold", "Storm Drain", "Sturdy",
    "Suction Cups", "Tangled Feet", "Thick Fat", "Unaware",
    "Vital Spirit", "Volt Absorb", "Water Absorb", "Water Veil",
    "White Smoke", "Wonder Guard"}
	
	if atkAbility == "moldbreaker" && !slices.Contains(ignoredAbilities, a.Id) {
		return false
	}

	return true
}

func (a *Ability) IsPreventingDamage(move *Move, defItem *Item, typeEffectivness float64, isGravity bool) bool {
	if ((move.HasType("???") || move.ID == "firefang") && a.Id == "wonderguard" && typeEffectivness <= 1.0) ||
	   (move.HasType("Fire") && a.Id == "flashfire") ||
	   (move.HasType("Water") && (a.Id == "dryskin" || a.Id == "waterabsorb")) ||
	   (move.HasType("Electric") && (a.Id == "motordrive" || a.Id == "voltabsorb")) ||
	   (move.HasType("Ground") && !isGravity && defItem != nil && defItem.Id != "ironball" && a.Id == "levitate") ||
	   (move.Flags.Sound != 0 && a.Id == "soundproof") {
		return true
	}
	
	return false
}

/* Item */
type ItemEffect string
type ItemConditionFlag uint16
type ItemDamageSubStage int

const (
	StatEffect		ItemEffect	= "stat"
	DamageEffect	ItemEffect	= "damage"
	MoveEffect		ItemEffect	= "move"
)

const (
	ItemConSuperEffective		ItemConditionFlag = 1 << iota
	ItemConNotTransformed
	ItemConPhysical
	ItemConSpecial
	ItemConMetronome
	ItemConLifeOrb
	ItemConTypeReduction
)

const (
	DamagePreLoop		ItemDamageSubStage = iota	// life orb, metronome => applied before final loop in calcModificationDamage
	DamageExpertBelt								// expert belt applied before tinted lens modifier (inside loop)
	DamageBerry										// berry mod applied after tinted lens modifier (inside loop)
)

type NaturalGiftItem struct {
	BasePower int    `json:"basePower"`
	Type      string `json:"type"`
}

type ItemRestriction struct {
	Species		[]string	`json:"species"`
}

type ItemCondition struct {
	SuperEffective	bool
	NotTransformed	bool
	Physical		bool
	Special			bool
	Metronome		bool
	LifeOrb			bool
	TypeReduction	bool
}

type ItemModifierOptions struct {
	MoveType	string
	TypeMod		float64
	SubStage	ItemDamageSubStage
}

type Item struct {
	Kind			string				`json:"kind"`
	Id				string				`json:"id"`
	Name			string				`json:"name"`
	IsBerry			bool				`json:"isBerry"`
	NaturalGift		NaturalGiftItem		`json:"naturalGift"`
	Effect			ItemEffect			`json:"effect"`
	Stats			[]string			`json:"stats"`
	Types			[]string			`json:"types"`
	Multiplier		float64				`json:"multiplier"`
	Restrictions	*ItemRestriction	`json:"restrictions"`
	Condition		*ItemCondition		`json:"condition"`
	FlingPower		int					`json:"flingPower"`
}

func (i Item) Key() string {
	return i.Name
}

func (c *ItemCondition) ToFlag() ItemConditionFlag {
	var flag ItemConditionFlag

	if c.SuperEffective { flag |= ItemConSuperEffective }
	if c.NotTransformed { flag |= ItemConNotTransformed }
	if c.Physical { flag |= ItemConPhysical }
	if c.Special { flag |= ItemConSpecial }
	if c.Metronome { flag |= ItemConMetronome }
	if c.LifeOrb { flag |= ItemConLifeOrb }
	if c.TypeReduction { flag |= ItemConTypeReduction }

	return flag
}

// TODO: maybe move this logic into json structure => add type json member like other type boosting items
func (i *Item) GetPlateBoostType() string {
	switch i.Id {
	case "fistplate":
		return "Fighting"
	case "skyplate":
		return "Flying"
	case "toxicplate":
		return "Posion"
	case "earthplate":
		return "Ground"
	case "stoneplate":
		return "Rock"
	case "insectplate":
		return "Bug"
	case "spookyplate":
		return "Ghost"
	case "ironplate":
		return "Steel"
	case "flameplate":
		return "Fire"
	case "splashplate":
		return "Water"
	case "meadowplate":
		return "Grass"
	case "zapplate":
		return "Electric"
	case "mindplate":
		return "Psychic"
	case "icicleplate":
		return "Ice"
	case "dracoplate":
		return "Dragon"
	case "dreadplate":
		return "Dark"
	default:
		return "Normal"
	}
}

/* Move */
type MoveFlags struct {
	Contact		int `json:"contact"`
	Bite		int `json:"bite"`
	Sound		int `json:"sound"`
	// TODO: Heal		int `json:"heal"`
	Punch		int `json:"punch"`
	Bullet		int `json:"bullet"`
	Pulse		int `json:"pulse"`
	Slicing		int `json:"slicing"`
	Wind		int `json:"wind"`
}

// TODO: maybe use stat type battle modifier later
type MoveSelf struct {
	Atk int `json:"atk"`
	Def int `json:"def"`
	Spa int `json:"spa"`
	Spd int `json:"spd"`
	Spe int `json:"spe"`
}

type Move struct {
	Kind  			string `json:"kind"`
	ID    			string `json:"id"`
	Name  			string `json:"name"`
	Flags 			MoveFlags `json:"flags"`
	Category       	string `json:"category"`
	Type           	string `json:"type"`
	BasePower      	int    `json:"basePower"`
	Drain          	[]int  `json:"drain,omitempty"`
	Target         	string `json:"target,omitempty"`
	Multihit       	[]int  `json:"multihit,omitempty"`
	Hits            int    `json:"hits,omitempty"`
	Priority       	int    `json:"priority,omitempty"`
	Recoil         	[]int  `json:"recoil,omitempty"`
	HasCrashDamage 	bool   `json:"hasCrashDamage,omitempty"`
	StruggleRecoil 	bool   `json:"struggleRecoil,omitempty"`
	Multiaccuracy  	bool   `json:"multiaccuracy,omitempty"`
	Self           	MoveSelf `json:"self,omitempty"`
	BreaksProtect 	bool `json:"breaksProtect,omitempty"`
	Critical		bool `json:"critical"`
}

func (m Move) Key() string {
	return m.Name
}

func (m *Move) HandlePainSplit(atkHP, defHP int) int {
	avg := math.Floor(float64(atkHP + defHP) / 2)
	dmg := math.Max(0, float64(defHP) - avg)
	return int(dmg)
}

func (m *Move) NeedSpecialMoveSetup() bool {
	specialMoves := []string{"weatherball", "judgement", "naturalgift", "brickbreak"}
	return slices.Contains(specialMoves, m.ID)
}

func (m *Move) HasType(ids ...string) bool {
	return slices.Contains(ids, m.Type)
}

func (m *Move) DropsStats() bool {
	boosts := []int{m.Self.Atk, m.Self.Def, m.Self.Spa, m.Self.Spd, m.Self.Spa}

	for _, s := range boosts {
		if s < 0 {
			return true
		}
	}

	return false
}

func (m *Move) GetBoostDrop(stat string) int {
	switch (stat) {
	case "atk":
		return m.Self.Atk
	case "def":
		return m.Self.Def
	case "spa":
		return m.Self.Spa
	case "spd":
		return m.Self.Spd
	case "spe":
		return m.Self.Spe
	default:
		return 0
	}
}

type MultihitItemUsage struct {
	HasAtkUsedItem bool
	HasDefUsedItem bool
}

func (m *Move) ValidMultihit(hits int) bool {
	maxHits := slices.Max(m.Multihit)
	minHits := slices.Min(m.Multihit)

	return hits >= minHits && hits <= maxHits
}

/* Nature */
type Nature struct {
	Kind  string `json:"kind"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Plus  string `json:"plus"`
	Minus string `json:"minus"`
}

func (n Nature) Key() string {
	return n.Name
}

/* Type */
type Type struct {
	Kind  			string `json:"kind"`
	ID    			string `json:"id"`
	Name  			string `json:"name"`
	Effectiveness	map[string]float64 `json:"effectiveness"`
}

func (t Type) Key() string {
	return t.Name
}

/* Specie */
type Specie struct {
	Kind			string				`json:"kind"`
	ID				string				`json:"id"`
	Name			string				`json:"name"`
	BaseStats		BaseStats			`json:"baseStats"`
	Types			[]string			`json:"types"`
	Weightkg		float64				`json:"weightkg"`
	Nfe				bool				`json:"nfe"`
	Gender			string				`json:"gender"`
	OtherFormes		[]string			`json:"otherFormes"`
	BaseSpecies		[]string			`json:"baseSpecies"`
	Abilities		map[string]string	`json:"abilities"`
}

func (s Specie) Key() string {
	return s.Name
}

/* Data */
type Data struct {
	Abilities	map[string]Ability
	Items		map[string]Item
	Moves		map[string]Move
	Natures		map[string]Nature
	Types		map[string]Type
	Species		map[string]Specie
}

/* Field */
type FieldType string
func (ft *FieldType) IsValid() bool {
	switch *ft {
	case Singles, Doubles:
		return true
	default:
		return false
	}
}

type FieldTerrain string
func (ft *FieldTerrain) IsValid() bool {
	switch *ft {
	case Clear, Sun, Rain, Sand, Hail:
		return true
	default:
		return false
	}
}

type FieldSpikesStack int
func (fss *FieldSpikesStack) IsValid() bool {
	switch *fss {
	case Spikes0, Spikes1, Spikes2, Spikes3:
		return true
	default:
		return false
	}
}

type FieldLevel int
func (fl *FieldLevel) IsValid() bool {
	switch *fl {
	case Level5, Level50, Level100:
		return true
	default:
		return false
	}
}

const (
	Singles		FieldType	= "Singles"
	Doubles		FieldType	= "Doubles"

	Clear		FieldTerrain	= "Clear"
	Sun			FieldTerrain	= "Sun"
	Rain		FieldTerrain	= "Rain"
	Sand		FieldTerrain	= "Sand"
	Hail		FieldTerrain	= "Hail"
	// Gravity		FieldTerrain	= "Gravity"

	Spikes0		FieldSpikesStack	= 0
	Spikes1		FieldSpikesStack	= 1
	Spikes2		FieldSpikesStack	= 2
	Spikes3		FieldSpikesStack	= 3

	Level5		FieldLevel			= 5
	Level50		FieldLevel			= 50
	Level100	FieldLevel			= 100
)

type FieldSide struct {
	Spikes			FieldSpikesStack
	StealthRocks	bool
	Reflect			bool
	LightScreen		bool
	Protect			bool
	LeechSeed		bool
	Foresight		bool
	HelpingHand		bool
	Tailwind		bool
	FlowerGift		bool
	PowerTrick		bool
	SwitchingOut	bool
}

type Field struct {
	Type		FieldType
	Level		FieldLevel
	Terrain		FieldTerrain
	IsGravity	bool
	Left		FieldSide
	Right		FieldSide
}

/* Calculation context */
type CalcContext struct {
	SuperEffective		bool
}
