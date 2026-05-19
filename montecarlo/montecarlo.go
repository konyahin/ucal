package montecarlo

import (
	"math/rand/v2"

	"golang.org/x/sync/errgroup"

	"github.com/caio/go-tdigest/v5"
)

const (
	defaultSize = 400_000
	defaultStep = 50_000
)

type Simulation struct {
	size int
	f    func(*rand.Rand) float64
}

type Result struct {
	digest *tdigest.TDigest
}

func New(f func(*rand.Rand) float64) *Simulation {
	return &Simulation{
		size: defaultSize,
		f:    f,
	}
}

func (s *Simulation) Run() (*Result, error) {
	step := min(s.size, defaultStep)
	times := s.size / step

	digests := make([]*tdigest.TDigest, times)
	var wg errgroup.Group

	for i := range times {
		d, err := tdigest.New()
		if err != nil {
			return nil, err
		}

		digests[i] = d
		rng := rngFactory()
		wg.Go(func() error {
			for range step {
				err := d.Add(s.f(rng))
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	final, err := tdigest.New()
	if err != nil {
		return nil, err
	}

	for _, d := range digests {
		if err := final.Merge(d); err != nil {
			return nil, err
		}
	}

	return &Result{digest: final}, nil
}

// Percentile returns the value at the p-th percentile of the simulated
// distribution, where p is in [0, 100].
func (r *Result) Percentile(p float64) float64 {
	return r.digest.Quantile(p / 100)
}

func rngFactory() *rand.Rand {
	return rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
}
