package calc

import "github.com/tmuelli/godac/internal/models"

func (g *Gender) isValid(input *string) bool {
	if input == nil {
		// no gender was specified => this is not an error => but gender logic can not be active than
		return true
	}

	*g = Gender(*input)

	switch Gender(*input) {
	case Male, Female, Neutral:
		return true
	default:
		return false
	}
}

func (g *Gender) isValidForSpecie(s *models.Specie) bool {
	if s.Gender == "" && *g != Neutral {
		// if no specific gender is set => means female and male are valid but not neutral
		return true
	}

	// if gender on specie is set => input gender has to match that
	return s.Gender == string(*g)
}