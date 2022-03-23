package burrow_test

import (
	"math/rand"

	"github.com/bdshroyer/burrow"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func testDistro() float64 {
	return rand.Float64()
}

var _ = Describe("Distros", func() {
	// first two numbers returned by rand.Float64() with seed 3

	BeforeEach(func() {
		rand.Seed(3)
	})

	Context("SampleDistribution generator", func() {
		When("Invoked with a positive number", func() {
			It("Returns the same number of samples drawn from the given distribution", func() {
				sample := burrow.SampleDistribution(testDistro).Sample(2)
				Eventually(sample).Should(Receive(BeNumerically("~", 0.71998, 1e-4)))
				Eventually(sample).Should(Receive(BeNumerically("~", 0.65263, 1e-4)))
				Eventually(sample).Should(BeClosed())
			})
		})

		When("Invoked with a zero", func() {
			It("Closes without returning any numbers", func() {
				sample := burrow.SampleDistribution(testDistro).Sample(0)
				Consistently(sample).ShouldNot(Receive())
				Eventually(sample).Should(BeClosed())
			})
		})
	})

	Context("MakeUniformDistribution", func() {
		When("Given a positive range", func() {
			It("Returns a sample distribution", func() {
				distro, err := burrow.MakeUniformDistribution(5.0)
				Expect(err).NotTo(HaveOccurred())

				sample := distro.Sample(2)
				Eventually(sample).Should(Receive(BeNumerically("~", 5*0.71998, 1e-4)))
				Eventually(sample).Should(Receive(BeNumerically("~", 5*0.65263, 1e-4)))
				Eventually(sample).Should(BeClosed())
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
})
