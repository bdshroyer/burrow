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
	Context("Create a 10000-stop delivery network with one hub", func() {
		It("Runs on one core", func() {
			experiment := gmeasure.NewExperiment("Network Creation")
			AddReportEntry(experiment.Name, experiment)

			experiment.Sample(
				func(idx int) {
					distro, err := burrow.UniformTimestampDistribution(Today(), 24*time.Hour)
					Expect(err).NotTo(HaveOccurred())

					stopwatch := experiment.NewStopwatch()
					G, err := burrow.MakeDeliveryNetwork(1, 10000, distro)
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
						edges.WeightedEdge()
					}
					stopwatch.Record("Edge Iterating Time", gmeasure.Precision(time.Microsecond))

					experiment.RecordValue("nEdges", nEdges, gmeasure.Units("edges"))
					experiment.RecordValue("edge density", nEdges/(2*nHubs*nStops+(nStops*(nStops-1))/2))
				},
				gmeasure.SamplingConfig{N: 10},
			)
		})
	})
})
