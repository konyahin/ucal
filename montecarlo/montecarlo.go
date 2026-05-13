package montecarlo

import "slices"

// montecarlo runs `f` `size` times and returns the results sorted ascending
func run(size int, f func() float64) []float64 {
	results := make([]float64, size)
	for i := range size {
		results[i] = f()
	}

	slices.Sort(results)
	return results
}
