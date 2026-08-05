package calc

type StatKind string
type MoveCategory string
type IvOption string
type EvOption string
type CalcDataType string
type StatusCondition string
type Gender string

const (
	// move category
	Physical	MoveCategory	= "Physical"
	Special		MoveCategory	= "Special"
	Status		MoveCategory	= "Status"

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

	// gender
	Male			Gender				= "M"
	Female			Gender				= "F"
	Neutral			Gender				= "N"

	// fixed damage
	FIXED_DMG_DRAGON_RAGE	= 40
	FIXED_DMG_SONIC_BOOM	= 20

	// speed modifiers
	COMBINE_MODS_MIN		= 410
	COMBINE_MODS_MAX		= 131172
	FIXED_POINT_SCALE_BASE	= 4096
	HALF_POINT_SCALE_BASE	= 2048

	SPEED_MOD_1X			= FIXED_POINT_SCALE_BASE
	SPEED_MOD_2X			= 2 * FIXED_POINT_SCALE_BASE
	SPEED_MOD_1_5X			= 6144
	SPEED_MOD_0_5X			= HALF_POINT_SCALE_BASE

	MAX_SPEED				= 10000
)