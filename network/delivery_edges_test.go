package network_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow/matchers"
	"github.com/bdshroyer/burrow/network"
)

var _ = Describe("DeliveryEdges", func() {
	testEdges := func() *network.DeliveryEdges {
		nodes := make(map[int64]*network.StopNode)
		for val := 1; val < 6; val++ {
			nodes[int64(val)] = &network.StopNode{Val: int64(val)}
		}

		edges := &network.DeliveryEdges{
			Payload: []*network.DeliveryEdge{
				&network.DeliveryEdge{Src: nodes[1], Dst: nodes[2], Wgt: 1.0},
				&network.DeliveryEdge{Src: nodes[1], Dst: nodes[3], Wgt: 1.4},
				&network.DeliveryEdge{Src: nodes[2], Dst: nodes[4], Wgt: 0.4},
				&network.DeliveryEdge{Src: nodes[4], Dst: nodes[5], Wgt: 2.7},
			},
			CurrentIdx: -1,
		}

		return edges
	}

	Describe("NewDeliveryEdges", func() {
		It("Returns a fresh delivery edges iterator", func() {
			Expect(network.NewDeliveryEdges()).NotTo(BeNil())
			Expect(network.NewDeliveryEdges().Len()).To(Equal(0))
		})
	})

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
				edges := &network.DeliveryEdges{Payload: []*network.DeliveryEdge{}}
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
				edges := &network.DeliveryEdges{Payload: []*network.DeliveryEdge{}}
				Expect(edges.WeightedEdge()).To(BeNil())
			})
		})
	})

	Describe("Current", func() {
		When("Iterator is not yet exhausted", func() {
			It("Returns the edge currently pointed to by the iterator", func() {
				edges := testEdges()
				edges.CurrentIdx = 2
				Expect(edges.Current().From().ID()).To(BeEquivalentTo(2))
				Expect(edges.Current().To().ID()).To(BeEquivalentTo(4))
				Expect(edges.Current().Weight()).To(BeEquivalentTo(0.4))
			})
		})

		When("Iterator is exhausted", func() {
			It("Returns a nil", func() {
				edges := testEdges()
				edges.CurrentIdx = 4
				Expect(edges.Current()).To(BeNil())
			})
		})

		When("Iterator is empty", func() {
			It("Returns a nil", func() {
				edges := &network.DeliveryEdges{Payload: []*network.DeliveryEdge{}}
				Expect(edges.Current()).To(BeNil())
			})
		})
	})

	Describe("Next", func() {
		When("Another item exists", func() {
			It("Returns the next item in the iterator", func() {
				edges := testEdges()
				edges.Next() // initialize the iterator
				Expect(edges.Next()).To(BeTrue())

				Expect(edges.WeightedEdge().From().ID()).To(BeEquivalentTo(1))
				Expect(edges.WeightedEdge().To().ID()).To(BeEquivalentTo(3))
				Expect(edges.WeightedEdge().Weight()).To(BeEquivalentTo(1.4))
			})
		})

		When("On the last item in the collection", func() {
			It("Returns true", func() {
				edges := testEdges()
				edges.Next() // initialize iterator
				Expect(edges.Next()).To(BeTrue())
				Expect(edges.Next()).To(BeTrue())
				Expect(edges.Next()).To(BeTrue())

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
				edges := &network.DeliveryEdges{Payload: []*network.DeliveryEdge{}}
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
				edges.Next() // initialize iterator
				edges.Next() // traverse one step

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
				edges := &network.DeliveryEdges{Payload: []*network.DeliveryEdge{}}

				Expect(edges.Len()).To(Equal(0))

				edges.Next()
				Expect(edges.Len()).To(Equal(0))
			})
		})
	})

	Describe("Reset", func() {
		It("Returns iterator to the start of the collection", func() {
			edges := testEdges()
			edges.Next() // initialize iterator

			firstEdge := edges.WeightedEdge()
			firstRemainder := edges.Len()

			edges.Next()

			edges.Reset()
			Expect(edges.Len()).To(Equal(firstRemainder))

			edges.Next() // reinitialize iterator
			Expect(edges.WeightedEdge()).To(matchers.MatchEdge(firstEdge))
		})
	})
})
