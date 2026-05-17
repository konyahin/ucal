package montecarlo

import (
	"math/rand/v2"
	"slices"
	"sync"

	"github.com/caio/go-tdigest/v5"
)

func run(size int, rng *rand.Rand, f func(rng *rand.Rand) float64) []float64 {
	results := make([]float64, size)
	for i := range size {
		results[i] = f(rng)
	}

	slices.Sort(results)
	return results
}

func parallelRun(times, size int, rngf func() *rand.Rand, f func(rng *rand.Rand) float64) (*tdigest.TDigest, error) {
	digests := make([]*tdigest.TDigest, times)
	var wg sync.WaitGroup

	for i := range times {
		d, err := tdigest.New()
		if err != nil {
			return nil, err
		}

		digests[i] = d
		rng := rngf()
		wg.Go(func() {
			for range size {
				err := d.Add(f(rng))
				if err != nil {
					panic(err)
				}
			}
		})
	}

	wg.Wait()

	final, err := tdigest.New()
	if err != nil {
		return nil, err
	}

	for _, d := range digests {
		if err := final.Merge(d); err != nil {
			return nil, err
		}
	}
	return final, nil
}
