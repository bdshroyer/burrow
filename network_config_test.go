package burrow_test

import (
	"math"
	"math/rand"
	"time"

	"github.com/bdshroyer/burrow"
	"github.com/bdshroyer/burrow/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gonum.org/v1/gonum/stat"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("NetworkConfig", func() {
	BeforeEach(func() {
		rand.Seed(3)
	})

	makeNetworkUniformDistro := func() *burrow.NetworkSpec_Uniform {
		return &burrow.NetworkSpec_Uniform{Uniform: &burrow.NetworkSpec_UniformDistro{}}
	}

	makeNetworkGaussianDistro := func(mean, stdDev int64) *burrow.NetworkSpec_Gaussian {
		return &burrow.NetworkSpec_Gaussian{
			Gaussian: &burrow.NetworkSpec_GaussianDistro{Mean: mean, StdDev: stdDev},
		}
	}

	Context("NewConfig", func() {
		var (
			spec         burrow.NetworkSpec
			tStart, tEnd *timestamppb.Timestamp
		)

		BeforeEach(func() {
			tStart = timestamppb.New(today())
			tEnd = timestamppb.New(today().Add(24 * time.Hour))

			spec = burrow.NetworkSpec{
				Hubs:         2,
				Stops:        5,
				Start:        tStart,
				End:          tEnd,
				ShortEdge:    durationpb.New(30 * time.Minute),
				LongEdge:     durationpb.New(6 * time.Hour),
				Distribution: makeNetworkUniformDistro(),
			}
		})

		When("Given a valid network spec", func() {
			It("Produces an artifact with the correct network specifications and edge boundaries", func() {
				cfg, err := burrow.NewNetworkConfig(spec)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.HubNodes).To(BeEquivalentTo(spec.Hubs))
				Expect(cfg.StopNodes).To(BeEquivalentTo(spec.Stops))
				Expect(cfg.EdgeBounds[0]).To(Equal(30 * time.Minute))
				Expect(cfg.EdgeBounds[1]).To(Equal(6 * time.Hour))
			})

			It("Returns a matching config on a uniform distro", func() {
				cfg, err := burrow.NewNetworkConfig(spec)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				varianceSample := make([]float64, 10000)

				for i := 0; i < len(varianceSample); i++ {
					varianceSample[i] = float64(cfg.Distro().UnixMilli())
				}

				By("Verifying the uniform distribution characteristics")

				mean, variance := stat.MeanVariance(varianceSample, nil)

				expectedMean := 0.5 * float64(tStart.AsTime().UnixMilli()+tEnd.AsTime().UnixMilli())
				expectedVariance := math.Pow(
					float64(tEnd.AsTime().UnixMilli()-tStart.AsTime().UnixMilli()),
					2.0,
				) / 12.0

				Expect((mean - expectedMean) / expectedMean).To(BeNumerically("~", 0, 5e-3))
				Expect((variance - expectedVariance) / expectedVariance).To(BeNumerically("~", 0, 5e-3))
			})

			It("Returns a matching config on a Gaussian distro", func() {
				mu, sigma := tStart.AsTime().Add(11*time.Hour), 2*time.Hour
				spec.Distribution = makeNetworkGaussianDistro(
					mu.UnixMicro(),
					sigma.Microseconds(),
				)

				cfg, err := burrow.NewNetworkConfig(spec)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				N := 10000
				varianceSample := make([]float64, N)

				for i := 0; i < len(varianceSample); i++ {
					varianceSample[i] = float64(cfg.Distro().UnixMilli())
				}

				By("Verifying the Gaussian distribution characteristics")

				pValue, err := testutils.AndersonDarlingTest(varianceSample)
				Expect(err).NotTo(HaveOccurred())
				Expect(pValue).To(And(BeNumerically(">=", 0.0), BeNumerically("<=", 0.95)))
			})
		})
	})
})
