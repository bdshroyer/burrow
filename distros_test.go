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
	Context("SampleGenerator", func() {
		BeforeEach(func() {
			rand.Seed(3)
		})

		When("Invoked with a positive number", func() {
			It("Returns the same number of samples drawn from the given distribution", func() {
				sample := burrow.SampleDistribution(testDistro).Sample(2)
				Eventually(sample).Should(Receive(BeNumerically("~", 0.71998, 1e-5)))
				Eventually(sample).Should(Receive(BeNumerically("~", 0.65263, 1e-5)))
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
})
