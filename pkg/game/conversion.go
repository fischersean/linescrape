package game

import (
	"math"
)

// Convert a decimal style odd to an american style
func DecimalToAmerican(dOdds float64) int64 {

	nearestFive := func(x float64) float64 {
		return math.Round(x/5.0) * 5.0
	}
	if dOdds >= 2.0 {
		return int64(nearestFive((dOdds - 1) * 100))
	}

	return -1 * int64(nearestFive(math.Abs((-100)/(dOdds-1))))
}
