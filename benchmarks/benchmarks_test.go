package benchmarks_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow"
	"github.com/onsi/gomega/gmeasure"
)

func toKB(B uint64) float64 {
	return float64(B) / 1024
}

func toMB(B uint64) float64 {
	return float64(B) / (1024 * 1024)
}

func Today() time.Time {
	payload := time.Now()
	year, month, day := payload.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, payload.Location())
}

var _ = Describe("Benchmark", func() {
	benchmarkNetworkCreation := func(experiment *gmeasure.Experiment, cfg burrow.DeliveryNetworkConfig) func(int) {
		return func(idx int) {
			stopwatch := experiment.NewStopwatch()
			G, err := burrow.MakeDeliveryNetwork(cfg)
			stopwatch.Record("Creation Time", gmeasure.Precision(time.Microsecond))

			Expect(err).NotTo(HaveOccurred())

			stopwatch.Reset()
			nodes := G.Nodes()
			stopwatch.Record("Node Gathering Time", gmeasure.Precision(time.Microsecond))

			stopwatch.Reset()
			for nodes.Next() {
			}
			stopwatch.Record("Node Iterating Time", gmeasure.Precision(time.Microsecond))

			stopwatch.Reset()
			edges := G.Edges()
			stopwatch.Record("Edge Gathering Time", gmeasure.Precision(time.Microsecond))

			nEdges := float64(edges.Len())
			nStops := float64(len(G.Stops))
			nHubs := float64(len(G.Hubs))

			stopwatch.Reset()
			for edges.Next() {
			}
			stopwatch.Record("Edge Iterating Time", gmeasure.Precision(time.Microsecond))

			meanOutDegree := float64(0.0)
			for _, edges := range G.DEdges {
				meanOutDegree += float64(len(edges)) / nStops
			}

			experiment.RecordValue("nEdges", nEdges, gmeasure.Units("edges"))
			experiment.RecordValue("delivery network edge density", nEdges/(2*nHubs*nStops+(nStops*(nStops-1))/2))
			experiment.RecordValue("mean out-degree", meanOutDegree)
		}
	}

	Context("Create a delivery network with five hubs and 10000 stops", func() {
		It("Runs on a uniform distribution", func() {
			experiment := gmeasure.NewExperiment("Network Creation [uniform]")
			AddReportEntry(experiment.Name, experiment)

			distro, err := burrow.UniformTimestampDistribution(Today(), 24*time.Hour)
			Expect(err).NotTo(HaveOccurred())

			cfg := burrow.DeliveryNetworkConfig{HubNodes: 5, StopNodes: 10000, Distro: distro}

			experiment.Sample(
				benchmarkNetworkCreation(experiment, cfg),
				gmeasure.SamplingConfig{N: 10},
			)
		})

		It("Runs on a Gaussian distribution with bounds", func() {
			experiment := gmeasure.NewExperiment("Network Creation [Gaussian]")
			AddReportEntry(experiment.Name, experiment)

			distro, err := burrow.GaussianTimestampDistribution(
				Today().Add(11*time.Hour),
				2*time.Hour,
			)

			Expect(err).NotTo(HaveOccurred())

			cfg := burrow.DeliveryNetworkConfig{
				HubNodes:   5,
				StopNodes:  10000,
				Distro:     distro,
				EdgeBounds: &burrow.TimeBox{0 * time.Hour, 6 * time.Hour},
			}

			experiment.Sample(
				benchmarkNetworkCreation(experiment, cfg),
				gmeasure.SamplingConfig{N: 10},
			)
		})

	})
})
