package burrow

import (
	"time"
	"fmt"
)


type TimeBox [2]time.Duration

type DeliveryNetworkConfig struct {
	HubNodes, StopNodes uint
	Distro SampleDistribution[time.Time]
	EdgeBounds *TimeBox
}

func (spec NetworkSpec) parseDistribution() (SampleDistribution[time.Time], error) {
	var distro SampleDistribution[time.Time]
	var err error

	if distroSpec := spec.GetUniform(); distroSpec != nil {
		tStart := spec.Start.AsTime()
		durationRange := spec.End.AsTime().Sub(spec.Start.AsTime())

		if distro, err = UniformTimestampDistribution(tStart, durationRange); err != nil {
			return nil, err
		}
	}

	if distroSpec := spec.GetGaussian(); distroSpec != nil {
		tMean := time.UnixMilli(int64(distroSpec.Mean))
		dStdDev := time.Duration(distroSpec.StdDev)
		if distro, err = GaussianTimestampDistribution(tMean, dStdDev); err != nil {
			return nil, err
		}
	}

	return distro, err
}

// Generates a NetworkConfig from a NetworkSpec. Returns an error if it's unable to convert the distribution
// specification into an actual distribution sampling function.
//
// Note that this conversion does no validation on the non-distribution attributes of the network spec. Bad values like negative node counts are expected to be handled by the generating method.
func NewNetworkConfig(spec NetworkSpec) (*DeliveryNetworkConfig, error) {

	distro, err := spec.parseDistribution()
	if err != nil {
		return nil, err
	}

	if distro == nil {
		return nil, fmt.Errorf("No distribution provided")
	}

	cfg := &DeliveryNetworkConfig{
		HubNodes: uint(spec.Hubs),
		StopNodes: uint(spec.Stops),
		EdgeBounds: &TimeBox{spec.ShortEdge.AsDuration(), spec.LongEdge.AsDuration()},
		Distro: distro,
	}

	return cfg, nil
}
