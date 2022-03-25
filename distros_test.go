package burrow_test

import (
	"math/rand"
	"time"

	"github.com/bdshroyer/burrow"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func today() time.Time {
	base := time.Now()
	return time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, base.Location())
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
})
