package network_test

import (
	"github.com/bdshroyer/burrow/network"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TestNode struct {
	Val int64
}

func (n *TestNode) ID() int64 {
	return n.Val
}

func (n *TestNode) IsHub() bool {
	return false
}

var _ = Describe("DeliveryEdge", func() {
	Describe("Edge", func() {
		It("Implements the Edge interface", func() {
			e := &network.DeliveryEdge{
				Src: &TestNode{Val: 1},
				Dst: &TestNode{Val: 2},
			}

			// From() and To() should return pointers to objects
			// that implement the Node interface.
			Expect(e.From().ID()).To(BeEquivalentTo(1))
			Expect(e.To().ID()).To(BeEquivalentTo(2))

			ePrime := e.ReversedEdge()

			Expect(ePrime.From().ID()).To(BeEquivalentTo(2))
			Expect(ePrime.To().ID()).To(BeEquivalentTo(1))
		})
	})

	Describe("Weight", func() {
		It("Returns the assigned edge weight", func() {
			e := network.DeliveryEdge{
				Src: &TestNode{Val: int64(1)},
				Dst: &TestNode{Val: int64(2)},
				Wgt: 3.0,
			}

			Expect(e.Weight()).To(BeEquivalentTo(3.0))
		})
	})
})
