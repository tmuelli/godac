package calc

type StatKind string
type MoveCategory string
type IvOption string
type EvOption string

const (
	// move category
	Physical	MoveCategory	= "Physical"
	Special		MoveCategory	= "Special"

	// stat kind
	StatHp		StatKind		= "hp"
	StatAtk		StatKind		= "atk"
	StatDef		StatKind		= "def"
	StatSpa		StatKind		= "spa"
	StatSpd		StatKind		= "spd"
	StatSpe		StatKind		= "spe"

	// iv option
	IvOptimal	IvOption		= "optimal"
	IvWorstCase	IvOption		= "worstCase"
	IvCustom	IvOption		= "custom"

	// ev option
	EvOptimal	EvOption		= "optimal"
	EvWorstCase	EvOption		= "worstCase"
	EvCustom	EvOption		= "custom"
)