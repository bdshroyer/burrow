package burrow

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"time"
)

// A distribution type intended to cover both Ordered types and other types like Time that behave in an ordered fashion, but don't conform to the programmatic rules of Go type constraints.
type Rangeable interface {
	constraints.Ordered | time.Time
}

// A SampleDistribution is a function that returns a rangeable value drawn from some distribution.
type SampleDistribution[T Rangeable] func() T

// A SampleGenerator tracks a distribution function to be sampled from, as well as providing a means to clean up the generator if it's interrupted. This is important because the generator spins off a goroutine to draw samples from the distribution, and if the full range of samples isn't consumed by a receiver then the goroutine hangs aroun in a blocked state.
type SampleGenerator[T Rangeable] struct {
	Distro SampleDistribution[T]
	Quit   chan bool
}

func NewSampleGenerator[T Rangeable] (distro SampleDistribution[T]) (*SampleGenerator[T], error) {
	if distro == nil {
		return nil, fmt.Errorf("Must receive a non-null sample distribution.")
	}

	return &SampleGenerator[T]{Distro: distro, Quit: make(chan bool)}, nil
}

// Sample generates `n` samples drawn from the probability distribution `distro`, which it feeds into the channel returned to the caller. In effect, this is a Python-style generator for distribution samples.
// Note: Make sure to call the sampler's Stop() methods if the Sample function is not going to be allowed to finish (if, for example, the consumer decides to exit its consumption loop).
func (sampler *SampleGenerator[T]) Sample (n uint) chan T {
	output := make(chan T)

	go func() {
		defer close(output)

		for i := uint(0); i < n; i++ {
			select {
			case <-sampler.Quit:
				return

			default:
				output <- sampler.Distro()
			}
		}
	}()

	return output
}

// If the consumer decides to prematurely cease collecting samples from the generator, Stop() tells the sample generator to stop producing more samples and exit.
func (sampler *SampleGenerator[T]) Stop() {
	sampler.Quit <- true
}
