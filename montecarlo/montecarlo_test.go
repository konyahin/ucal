package montecarlo

import (
	"math/rand/v2"
	"testing"
)

func TestMonteCarlo(t *testing.T) {
	distr := newNormalDistribution(4, 6)
	
	r := rand.New(rand.NewPCG(1, 2))

	results := run(1_000_000, func() float64 {
		return 100 / distr.sample(r)
	})

	if !floatEqual(results[25_000], 16.667473385) {
		t.Fatalf("2.5 percentile should be '%v', but we got '%v'", 16.667473385542948, results[25_000])
	}
	if !floatEqual(results[975_000], 25.005009992) {
		t.Fatalf("97.5 percentile should be '%v', but we got '%v'", 25.005009992533125, results[975_000])
	}
}
