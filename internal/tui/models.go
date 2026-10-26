package tui

type Side int

var DEFAULT_HITS = hitOptions[0].label

const (
	LeftSide Side = iota
	RightSide

	DEFAULT_GENDER = "Female"
	DEFAULT_NATURE = "Hardy"

	// TODO: check if they can be combine with constants from models or calc
	DISPLAY_FEMALE_NAME     = "Female"
	DISPLAY_MALE_NAME       = "Male"
	DISPLAY_GENDERLESS_NAME = "Genderless"

	// TODO: check if they can be combine with constants from models or calc
	DISPLAY_STATUS_HEALTHY        = "Healthy"
	DISPLAY_STATUS_POISONED       = "Poisoned"
	DISPLAY_STATUS_BADLY_POISONED = "Badly Poisoned"
	DISPLAY_STATUS_BURNED         = "Burned"
	DISPLAY_STATUS_PARALYZED      = "Paralyzed"
	DISPLAY_STATUS_ASLEEP         = "Asleep"
	DISPLAY_STATUS_FROZEN         = "Frozen"
)

var STATUS_DISPLAY_NAMES = []string{
	DISPLAY_STATUS_HEALTHY,
	DISPLAY_STATUS_POISONED,
	DISPLAY_STATUS_BADLY_POISONED,
	DISPLAY_STATUS_BURNED,
	DISPLAY_STATUS_PARALYZED,
	DISPLAY_STATUS_ASLEEP,
	DISPLAY_STATUS_FROZEN,
}
