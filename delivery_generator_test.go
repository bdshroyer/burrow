package burrow_test

import (
	"math/rand"
	"time"
	"container/heap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow"
	"github.com/bdshroyer/burrow/matchers"
)

func testNewStopFactory(counterSeed int64) *burrow.StopFactory {
	return &burrow.StopFactory{Counter: counterSeed}
}

func testTimeDist(t0 time.Time, window time.Duration) burrow.SampleDistribution[time.Time] {
	payload := func() time.Time {
		rangeSample := time.Duration(rand.Int63n(int64(window)))
		return t0.Add(rangeSample)
	}

	return burrow.SampleDistribution[time.Time](payload)
}

var _ = Describe("DeliveryGenerator", func() {
	Describe("StopNodeHeap", func() {

		var (
			rawHeap *burrow.StopNodeHeap
			Payload []*burrow.StopNode
		)

		BeforeEach(func() {
			Payload = []*burrow.StopNode{
						&burrow.StopNode{Val: 1, Timestamp: time.Date(2022, 3, 29, 16, 11, 8, 0, time.UTC)},
						&burrow.StopNode{Val: 2, Timestamp: time.Date(2022, 3, 29, 4, 32, 19, 0, time.UTC)},
						&burrow.StopNode{Val: 3, Timestamp: time.Date(2022, 3, 29, 5, 51, 26, 0, time.UTC)},
						&burrow.StopNode{Val: 4, Timestamp: time.Date(2022, 3, 29, 5, 51, 26, 0, time.UTC)},
					}
			rawHeap = new(burrow.StopNodeHeap)
			rawHeap.Payload = Payload
		})

		Describe("Len", func() {
			When("Called on a non-empty heap", func() {
				It("Returns the length of the heap", func() {
					Expect(rawHeap.Len()).To(Equal(4))
				})
			})

			When("Called on an empty heap", func() {
				rawHeap := new(burrow.StopNodeHeap)
				Expect(rawHeap.Len()).To(Equal(0))
			})
		})

		Describe("Less", func() {
			It("Returns true if item i is less than item j", func() {
				Expect(rawHeap.Less(0, 1)).To(BeFalse()) // t < u
				Expect(rawHeap.Less(1, 2)).To(BeTrue())  // t > u
				Expect(rawHeap.Less(2, 3)).To(BeFalse()) // t = u
			})
		})

		Describe("Swap", func() {
			It("Swaps the elements' positions in the heap", func() {
				rawHeap2 := &burrow.StopNodeHeap{ Payload: make([]*burrow.StopNode, rawHeap.Len()) }
				copy(rawHeap2.Payload, rawHeap.Payload)

				rawHeap2.Swap(0, 3)
				Expect(rawHeap2.Payload[0].Timestamp).To(BeTemporally("==", rawHeap.Payload[3].Timestamp))
				Expect(rawHeap2.Payload[3].Timestamp).To(BeTemporally("==", rawHeap.Payload[0].Timestamp))

				rawHeap2.Swap(2, 2)
				Expect(rawHeap2.Payload[2].Timestamp).To(BeTemporally("==", rawHeap.Payload[2].Timestamp))
			})
		})

		Describe("Push", func() {
			It("Appends a new node to the heap", func() {
				originalLength := rawHeap.Len()
				newNode := &burrow.StopNode{Val: 5, Timestamp: time.Date(2022, 3, 29, 21, 48, 54, 0, time.UTC)}

				rawHeap.Push(newNode)
				Expect(rawHeap.Len()).To(Equal(originalLength + 1))
				Expect(rawHeap.Payload).To(ContainElement(matchers.MatchNode(newNode)))
			})
		})

		Describe("Pop", func() {
			It("Removes and returns the top item from the heap", func() {
				originalLength := rawHeap.Len()
				originalNode := rawHeap.Payload[0]
				nextNode := rawHeap.Payload[1]

				popNode := rawHeap.Pop()

				Expect(rawHeap.Len()).To(Equal(originalLength - 1))
				Expect(popNode).To(matchers.MatchNode(originalNode))
				Expect(rawHeap.Payload[0]).To(matchers.MatchNode(nextNode))
			})
		})

		It("Observes the min-heap property", func() {
			testHeap := new(burrow.StopNodeHeap)
			heap.Init(testHeap)

			for i := 0; i < rawHeap.Len(); i++ {
				testHeap.Push(rawHeap.Payload[i])
			}

			Expect(testHeap.Len()).To(Equal(rawHeap.Len()))
			prev := testHeap.Pop().(*burrow.StopNode)

			for i := 1; i < rawHeap.Len(); i++ {
				current := testHeap.Pop().(*burrow.StopNode)

				Expect(prev.Timestamp).To(BeTemporally("<=", current.Timestamp))
				prev = current
			}
		})
	})

	Context("StopFactory", func() {
		Context("NewStopFactory", func() {
			It("Returns a StopFactory struct with a counter of 1", func() {
				sg := burrow.NewStopFactory()
				Expect(sg).NotTo(BeNil())
				Expect(sg.Counter).To(BeEquivalentTo(1))
			})
		})

		When("MakeStop is called", func() {
			It("Returns a new StopNode with an ID matching the generator's counter", func() {
				sg := testNewStopFactory(1)
				firstCount := sg.Counter
				stamp := time.Now()

				newNode := sg.MakeStop(stamp)
				Expect(newNode).NotTo(BeNil())
				Expect(newNode.ID()).To(BeEquivalentTo(firstCount))
				Expect(newNode.Timestamp).To(BeEquivalentTo(stamp))
			})

			It("Increments the generator's counter", func() {
				sg := testNewStopFactory(1)
				firstCount := sg.Counter
				stamp := time.Now()

				sg.MakeStop(stamp)
				Expect(sg.Counter).To(BeEquivalentTo(firstCount + 1))
			})
		})
	})

	Context("MakeStopDAG", func() {
		When("Invoked with a positive integer", func() {
			It("Returns a DAG of stop nodes", func() {
				window := time.Duration(rand.Int63n(int64(24 * time.Hour)))
				today, err := time.Parse(time.RFC3339, "2022-03-25T00:00:00-04:00")
				Expect(err).NotTo(HaveOccurred())

				dag, err := burrow.MakeStopDAG(3, testTimeDist(today, window))
				Expect(err).NotTo(HaveOccurred())
				Expect(dag.Nodes().Len()).To(Equal(3))
			})

			It("Has edges that comply with the happens-before relation", func() {
				window := 24 * time.Hour
				today, err := time.Parse(time.RFC3339, "2022-03-25T00:00:00-04:00")
				Expect(err).NotTo(HaveOccurred())

				dag, err := burrow.MakeStopDAG(3, testTimeDist(today, window))
				Expect(err).NotTo(HaveOccurred())

				edges := dag.Edges().(*burrow.DeliveryEdges)

				// Edges should be created if timestamps are not identical.
				Expect(edges.Len()).To(BeNumerically(">", 0))

				// Each edge must have the property that the source is earlier than the destination.
				for edges.Next() {
					from := edges.Current().From().(*burrow.StopNode)
					to := edges.Current().To().(*burrow.StopNode)

					Expect(from.Timestamp).To(BeTemporally("<", to.Timestamp))
				}
			})

			Context("When one or more stops occurs at the same time", func() {
			})
		})

		When("Invoked with a zero", func() {
			It("Creats an empty DAG", func() {
				window := time.Duration(rand.Int63n(int64(24 * time.Hour)))
				today, err := time.Parse(time.RFC3339, "2022-03-25T00:00:00-04:00")
				Expect(err).NotTo(HaveOccurred())

				dag, err := burrow.MakeStopDAG(0, testTimeDist(today, window))
				Expect(err).NotTo(HaveOccurred())
				Expect(dag.Nodes().Len()).To(Equal(0))
			})
		})

		When("Passed an empty distribution", func() {
			It("Returns an error", func() {
				dag, err := burrow.MakeStopDAG(3, nil)
				Expect(err).To(MatchError("Must receive a non-null sample distribution."))
				Expect(dag).To(BeNil())
			})
		})
	})
})
