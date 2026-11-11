package calc

import "github.com/tmuelli/godac/models"

func (g *gender) isValid(input *string) bool {
	if input == nil {
		// no gender was specified => this is not an error => but gender logic can not be active than
		return true
	}

	*g = gender(*input)

	switch gender(*input) {
	case male, female, neutral:
		return true
	default:
		return false
	}
}

func (g *gender) isValidForSpecie(s *models.Specie) bool {
	if s.Gender == "" && *g != neutral {
		// if no specific gender is set => means female and male are valid but not neutral
		return true
	}

	// if gender on specie is set => input gender has to match that
	return s.Gender == string(*g)
}
