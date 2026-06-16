package models

type BattleMon struct {
	Name  string `json:"name"`
	Types []struct {
		Name string `json:"name"`
	} `json:"types"`
	Gender    string `json:"gender"`
	Level     int    `json:"level"`
	BaseStats struct {
		Hp  int `json:"hp"`
		Atk int `json:"atk"`
		Def int `json:"def"`
		Spa int `json:"spa"`
		Spd int `json:"spd"`
		Spe int `json:"spe"`
	} `json:"baseStats"`
	Ivs struct {
		Hp  int `json:"hp"`
		Atk int `json:"atk"`
		Def int `json:"def"`
		Spa int `json:"spa"`
		Spd int `json:"spd"`
		Spe int `json:"spe"`
	} `json:"ivs"`
	Evs struct {
		Hp  int `json:"hp"`
		Atk int `json:"atk"`
		Def int `json:"def"`
		Spa int `json:"spa"`
		Spd int `json:"spd"`
		Spe int `json:"spe"`
	} `json:"evs"`
	BattleMods struct {
		Hp  int `json:"hp"`
		Atk int `json:"atk"`
		Def int `json:"def"`
		Spa int `json:"spa"`
		Spd int `json:"spd"`
		Spe int `json:"spe"`
	} `json:"battleMods"`
	Nature    string      `json:"nature"`
	Ability   string      `json:"ability"`
	Item      interface{} `json:"item"`
	Status    string      `json:"status"`
	CurrentHp struct {
		Hp   int `json:"hp"`
		Perc int `json:"perc"`
	} `json:"currentHp"`
	Moves []struct {
		Name      string `json:"name"`
		BasePower int    `json:"basePower"`
		Type      struct {
			Name string `json:"name"`
		} `json:"type"`
		Category string `json:"category"`
		IsCrit   bool   `json:"isCrit"`
		Streak   int    `json:"streak"`
	} `json:"moves"`
}

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

/* Ability */
type Ability struct {
	Kind	string
	Id		string
	Name	string
}

/* Item */
type NaturalGiftItem struct {
	BasePower int    `json:"basePower"`
	Type      string `json:"type"`
}

type Item struct {
	Kind			string			`json:"kind"`
	Id				string			`json:"id"`
	Name			string			`json:"name"`
	IsBerry			bool			`json:"isBerry"`
	NaturalGift		NaturalGiftItem	`json:"naturalGift"`
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

/* Nature */
type Nature struct {
	Kind  string `json:"kind"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Plus  string `json:"plus"`
	Minus string `json:"minus"`
}

/* Type */
type Type struct {
	Kind  			string `json:"kind"`
	ID    			string `json:"id"`
	Name  			string `json:"name"`
	Effectiveness	map[string]float64 `json:"effectiveness"`
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

type Data struct {
	Abilities	[]Ability
	Items		[]Item
	Moves		[]Move
	Natures		[]Nature
	Types		[]Type
	Species		[]Specie
}

