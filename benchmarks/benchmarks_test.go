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
					_, err = burrow.MakeDeliveryNetwork(1, 10000, distro)
					stopwatch.Record("Creation Time", gmeasure.Precision(time.Microsecond))

					Expect(err).NotTo(HaveOccurred())
				},
				gmeasure.SamplingConfig{N: 10},
			)
		})
	})
})
