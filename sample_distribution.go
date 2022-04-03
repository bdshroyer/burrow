package burrow

import (
	"golang.org/x/exp/constraints"
	"time"
)

// A distribution type intended to cover both Ordered types and other types like Time that behave in an ordered fashion, but don't conform to the programmatic rules of Go type constraints.
type Rangeable interface {
	constraints.Ordered | time.Time
}

// A SampleDistribution is a function that returns a rangeable value drawn from some distribution.
type SampleDistribution[T Rangeable] func() T
