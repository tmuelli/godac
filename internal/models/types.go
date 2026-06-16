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
