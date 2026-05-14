package montecarlo

import (
	"math/rand/v2"
	"slices"
)

func run(size int, rng *rand.Rand, f func(rng *rand.Rand) float64) []float64 {
	results := make([]float64, size)
	for i := range size {
		results[i] = f(rng)
	}

	slices.Sort(results)
	return results
}

func parallelRun(times, size int, rngf func() *rand.Rand, f func(rng *rand.Rand) float64) []float64 {
	allResults := make([]float64, 0, size*times)
	resChan := make(chan []float64)

	for range times {
		rng := rngf()
		go func() {
			results := make([]float64, size)
			for i := range size {
				results[i] = f(rng)
			}
			resChan <- results
		}()
	}

	for range times {
		allResults = append(allResults, <-resChan...)
	}
	slices.Sort(allResults)
	return allResults
}
