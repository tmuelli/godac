package models

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
	Possibilities []int
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

/* Ability */
type Ability struct {
	Kind	string
	Id		string
	Name	string
}

func (a Ability) Key() string {
	return a.Name
}

/* Item */
type ItemEffect string
type ItemConditionFlag uint16

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
type MoveSelfBoosts struct {
	Atk int `json:"atk"`
	Def int `json:"def"`
	Spa int `json:"spa"`
	Spd int `json:"spd"`
	Spe int `json:"spe"`
}

// TODO: maybe use stat type battle modifier later
type MoveSelf struct {
	Boosts MoveSelfBoosts `json:"boosts"`
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
	Priority       	int    `json:"priority,omitempty"`
	Recoil         	[]int  `json:"recoil,omitempty"`
	HasCrashDamage 	bool   `json:"hasCrashDamage,omitempty"`
	StruggleRecoil 	bool   `json:"struggleRecoil,omitempty"`
	Multiaccuracy  	bool   `json:"multiaccuracy,omitempty"`
	Self           	MoveSelf `json:"self,omitempty"`
	BreaksProtect 	bool `json:"breaksProtect,omitempty"`
}

func (m Move) Key() string {
	return m.Name
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


/* Calculation context */
type CalcContext struct {
	SuperEffective		bool
}
