package distribution

import (
	"fmt"
	"math"
	"testing"
)

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

func TestNormalDistribution(t *testing.T) {
	tests := []struct {
		from  float64
		to    float64
		mu    float64
		sigma float64
	}{
		{5, 10, 7.5, 1.27551020408},
		{-2.5, 2.5, 0, 1.27551020408},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("from %v to %v", test.from, test.to), func(t *testing.T) {
			distr := NewNormal(test.from, test.to)

			if !floatEqual(distr.mu, test.mu) {
				t.Fatalf("mu should be '%v', but we got '%v'", test.mu, distr.mu)
			}

			if !floatEqual(distr.sigma, test.sigma) {
				t.Fatalf("sigma should be '%v', but we got '%v'", test.sigma, distr.sigma)
			}
		})
	}
}
