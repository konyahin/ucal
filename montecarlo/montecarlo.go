package montecarlo

import (
	"math/rand/v2"
	"slices"
	"sync"
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
	results := make([]float64, size*times)
	var wg sync.WaitGroup

	for i := range times {
		rng := rngf()
		start := i * size
		wg.Go(func() {
			for j := range size {
				results[start+j] = f(rng)
			}
		})
	}

	wg.Wait()
	slices.Sort(results)
	return results
}
