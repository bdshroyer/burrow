package testutils_test

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gonum.org/v1/gonum/stat/distuv"

	"github.com/bdshroyer/burrow/testutils"
)

// From Section 3 of Marsaglia & Marsaglia '04.
func testLeftCrossing(N float64) float64 {
	return 0.01265 + 0.1757/N
}

var _ = Describe("AndersonDarling", func() {
	// Values sourced from Marsaglia and Marsaglia, indicating the 90%, 95%, and 99%
	// confidence thresholds for the Anderson-Darling a-statistic.
	const AD90 float64 = 1.933
	const AD95 float64 = 2.492
	const AD99 float64 = 3.880

	Describe("AndersonDarlingTest", func() {
		var mu, sigma float64 = 0.0, 1.0

		When("Called on a Gaussian distributed sample", func() {
			It("Returns a p-value less than 0.95", func() {
				dist := distuv.Normal{Mu: mu, Sigma: sigma}

				samples := make([]float64, 100)
				for i := 0; i < 100; i++ {
					samples[i] = dist.Rand()
				}

				pValue, err := testutils.AndersonDarlingTest(samples)
				Expect(err).NotTo(HaveOccurred())
				Expect(pValue).To(And(BeNumerically(">=", 0.0), BeNumerically("<", 0.95)))
			})

			It("Is not thrown by non-normalized distributions", func() {
				dist := distuv.Normal{Mu: mu + 1.0, Sigma: sigma + 3.0}

				samples := make([]float64, 100)
				for i := 0; i < 100; i++ {
					samples[i] = dist.Rand()
				}

				pValue, err := testutils.AndersonDarlingTest(samples)
				Expect(err).NotTo(HaveOccurred())
				Expect(pValue).To(And(BeNumerically(">=", 0.0), BeNumerically("<", 0.95)))
			})
		})

		When("Called on a non-Gaussian distributed case", func() {
			It("Returns a p-value greater than or equal to 0.95", func() {
				dist := distuv.LogNormal{Mu: mu, Sigma: sigma}

				samples := make([]float64, 100)
				for i := 0; i < 100; i++ {
					samples[i] = dist.Rand()
				}

				pValue, err := testutils.AndersonDarlingTest(samples)
				Expect(err).NotTo(HaveOccurred())
				Expect(pValue).To(BeNumerically(">=", 0.95))
			})
		})

		When("Called on an empty sample", func() {
			It("Returns a -1 and throws an error", func() {
				pValue, err := testutils.AndersonDarlingTest([]float64{})
				Expect(pValue).To(BeEquivalentTo(-1))
				Expect(err).To(MatchError("Cannot compute a-statistic on an empty sample"))
			})
		})
	})

	Describe("ADStatistic", func() {
		var mu, sigma float64 = 0.0, 1.0

		When("Given Gaussian distributed data", func() {
			It("Computes the A-Statistic", func() {
				dist := distuv.Normal{Mu: mu, Sigma: sigma}

				samples := make([]float64, 100)
				for i := 0; i < 100; i++ {
					samples[i] = dist.Rand()
				}

				aStat, err := testutils.ADStatistic(samples)
				Expect(err).NotTo(HaveOccurred())
				Expect(aStat).To(BeNumerically("<", AD95))
			})

			It("Is not thrown by non-centered distributions", func() {
				dist := distuv.Normal{Mu: mu + 5.0, Sigma: sigma + 1.5}

				samples := make([]float64, 100)
				for i := 0; i < 100; i++ {
					samples[i] = dist.Rand()
				}

				aStat, err := testutils.ADStatistic(samples)
				Expect(err).NotTo(HaveOccurred())
				Expect(aStat).To(BeNumerically("<", AD95))
			})
		})

		When("Given non-normally distributed data", func() {
			It("Computes the A-Statistic", func() {
				dist := distuv.LogNormal{Mu: mu, Sigma: sigma}

				samples := make([]float64, 100)
				for i := 0; i < 100; i++ {
					samples[i] = dist.Rand()
				}

				aStat, err := testutils.ADStatistic(samples)
				Expect(err).NotTo(HaveOccurred())
				Expect(aStat).To(BeNumerically(">=", AD95))
			})
		})

		When("Passed a uniform sample", func() {
			It("returns an infinite a-statistic", func() {
				samples := make([]float64, 100)
				for i := 0; i < 100; i++ {
					samples[i] = 0.1
				}

				aStat, err := testutils.ADStatistic(samples)
				Expect(err).NotTo(HaveOccurred())
				Expect(aStat).To(BeEquivalentTo(math.Inf(1)))
			})
		})

		When("Passed an empty sample", func() {
			It("Returns an error and a negative a-statistic", func() {
				samples := []float64{}

				aStat, err := testutils.ADStatistic(samples)
				Expect(err).To(MatchError("Cannot compute a-statistic on an empty sample"))
				Expect(aStat).To(BeEquivalentTo(-1))
			})
		})
	})

	Describe("ADPValue", func() {
		It("Approximates thresholds achieved by Anderson & Darling", func() {
			pValue, err := testutils.ADPValue(AD90)
			Expect(err).NotTo(HaveOccurred())

			Expect(pValue).To(BeNumerically("~", 0.90, 1e-3))

			pValue, err = testutils.ADPValue(2.492)
			Expect(err).NotTo(HaveOccurred())

			Expect(pValue).To(BeNumerically("~", 0.95, 1e-3))

			// This is using the threshold given by Sinclair and Spur as cited in
			// Marsaglia & Marsaglia, as it is more accuate than the Anderson-Darling
			// estimate.
			pValue, err = testutils.ADPValue(AD99)
			Expect(err).NotTo(HaveOccurred())

			Expect(pValue).To(BeNumerically("~", 0.99, 1e-4))
		})

		Context("A-statistic limits", func() {
			When("Called on an AD-statistic close to zero", func() {
				It("Returns a value close to 0", func() {
					pValue, err := testutils.ADPValue(math.SmallestNonzeroFloat64)
					Expect(err).NotTo(HaveOccurred())

					Expect(pValue).To(BeNumerically("~", 0.0, 1e-6))
				})
			})

			When("Called on a very large AD-statistic", func() {
				It("Returns a value close to one", func() {
					pValue, err := testutils.ADPValue(math.MaxFloat64)
					Expect(err).NotTo(HaveOccurred())

					Expect(pValue).To(BeNumerically("~", 1.0, 1e-6))
				})
			})
		})

		Context("Input limits", func() {
			When("Called on a zero", func() {
				It("Throws an error and returns -1", func() {
					pValue, err := testutils.ADPValue(0.0)
					Expect(err).To(MatchError("AD-statistic must be greater than 0"))
					Expect(pValue).To(Equal(-1.0))
				})
			})

			When("Called on infinity", func() {
				It("Returns a value close to one", func() {
					pValue, err := testutils.ADPValue(math.Inf(1))
					Expect(err).NotTo(HaveOccurred())

					Expect(pValue).To(BeNumerically("~", 1.0, 1e-6))
				})
			})
		})
	})

	Describe("ADSErrFix", func() {
		var Ntest []float64 = []float64{8.0, 16.0, 32.0, 64.0}

		// Marsaglia & Marsaglia's errfix approximation was calculated based on simulations
		// that showed two x-intercepts for the error curve of A for any given N. This set
		// of tests checks against that result. The target accuracy of 5e-5 is given by [1].
		When("Tested against the x-intercepts of the error curves", func() {
			It("Returns zero for a p-value of 0", func() {
				for i := 0; i < len(Ntest); i++ {
					fix, err := testutils.ADErrFix(Ntest[i], 0.0)

					Expect(err).NotTo(HaveOccurred())
					Expect(fix).To(BeNumerically("~", 0.0, 5e-5))
				}
			})

			It("Returns zero for the first x-intercept", func() {
				for i := 0; i < len(Ntest); i++ {
					xIntercept := testLeftCrossing(Ntest[i])
					fix, err := testutils.ADErrFix(Ntest[i], xIntercept)

					Expect(err).NotTo(HaveOccurred())
					Expect(fix).To(BeNumerically("~", 0.0, 5e-5))
				}
			})

			It("Returns zero for the second x-intercept", func() {
				for i := 0; i < len(Ntest); i++ {
					fix, err := testutils.ADErrFix(Ntest[i], 0.8)
					Expect(err).NotTo(HaveOccurred())
					Expect(fix).To(BeNumerically("~", 0.0, 5e-5))
				}
			})
		})

		// As N approaches infinity, the AD distribution should converge to its infinite version,
		// and the error correction term should vanish.
		When("An infinite sample", func() {
			It("Returns zero", func() {
				fix, err := testutils.ADErrFix(math.Inf(1), 0.8)
				Expect(err).NotTo(HaveOccurred())
				Expect(fix).To(BeZero())
			})
		})

		FIt("Does stuff", func() {
			pValue, err := testutils.ADPValue(AD95)
			Expect(err).NotTo(HaveOccurred())
			fix, err := testutils.ADErrFix(100.0, pValue)
			Expect(err).NotTo(HaveOccurred())
			fmt.Printf("AD95 at 100 = %f\n", fix)
		})

		When("Passed bad inputs", func() {
			It("Returns an error", func() {
				for i := 0; i < len(Ntest); i++ {
					fix, err := testutils.ADErrFix(Ntest[i], -1.0)
					Expect(fix).To(BeZero())
					Expect(err).To(MatchError("Probability must be between 0 and 1."))

					fix, err = testutils.ADErrFix(Ntest[i], 3.0)
					Expect(fix).To(BeZero())
					Expect(err).To(MatchError("Probability must be between 0 and 1."))

					fix, err = testutils.ADErrFix(0.0, 0.5)
					Expect(fix).To(BeZero())
					Expect(err).To(MatchError("N must be greater than 0."))
				}
			})
		})
	})
})
