package calc

type statKind string
type moveCategory string
type ivOption string
type evOption string
type statusCondition string
type gender string

const (
	// move category
	physical moveCategory = "Physical"
	special  moveCategory = "Special"
	status   moveCategory = "Status"

	// stat kind
	statHp  statKind = "hp"
	statAtk statKind = "atk"
	statDef statKind = "def"
	statSpa statKind = "spa"
	statSpd statKind = "spd"
	statSpe statKind = "spe"

	// iv option
	ivOptimal   ivOption = "optimal"
	ivWorstCase ivOption = "worstCase"
	ivCustom    ivOption = "custom"

	// ev option
	evOptimal   evOption = "optimal"
	evWorstCase evOption = "worstCase"
	evCustom    evOption = "custom"

	// status condition
	noneStatus    statusCondition = "none"
	healthy       statusCondition = "healthy"
	poisoned      statusCondition = "poisoned"
	badlyPoisoned statusCondition = "badly poisoned"
	burned        statusCondition = "burned"
	paralyzed     statusCondition = "paralyzed"
	asleep        statusCondition = "asleep"
	frozen        statusCondition = "frozen"

	// gender
	male    gender = "M"
	female  gender = "F"
	neutral gender = "N"

	// fixed damage
	fixedDmgDragonRage = 40
	fixedDmgSonicBoom  = 20

	// speed modifiers
	combineModsMin      = 410
	combineModsMax      = 131172
	fixedPointScaleBase = 4096
	halfPointScaleBase  = 2048

	speedMod1X   = fixedPointScaleBase
	speedMod2X   = 2 * fixedPointScaleBase
	speedMod1_5X = 6144
	speedMod0_5X = halfPointScaleBase

	maxSpeed = 10000

	resultDamageSliceSize = 16

	defaultHitCount = 1
)
