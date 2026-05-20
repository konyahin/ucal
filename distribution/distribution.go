package distribution

import (
	"math/rand/v2"
)

// this is the 97.5% quantile of the standard normal distribution
// used to construct a symmetric 95% confidence interval
const z975 = 1.96

func Normal(from, to float64, r *rand.Rand) float64 {
	return (from+to)/2 + (to-from)/(z975*2)*r.NormFloat64()
}
