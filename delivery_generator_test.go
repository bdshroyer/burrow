package burrow_test

import (
	"math/rand"
	"time"
	"sort"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow"
)

func testNewNodeFactory(counterSeed int64) *burrow.NodeFactory {
	return &burrow.NodeFactory{Counter: counterSeed}
}

func testTimeDist(t0 time.Time, window time.Duration) burrow.SampleDistribution[time.Time] {
	payload := func() time.Time {
		rangeSample := time.Duration(rand.Int63n(int64(window)))
		return t0.Add(rangeSample)
	}

	return burrow.SampleDistribution[time.Time](payload)
}

const window time.Duration = 24 * time.Hour

var _ = Describe("DeliveryGenerator", func() {
	Describe("StopNodeSortable", func() {

		var (
			rawSortable *burrow.StopNodeSortable
			Payload []*burrow.StopNode
		)

		BeforeEach(func() {
			Payload = []*burrow.StopNode{
						&burrow.StopNode{Val: 1, Timestamp: time.Date(2022, 3, 29, 16, 11, 8, 0, time.UTC)},
						&burrow.StopNode{Val: 2, Timestamp: time.Date(2022, 3, 29, 4, 32, 19, 0, time.UTC)},
						&burrow.StopNode{Val: 3, Timestamp: time.Date(2022, 3, 29, 5, 51, 26, 0, time.UTC)},
						&burrow.StopNode{Val: 4, Timestamp: time.Date(2022, 3, 29, 5, 51, 26, 0, time.UTC)},
					}
			rawSortable = new(burrow.StopNodeSortable)
			rawSortable.Payload = Payload
		})

		Describe("Len", func() {
			When("Called on a non-empty heap", func() {
				It("Returns the length of the heap", func() {
					Expect(rawSortable.Len()).To(Equal(4))
				})
			})

			When("Called on an empty heap", func() {
				rawSortable := new(burrow.StopNodeSortable)
				Expect(rawSortable.Len()).To(Equal(0))
			})
		})

		Describe("Less", func() {
			It("Returns true if item i is less than item j", func() {
				Expect(rawSortable.Less(0, 1)).To(BeFalse()) // t < u
				Expect(rawSortable.Less(1, 2)).To(BeTrue())  // t > u
				Expect(rawSortable.Less(2, 3)).To(BeFalse()) // t = u
			})
		})

		Describe("Swap", func() {
			It("Swaps the elements' positions in the heap", func() {
				rawSortable2 := &burrow.StopNodeSortable{ Payload: make([]*burrow.StopNode, rawSortable.Len()) }
				copy(rawSortable2.Payload, rawSortable.Payload)

				rawSortable2.Swap(0, 3)
				Expect(rawSortable2.Payload[0].Timestamp).To(BeTemporally("==", rawSortable.Payload[3].Timestamp))
				Expect(rawSortable2.Payload[3].Timestamp).To(BeTemporally("==", rawSortable.Payload[0].Timestamp))

				rawSortable2.Swap(2, 2)
				Expect(rawSortable2.Payload[2].Timestamp).To(BeTemporally("==", rawSortable.Payload[2].Timestamp))
			})
		})

		It("Produces an ordered list when passed to the sort.Sort() function", func() {
			sort.Sort(rawSortable)

			for i := 1; i < rawSortable.Len(); i++ {
				Expect(rawSortable.Payload[i].Timestamp).To(BeTemporally(">=", rawSortable.Payload[i-1].Timestamp))
			}
		})
	})

	Describe("sortInPlace", func() {
		It("Sorts the passed-in array of stop nodes", func() {
			payload := []*burrow.StopNode{
				&burrow.StopNode{Val: 1, Timestamp: time.Date(2022, 3, 29, 16, 11, 8, 0, time.UTC)},
				&burrow.StopNode{Val: 2, Timestamp: time.Date(2022, 3, 29, 4, 32, 19, 0, time.UTC)},
				&burrow.StopNode{Val: 3, Timestamp: time.Date(2022, 3, 29, 5, 51, 26, 0, time.UTC)},
				&burrow.StopNode{Val: 4, Timestamp: time.Date(2022, 3, 29, 5, 51, 26, 0, time.UTC)},
			}

			burrow.SortInPlace(payload)

			for i := 1; i < len(payload); i++ {
				Expect(payload[i].Timestamp).To(BeTemporally(">=", payload[i-1].Timestamp))
			}
		})
	})

	Context("NodeFactory", func() {
		Context("NewNodeFactory", func() {
			It("Returns a NodeFactory struct with a counter of 1", func() {
				nf := burrow.NewNodeFactory()
				Expect(nf).NotTo(BeNil())
				Expect(nf.Counter).To(BeEquivalentTo(1))
			})
		})

		When("MakeStop is called", func() {
			It("Returns a new StopNode with an ID matching the generator's counter", func() {
				nf := testNewNodeFactory(1)
				firstCount := nf.Counter
				stamp := time.Now()

				newNode := nf.MakeStop(stamp)
				Expect(newNode).NotTo(BeNil())
				Expect(newNode.ID()).To(Equal(firstCount))
				Expect(newNode.Timestamp).To(Equal(stamp))
			})

			It("Increments the generator's counter", func() {
				nf := testNewNodeFactory(1)
				firstCount := nf.Counter
				stamp := time.Now()

				nf.MakeStop(stamp)
				Expect(nf.Counter).To(Equal(firstCount + 1))
			})
		})

		When("MakeHub is called", func() {
			It("Increments the generator's counter", func() {
				nf := testNewNodeFactory(1)
				firstCount := nf.Counter

				newHub := nf.MakeHub()
				Expect(newHub).NotTo(BeNil())
				Expect(newHub.ID()).To(Equal(firstCount))
				Expect(nf.Counter).To(Equal(firstCount + 1))
			})
		})
	})

	Context("MakeDeliveryNetwork", func() {
		var (
			today time.Time
			cfg burrow.DeliveryNetworkConfig
			err error
		)

		BeforeEach(func() {
				today, err = time.Parse(time.RFC3339, "2022-03-25T00:00:00-04:00")
				Expect(err).NotTo(HaveOccurred())

				cfg = burrow.DeliveryNetworkConfig {
					HubNodes: 2,
					StopNodes: 3,
					Distro: testTimeDist(today, window),
				}
		})

		When("Given a distro and a non-zero number of stop and hub nodes", func() {
			It("Creates a delivery network", func() {
				G, err := burrow.MakeDeliveryNetwork(cfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(G).NotTo(BeNil())

				Expect(len(G.Hubs)).To(Equal(2))
				Expect(len(G.Stops)).To(Equal(3))

				// hub->stop + stop->hub + stop->stop
				//   = 2 * hubs * stops + (stops-1) * stops / 2
				//   = 2 * 2 * 3 + (3-1) * 3 / 2
				//   = 15
				Expect(G.Edges().Len()).To(Equal(15))
			})

			It("Has stop-to-stop edges that all comply with the happens-before relation", func() {
				G, err := burrow.MakeDeliveryNetwork(cfg)
				Expect(err).NotTo(HaveOccurred())

				edges := G.Edges().(*burrow.DeliveryEdges)
				stopEdgeCount := 0

				// Each edge must have the property that the source is earlier than the destination.
				for edges.Next() {
					from := edges.Current().From().(burrow.DeliveryNode)
					to := edges.Current().To().(burrow.DeliveryNode)

					if from.IsHub() || to.IsHub() {
						continue
					}

					stopEdgeCount++
					fromStop := from.(*burrow.StopNode)
					toStop := to.(*burrow.StopNode)
					Expect(fromStop.Timestamp).To(BeTemporally("<", toStop.Timestamp))
				}

				// Edges should be created if timestamps are not identical.
				Expect(stopEdgeCount).To(BeNumerically(">", 0))
			})
		})

		When("Invoked with a hub count of 0", func() {
			It("Creates a stop-node DAG", func() {
				cfg.HubNodes = 0

				dag, err := burrow.MakeDeliveryNetwork(cfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(dag.Hubs)).To(Equal(0))
				Expect(len(dag.Stops)).To(Equal(3))
				Expect(dag.Edges().Len()).To(Equal(3))
			})
		})

		When("Passed an empty distribution", func() {
			It("Returns an error", func() {
				cfg.Distro = nil

				dag, err := burrow.MakeDeliveryNetwork(cfg)
				Expect(err).To(MatchError("Must receive a non-null sample distribution."))
				Expect(dag).To(BeNil())
			})
		})
	})
})
