package montecarlo

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

const size = 1_000_000
const times = 12

var tests = []struct {
	percentile float64
	need       float64
}{
	{2.5, 17},
	{97.5, 25},
}

func TestMonteCarlo(t *testing.T) {
	distr := newNormalDistribution(4, 6)

	results := run(times*size, rand.New(rand.NewPCG(1, 2)), func(r *rand.Rand) float64 {
		return 100 / distr.sample(r)
	})

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v percentile", test.percentile), func(t *testing.T) {
			idx := int(times * size * test.percentile / 100)
			if !closeEnough(results[idx], test.need) {
				t.Fatalf("%v percentile should be '%v', but we got '%v'", test.percentile, test.need, results[idx])
			}
		})
	}
}

func TestMonteCarloParallel(t *testing.T) {
	distr := newNormalDistribution(4, 6)

	seed := uint64(0)
	rngf := func() *rand.Rand {
		seed += 1
		return rand.New(rand.NewPCG(seed, seed+1))
	}

	results := parallelRun(times, size, rngf, func(r *rand.Rand) float64 {
		return 100 / distr.sample(r)
	})

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v percentile", test.percentile), func(t *testing.T) {
			idx := int(times * size * test.percentile / 100)
			if !closeEnough(results[idx], test.need) {
				t.Fatalf("%v percentile should be '%v', but we got '%v'", test.percentile, test.need, results[idx])
			}
		})
	}
}
