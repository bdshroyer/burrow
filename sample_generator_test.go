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

func MakeTestGenerator[T burrow.Rangeable] (distro burrow.SampleDistribution[T]) *burrow.SampleGenerator[T] {
	return &burrow.SampleGenerator[T]{
		Distro: distro,
		Quit:   make(chan bool),
	}
}

var _ = Describe("SampleGenerator", func() {
	sampleDistro := burrow.SampleDistribution[float64](testDistro)

	BeforeEach(func() {
		rand.Seed(3)
	})

	Context("NewSampleGenerator", func() {
		When("Called on a valid sample distribution", func() {
			It("Returns a new distribution generator struct", func() {
				generator, err := burrow.NewSampleGenerator(sampleDistro)
				Expect(err).NotTo(HaveOccurred())
				Expect(generator).NotTo(BeNil())
				Expect(generator.Quit).NotTo(BeClosed())
			})
		})

		When("Called on a nilsample distribution", func() {
			It("returns nil and raises an error", func() {
				generator, err := burrow.NewSampleGenerator[float64](nil)
				Expect(err).To(MatchError("Must receive a non-null sample distribution."))
				Expect(generator).To(BeNil())
			})
		})
	})

	Context("SampleDistribution generator", func() {
		When("Invoked with a positive number", func() {
			It("Returns the same number of samples drawn from the given distribution", func() {
				generator := MakeTestGenerator(sampleDistro)
				sample := generator.Sample(2)

				Eventually(sample).Should(Receive(BeNumerically("~", 0.71998, 1e-4)))
				Eventually(sample).Should(Receive(BeNumerically("~", 0.65263, 1e-4)))
				Eventually(sample).Should(BeClosed())
			})
		})

		When("Invoked with a zero", func() {
			It("Closes without returning any numbers", func() {
				generator := MakeTestGenerator(sampleDistro)
				sample := generator.Sample(0)

				Consistently(sample).ShouldNot(Receive())
				Eventually(sample).Should(BeClosed())
			})
		})

		When("The Stop() command is called", func() {
			It("Stops generating new samples and exits", func() {
				generator := MakeTestGenerator(sampleDistro)
				sample := generator.Sample(4)

				Eventually(sample).Should(Receive())

				generator.Stop()

				Consistently(sample).ShouldNot(Receive())
				Eventually(sample).Should(BeClosed())
			})
		})
	})
})
