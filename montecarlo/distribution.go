package montecarlo

import (
	"math/rand/v2"
)

// this is the 97.5% quantile of the standard normal distribution
// used to construct a symmetric 95% confidence interval
const z975 = 1.96

type normalDistribution struct {
	mu    float64
	sigma float64
}

// 95% confidence interval assumed for [from, to]
func newNormalDistribution(from, to float64) normalDistribution {
	return normalDistribution{
		mu:    (from + to) / 2,
		sigma: (to - from) / (z975 * 2),
	}
}

func (n normalDistribution) sample(r *rand.Rand) float64 {
	return n.mu + n.sigma*r.NormFloat64()
}
