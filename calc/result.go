package calc

import (
	"math"
)

func toPercent(value int, maxHP int) float64 {
	return math.Floor(float64(value*1000)/float64(maxHP)) / 10
}

func (c *Calculation) percentRange() (minPerc float64, maxPerc float64) {
	min := toPercent(c.Result.Range.Min, c.defender.effectiveHpStat())
	max := toPercent(c.Result.Range.Max, c.defender.effectiveHpStat())

	return min, max
}
