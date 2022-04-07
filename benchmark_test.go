package burrow_test

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

var _ = Describe("Benchmark", func() {
	Context("Create a 1000-stop delivery network with one hub", func() {
		It("Runs on one core", func() {
			experiment := gmeasure.NewExperiment("Network Creation")
			AddReportEntry(experiment.Name, experiment)

			experiment.Sample(
				func(idx int) {
					distro, err := burrow.UniformTimestampDistribution(Today(), 24*time.Hour)
					Expect(err).NotTo(HaveOccurred())

					stopwatch := experiment.NewStopwatch()
					_, err = burrow.MakeDeliveryNetwork(1, 1000, distro)
					stopwatch.Record("Creation Time", gmeasure.Precision(time.Microsecond))

					Expect(err).NotTo(HaveOccurred())
				},
				gmeasure.SamplingConfig{N: 10},
			)
		})
	})
})
