package network_test

import (
	"github.com/bdshroyer/burrow/network"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func testNewDeliveryNodes() *network.DeliveryNodes {
	return &network.DeliveryNodes{
		Payload: []network.DeliveryNode{
			&network.HubNode{Val: 1},
			&network.StopNode{Val: 4},
			&network.StopNode{Val: 3},
		},
		CurrentIdx: -1,
	}

}

var _ = Describe("DeliveryNodes", func() {
	var nodes *network.DeliveryNodes

	Context("NewDeliveryNodes", func() {
		It("Returns a new DeliveryNodes structure", func() {
			Expect(network.NewDeliveryNodes()).NotTo(BeNil())
			Expect(network.NewDeliveryNodes().Len()).To(BeZero())
		})
	})

	Context("Methods", func() {
		BeforeEach(func() {
			nodes = testNewDeliveryNodes()
		})

		Context("Node", func() {
			It("Returns the current node pointed to by the iterator", func() {
				nodes.Next()
				Expect(nodes.Node().ID()).To(BeEquivalentTo(1))
			})

			It("Returns nil if the collection is empty", func() {
				nodes := new(network.DeliveryNodes)
				nodes.Next()
				Expect(nodes.Node()).To(BeNil())
			})

			It("Returns nil if the iterator has been exhausted", func() {
				for nodes.Next() {
				}
				nodes.Next()

				Expect(nodes.Node()).To(BeNil())
			})

			It("Returns nil if the iterator hasn't been initialized by calling Next()", func() {
				Expect(nodes.Node()).To(BeNil())
			})
		})

		Context("Current", func() {
			It("Returns the current node pointed to by the iterator", func() {
				nodes.Next()
				Expect(nodes.Current().ID()).To(BeEquivalentTo(1))
			})

			It("Returns nil if the collection is empty", func() {
				nodes := new(network.DeliveryNodes)
				nodes.Next()
				Expect(nodes.Current()).To(BeNil())
			})

			It("Returns nil if the iterator has been exhausted", func() {
				for nodes.Next() {
				}
				nodes.Next()

				Expect(nodes.Current()).To(BeNil())
			})

			It("Returns nil if the iterator hasn't been initializedby calling Next()", func() {
				Expect(nodes.Node()).To(BeNil())
			})
		})

		Context("Len", func() {
			It("Returns the number of nodes in a new collection", func() {
				nodes.Next()
				Expect(nodes.Len()).To(Equal(3))
			})

			It("Returns the number of nodes remaining in a partially traversed collection", func() {
				nodes.Next()
				nodes.Next()
				Expect(nodes.Len()).To(Equal(2))
			})

			It("Returns 0 on an exhausted iterator in perpetuity", func() {
				for nodes.Next() {
				}
				nodes.Next()

				Expect(nodes.Len()).To(Equal(0))

				nodes.Next()
				Expect(nodes.Len()).To(Equal(0))
			})

			It("Returns nil on an empty collection", func() {
				nodes := new(network.DeliveryNodes)
				Expect(nodes.Len()).To(Equal(0))
			})

			It("Returns the full number of nodes on a collection that hasn't been initialized", func() {
				Expect(nodes.Len()).To(Equal(3))
			})
		})

		Context("Next", func() {
			It("Advances the iterator, returning true iff the iterator hasn't traversed the full collection", func() {
				nodes.Next()

				Expect(nodes.Node().ID()).To(BeEquivalentTo(1))

				// true because one more node remains after the second one
				Expect(nodes.Next()).To(BeTrue())
				Expect(nodes.Node().ID()).To(BeEquivalentTo(4))

				// true because the third node is the last one
				Expect(nodes.Next()).To(BeTrue())
				Expect(nodes.Node().ID()).To(BeEquivalentTo(3))
			})

			It("Returns false if an attempt is made to advance past the end of the collection", func() {
				// Iterate over the three elements in the list
				nodes.Next()
				nodes.Next()
				nodes.Next()

				// The fourth and subsequent calls return false
				Expect(nodes.Next()).To(BeFalse())
				Expect(nodes.Next()).To(BeFalse())
			})

			It("Simply returns false on an empty collection", func() {
				nodes := new(network.DeliveryNodes)
				Expect(nodes.Next()).To(BeEquivalentTo(false))
			})
		})

		Context("Reset", func() {
			It("Returns an iterator to its initial position within the collection", func() {
				nodes.Next()
				firstNode := nodes.Node()
				nodes.Next()

				nodes.Reset()
				nodes.Next()
				Expect(nodes.Node().ID()).To(BeEquivalentTo(firstNode.ID()))
			})

			It("Is idempotent", func() {
				nodes.Next()
				firstNode := nodes.Node()

				nodes.Reset()
				nodes.Next()

				Expect(nodes.Node().ID()).To(BeEquivalentTo(firstNode.ID()))

				nodes.Next()
				nodes.Reset()
				nodes.Next()
				Expect(nodes.Node().ID()).To(BeEquivalentTo(firstNode.ID()))
			})
		})
	})
})
