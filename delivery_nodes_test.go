package burrow_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow"
)

var _ = Describe("DeliveryNodes", func() {
	var nodes *burrow.DeliveryNodes

	BeforeEach(func() {
		nodes = &burrow.DeliveryNodes{
			Payload: []burrow.DeliveryNode{
				&burrow.HubNode{Val: 1},
				&burrow.StopNode{Val: 4},
				&burrow.StopNode{Val: 3},
			},
		}
	})

	Context("Node", func() {
		It("Returns the current node pointed to by the iterator", func() {
			Expect(nodes.Node().ID()).To(BeEquivalentTo(1))
		})

		It("Returns nil if the collection is empty", func() {
			nodes := new(burrow.DeliveryNodes)
			Expect(nodes.Node()).To(BeNil())
		})

		It("Returns nil if the iterator has been exhausted", func() {
			for nodes.Next() {
			}
			nodes.Next()

			Expect(nodes.Node()).To(BeNil())
		})
	})

	Context("Len", func() {
		It("Returns the number of nodes in a new collection", func() {
			Expect(nodes.Len()).To(BeEquivalentTo(3))
		})

		It("Returns the number of nodes remaining in a partially traversed collection", func() {
			nodes.Next()
			Expect(nodes.Len()).To(BeEquivalentTo(2))
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
			nodes := new(burrow.DeliveryNodes)
			Expect(nodes.Len()).To(BeEquivalentTo(0))
		})
	})

	Context("Next", func() {
		It("Advances the iterator, returning true iff the iterator hasn't traversed the full collection", func() {

			Expect(nodes.Node().ID()).To(BeEquivalentTo(1))

			// true because one more node remains after the second one
			Expect(nodes.Next()).To(BeTrue())
			Expect(nodes.Node().ID()).To(BeEquivalentTo(4))

			// false because the third node is the last one
			Expect(nodes.Next()).To(BeFalse())
			Expect(nodes.Node().ID()).To(BeEquivalentTo(3))
		})

		It("Returns false if an attempt is made to advance past the end of the collection", func() {
			for nodes.Next() {
			}

			Expect(nodes.Next()).To(BeFalse())
		})

		It("Simply returns false on an empty collection", func() {
			nodes := new(burrow.DeliveryNodes)
			Expect(nodes.Next()).To(BeEquivalentTo(false))
		})
	})

	Context("Reset", func() {
		It("Returns an iterator to its initial position within the collection", func() {
			firstNode := nodes.Node()
			nodes.Next()

			nodes.Reset()
			Expect(nodes.Node().ID()).To(BeEquivalentTo(firstNode.ID()))
		})

		It("Is idempotent", func() {
			firstNode := nodes.Node()
			nodes.Reset()
			Expect(nodes.Node().ID()).To(BeEquivalentTo(firstNode.ID()))

			nodes.Next()
			nodes.Reset()
			Expect(nodes.Node().ID()).To(BeEquivalentTo(firstNode.ID()))
		})
	})
})
