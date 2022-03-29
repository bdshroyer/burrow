package burrow_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow"
	"github.com/bdshroyer/burrow/matchers"
)

var _ = Describe("DeliveryEdges", func() {
	testEdges := func() *burrow.DeliveryEdges {
		nodes := make(map[int64]*burrow.StopNode)
		for val := 1; val < 6; val++ {
			nodes[int64(val)] = &burrow.StopNode{Val: int64(val)}
		}

		edges := &burrow.DeliveryEdges{
			Payload: []*burrow.DeliveryEdge{
				&burrow.DeliveryEdge{Src: nodes[1], Dst: nodes[2], Wgt: 1.0},
				&burrow.DeliveryEdge{Src: nodes[1], Dst: nodes[3], Wgt: 1.4},
				&burrow.DeliveryEdge{Src: nodes[2], Dst: nodes[4], Wgt: 0.4},
				&burrow.DeliveryEdge{Src: nodes[4], Dst: nodes[5], Wgt: 2.7},
			},
		}

		return edges
	}

	Describe("Edge", func() {
		When("Iterator is not yet exhausted", func() {
			It("Returns the edge currently pointed to by the iterator", func() {
				edges := testEdges()
				edges.CurrentIdx = 2
				Expect(edges.Edge().From().ID()).To(BeEquivalentTo(2))
				Expect(edges.Edge().To().ID()).To(BeEquivalentTo(4))
			})
		})

		When("Iterator is exhausted", func() {
			It("Returns a nil", func() {
				edges := testEdges()
				edges.CurrentIdx = 4
				Expect(edges.Edge()).To(BeNil())
			})
		})

		When("Iterator is empty", func() {
			It("Returns a nil", func() {
				edges := &burrow.DeliveryEdges{Payload: []*burrow.DeliveryEdge{}}
				Expect(edges.Edge()).To(BeNil())
			})
		})
	})

	Describe("WeightedEdge", func() {
		When("Iterator is not yet exhausted", func() {
			It("Returns the edge currently pointed to by the iterator", func() {
				edges := testEdges()
				edges.CurrentIdx = 2
				Expect(edges.WeightedEdge().From().ID()).To(BeEquivalentTo(2))
				Expect(edges.WeightedEdge().To().ID()).To(BeEquivalentTo(4))
				Expect(edges.WeightedEdge().Weight()).To(BeEquivalentTo(0.4))
			})
		})

		When("Iterator is exhausted", func() {
			It("Returns a nil", func() {
				edges := testEdges()
				edges.CurrentIdx = 4
				Expect(edges.WeightedEdge()).To(BeNil())
			})
		})

		When("Iterator is empty", func() {
			It("Returns a nil", func() {
				edges := &burrow.DeliveryEdges{Payload: []*burrow.DeliveryEdge{}}
				Expect(edges.WeightedEdge()).To(BeNil())
			})
		})
	})

	Describe("Next", func() {
		When("Another item exists", func() {
			It("Returns the next item in the iterator", func() {
				edges := testEdges()
				Expect(edges.Next()).To(BeTrue())

				Expect(edges.WeightedEdge().From().ID()).To(BeEquivalentTo(1))
				Expect(edges.WeightedEdge().To().ID()).To(BeEquivalentTo(3))
				Expect(edges.WeightedEdge().Weight()).To(BeEquivalentTo(1.4))
			})
		})

		When("Only one more item remains in the iterator", func() {
			It("Returns false", func() {
				edges := testEdges()
				Expect(edges.Next()).To(BeTrue())
				Expect(edges.Next()).To(BeTrue())
				Expect(edges.Next()).To(BeFalse())

				Expect(edges.WeightedEdge().From().ID()).To(BeEquivalentTo(4))
				Expect(edges.WeightedEdge().To().ID()).To(BeEquivalentTo(5))
				Expect(edges.WeightedEdge().Weight()).To(BeEquivalentTo(2.7))
			})
		})

		When("When iterator is exhausted", func() {
			It("Returns false and iterator returns a nil", func() {
				edges := testEdges()
				for edges.Next() {
				}

				Expect(edges.Next()).To(BeFalse())
				Expect(edges.WeightedEdge()).To(BeNil())
			})
		})

		When("Iterator is empty", func() {
			It("Returns false and the iterator itself returns nil", func() {
				edges := &burrow.DeliveryEdges{Payload: []*burrow.DeliveryEdge{}}
				Expect(edges.Next()).To(BeFalse())
				Expect(edges.WeightedEdge()).To(BeNil())
			})
		})
	})

	Describe("Len", func() {
		When("Iterator is initialized", func() {
			It("Returns the number of items over which the iterator will traverse", func() {
				edges := testEdges()

				Expect(edges.Len()).To(Equal(4))
			})
		})

		When("Iterator is partially traversed", func() {
			It("Returns the number of items remaining", func() {
				edges := testEdges()
				edges.Next()

				Expect(edges.Len()).To(Equal(3))
			})
		})

		When("Iterator is exhausted", func() {
			It("Returns a 0 in perpetuity", func() {
				edges := testEdges()
				for edges.Next() {
				}
				edges.Next()

				Expect(edges.Len()).To(Equal(0))

				edges.Next()
				Expect(edges.Len()).To(Equal(0))
			})
		})

		When("Iterator is empty", func() {
			It("Returns a 0 in perpetuity", func() {
				edges := &burrow.DeliveryEdges{Payload: []*burrow.DeliveryEdge{}}

				Expect(edges.Len()).To(Equal(0))

				edges.Next()
				Expect(edges.Len()).To(Equal(0))
			})
		})
	})

	Describe("Reset", func() {
		It("Returns iterator to the start of the collection", func() {
			edges := testEdges()
			firstEdge := edges.WeightedEdge()
			firstRemainder := edges.Len()

			edges.Next()

			edges.Reset()
			Expect(edges.Len()).To(Equal(firstRemainder))
			Expect(edges.WeightedEdge()).To(matchers.MatchEdge(firstEdge))
		})
	})
})
