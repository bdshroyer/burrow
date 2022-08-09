package burrow_test

import (
	"math"
	"math/rand"
	"time"

	"github.com/bdshroyer/burrow"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func today() time.Time {
	base := time.Now()
	return time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, base.Location())
}

func rawStdDev(xs []float64) float64 {
	var mu float64 = 0.0
	var variance float64 = 0.0

	N := float64(len(xs))

	for _, x := range xs {
		mu += x
	}

	mu /= N

	for _, x := range xs {
		variance += math.Pow(x-mu, 2)
	}

	variance /= N - 1
	return math.Sqrt(variance)
}

func oneSampleTtest(samples []float64, popMu float64) float64 {
	nSamples := len(samples)

	sampleMu, sampleStdDev := stat.MeanStdDev(samples, nil)
	sampleStdErr := sampleStdDev / math.Sqrt(float64(nSamples))

	tDistro := distuv.StudentsT{
		Mu:    popMu,
		Sigma: sampleStdErr,
		Nu:    float64(nSamples - 1),
	}

	pValue := 2 * tDistro.CDF(sampleMu)

	return pValue
}

var _ = Describe("Distros", func() {
	BeforeEach(func() {
		rand.Seed(3)
	})

	Context("MakeUniformDistribution", func() {
		When("Given a positive range", func() {
			It("Returns a sample distribution", func() {
				distro, err := burrow.MakeUniformDistribution(5.0)
				Expect(err).NotTo(HaveOccurred())

				Expect(distro()).To(BeNumerically("~", 5*0.71998, 1e-4))
				Expect(distro()).To(BeNumerically("~", 5*0.65263, 1e-4))
			})
		})

		When("Given a non-positive range", func() {
			It("Returns an error message", func() {
				distro, err := burrow.MakeUniformDistribution(0.0)
				Expect(distro).To(BeNil())
				Expect(err).To(MatchError("Distribution range must be a positive number."))
			})
		})
	})

	Context("UniformTimestampDistribution", func() {
		When("Given a valid time t0 and a positive range T", func() {
			It("Returns a sample distribution of timestamps from t0 to t0 + T", func() {
				t0 := today()
				window := 24 * time.Hour

				distro, err := burrow.UniformTimestampDistribution(t0, window)
				Expect(err).NotTo(HaveOccurred())

				t1 := distro()
				t2 := distro()

				// t1 is within the window range
				Expect(t1).To(BeTemporally(">=", today()))
				Expect(t1).To(BeTemporally("<", today().Add(window)))

				// t2 is within the window range
				Expect(t2).To(BeTemporally(">=", today()))
				Expect(t2).To(BeTemporally("<", today().Add(window)))

				// t1 != t2 (super-klugey, but verifiable for seed == 3)
				Expect(t1).NotTo(BeTemporally("==", t2))
			})
		})

		When("Given an invalid duration", func() {
			It("Returns a nil generator and throws an error", func() {
				t0 := today()
				distro, err := burrow.UniformTimestampDistribution(t0, -1*time.Hour)
				Expect(err).To(HaveOccurred())
				Expect(distro).To(BeNil())

				distro, err = burrow.UniformTimestampDistribution(t0, 0*time.Hour)
				Expect(err).To(HaveOccurred())
				Expect(distro).To(BeNil())

			})
		})
	})

	Context("GaussianTimestampDistribution", func() {
		var (
			nSamples int = 100
		)
		When("Given a valid time tMu and a standard deviation tSigma", func() {
			It("Returns a gaussian distro function of distribution type N(tMu, tSigma)", func() {
				t0 := today()
				tMu := t0.Add(11 * time.Hour)
				tSigma := 2 * time.Hour

				distro, err := burrow.GaussianTimestampDistribution(tMu, tSigma)
				Expect(err).NotTo(HaveOccurred())
				Expect(distro).NotTo(BeNil())

				By("Failing to reject the null hypothesis")

				samples := make([]float64, 0, nSamples)
				for i := 0; i < nSamples; i++ {
					samples = append(samples, float64(distro().UnixMilli()))
				}

				pValue := oneSampleTtest(samples, float64(tMu.UnixMilli()))

				Expect(pValue).To(And(BeNumerically(">=", 0.05), BeNumerically("<=", 1.0)))
			})
		})

		When("Given a negative standard deviation", func() {
			It("Errors out with a message", func() {
				t0 := today()
				tMu := t0.Add(11 * time.Hour)
				tSigma := 2 * time.Hour

				distro, err := burrow.GaussianTimestampDistribution(tMu, -tSigma)
				Expect(err).To(MatchError("Standard deviation should not be negative"))
				Expect(distro).To(BeNil())
			})
		})
	})
})
