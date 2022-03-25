package burrow

import (
	"fmt"
	"math/rand"
	"time"
)

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

func UniformTimestampDistribution (tStart time.Time, uniformRange time.Duration) (SampleDistribution[time.Time], error) {
	if uniformRange <= 0 {
		return nil, fmt.Errorf("Requires a non-zero duration")
	}

	distroFunc := func() time.Time {
		window := time.Duration(rand.Int63n(int64(uniformRange)))
		return tStart.Add(window)
	}

	return SampleDistribution[time.Time](distroFunc), nil
}
