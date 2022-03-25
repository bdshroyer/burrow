package burrow

import (
	"fmt"
	"math/rand"
	"golang.org/x/exp/constraints"
	"time"
)

// A distribution type intended to cover both Ordered types and other types like Time that behave in an ordered fashion, but don't conform to the programmatic rules of Go type constraints.
type Rangeable interface {
	constraints.Ordered | time.Time
}

// A SampleDistribution is a function that returns a rangeable value drawn from some distribution.
type SampleDistribution[T Rangeable] func() T

// MakeUniformDistribution produces a SampleDistribution function with a uniform probability over the given range.
func MakeUniformDistribution (uniformRange float64) (SampleDistribution[float64], error) {
	if uniformRange <= 0 {
		return nil, fmt.Errorf("Distribution range must be a positive number.")
	}

	distroFunc := func() float64 {
		return rand.Float64() * uniformRange
	}

	return SampleDistribution[float64](distroFunc), nil
}

// A DistroGenerator tracks a distribution function to be sampled from, as well as providing a means to clean up the generator if it's interrupted. This is important because the generator spins off a goroutine to draw samples from the distribution, and if the full range of samples isn't consumed by a receiver then the goroutine hangs aroun in a blocked state.
type DistroGenerator[T Rangeable] struct {
	Distro SampleDistribution[T]
	Quit   chan bool
}

func NewDistroGenerator[T Rangeable] (distro SampleDistribution[T]) (*DistroGenerator[T], error) {
	if distro == nil {
		return nil, fmt.Errorf("Must receive a non-null sample distribution.")
	}

	return &DistroGenerator[T]{Distro: distro, Quit: make(chan bool)}, nil
}

// Sample generates `n` samples drawn from the probability distribution `distro`, which it feeds into the channel returned to the caller. In effect, this is a Python-style generator for distribution samples.
// Note: Make sure to call the distribution's Stop() methods if the Sample function is not going to be allowed to finish (if, for example, the consumer decides to exit its consumption loop).
func (distro *DistroGenerator[T]) Sample (n uint) chan T {
	output := make(chan T)

	go func() {
		defer close(output)

		for i := uint(0); i < n; i++ {
			select {
			case <-distro.Quit:
				return

			default:
				output <- distro.Distro()
			}
		}
	}()

	return output
}

// If the consumer decides to prematurely cease collecting samples from the generator, Stop() tells the sample generator to stop producing more samples and exit.
func (distro *DistroGenerator[T]) Stop() {
	distro.Quit <- true
}
