package calc

type StatKind string
type MoveCategory string
type IvOption string
type EvOption string
type CalcDataType string
type StatusCondition string

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

	// calc data type
	Attacker	CalcDataType	= "attacker"
	Defender	CalcDataType	= "defender"

	// status condition
	NoneStatus		StatusCondition		= "none"
	Healthy			StatusCondition		= "healthy"
	Poisoned		StatusCondition		= "poisoned"
	BadlyPoisoned	StatusCondition		= "badly poisoned"
	Burned			StatusCondition		= "burned"
	Paralyzed		StatusCondition		= "paralyzed"
	Asleep			StatusCondition		= "asleep"
	Frozen			StatusCondition		= "frozen"
)