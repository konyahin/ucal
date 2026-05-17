package montecarlo

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

const size = 500_000
const times = 2

const benchmarkSize = 1_000_000
const benchmarkTimes = 10

var tests = []struct {
	percentile float64
	need       float64
}{
	{2.5, 16.7},
	{10, 17.7},
	{20, 18.4},
	{30, 19},
	{40, 19.5},
	{50, 20},
	{60, 20.5},
	{70, 21.1},
	{80, 21.8},
	{90, 23},
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
	seed := uint64(0)
	rngf := func() *rand.Rand {
		seed += 1
		return rand.New(rand.NewPCG(seed, seed+1))
	}

	distr := newNormalDistribution(4, 6)
	digest, err := parallelRun(times, size, rngf, func(r *rand.Rand) float64 {
		return 100 / distr.sample(r)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v percentile", test.percentile), func(t *testing.T) {
			result := digest.Quantile(test.percentile / 100)
			if !closeEnough(result, test.need) {
				t.Fatalf("%v percentile should be '%v', but we got '%v'", test.percentile, test.need, result)
			}
		})
	}
}

func BenchmarkSequential(b *testing.B) {
	distr := newNormalDistribution(4, 6)
	rng := rand.New(rand.NewPCG(1, 2))
	f := func(r *rand.Rand) float64 { return 100 / distr.sample(r) }
	for b.Loop() {
		run(benchmarkTimes*benchmarkSize, rng, f)
	}
}

func BenchmarkParallel(b *testing.B) {
	distr := newNormalDistribution(4, 6)
	seed := uint64(0)
	rngf := func() *rand.Rand { seed++; return rand.New(rand.NewPCG(seed, 0)) }
	f := func(r *rand.Rand) float64 { return 100 / distr.sample(r) }
	for b.Loop() {
		parallelRun(benchmarkTimes, benchmarkSize, rngf, f)
	}
}
