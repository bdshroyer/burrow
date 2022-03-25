package burrow

import (
	"fmt"
	"math/rand"
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
