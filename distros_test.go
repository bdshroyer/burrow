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

func MakeTestGenerator[T burrow.Rangeable] (distro burrow.SampleDistribution[T]) *burrow.DistroGenerator[T] {
	return &burrow.DistroGenerator[T]{
		Distro: distro,
		Quit:   make(chan bool),
	}
}

var _ = Describe("Distros", func() {
	sampleDistro := burrow.SampleDistribution[float64](testDistro)

	BeforeEach(func() {
		rand.Seed(3)
	})

	Context("NewDistroGenerator", func() {
		When("Called on a valid sample distribution", func() {
			It("Returns a new distribution generator struct", func() {
				generator, err := burrow.NewDistroGenerator(sampleDistro)
				Expect(err).NotTo(HaveOccurred())
				Expect(generator).NotTo(BeNil())
				Expect(generator.Quit).NotTo(BeClosed())
			})
		})

		When("Called on a nilsample distribution", func() {
			It("returns nil and raises an error", func() {
				generator, err := burrow.NewDistroGenerator[float64](nil)
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
})
