package montecarlo

import "math"

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

func closeEnough(a, b float64) bool {
	return math.Abs(a-b) <= 0.1
}
