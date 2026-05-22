package montecarlo

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"ucal/distribution"
)

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
	f := func(r *rand.Rand) (float64, error) {
		return 100 / distribution.Normal(4, 6, r), nil
	}

	mc := New(f)
	results, err := mc.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v percentile", test.percentile), func(t *testing.T) {
			result := results.Percentile(test.percentile)
			if !closeEnough(result, test.need) {
				t.Fatalf("%v percentile should be '%v', but we got '%v'", test.percentile, test.need, result)
			}
		})
	}
}

func BenchmarkMonteCarlo(b *testing.B) {
	f := func(r *rand.Rand) (float64, error) {
		return 100 / distribution.Normal(4, 6, r), nil
	}
	mc := New(f)
	for b.Loop() {
		_, _ = mc.Run(context.Background())
	}
}

func closeEnough(a, b float64) bool {
	return math.Abs(a-b) <= 0.1
}
