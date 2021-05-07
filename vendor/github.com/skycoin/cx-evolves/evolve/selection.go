package evolve

import (
	"math/rand"
)

// Codes associated to each of the selection functions.
// const (
// 	tournamentRoulette = iota // Default
// 	selectionRoulette
// )

func tournamentSelection(errors []float64, chance float32, isMinimizing bool) (int, int) {
	idx := 0
	secondIdx := 0
	lowest := errors[0]
	for i, err := range errors {
		if isMinimizing {
			if err <= lowest && rand.Float32() <= chance {
				lowest = err
				secondIdx = idx
				idx = i
			}
		} else {
			if err >= lowest && rand.Float32() <= chance {
				lowest = err
				secondIdx = idx
				idx = i
			}
		}

	}
	return idx, secondIdx
}

// func rouletteSelection(errors []float64, isMinimizing bool) int {
// 	// var prevProb float64
// 	var errSum float64
// 	for _, err := range errors {
// 		errSum += err
// 	}
// 	value := rand.Float64() * errSum
// 	for i, err := range errors {
// 		value -= err
// 		if value <= 0 {
// 			return i
// 		}
// 	}

// 	// _ = prevProb

// 	return 0
// }
