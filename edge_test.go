package burrow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow"
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

var _ = Describe("Edge", func() {
	Describe("DeliveryEdge", func() {
		It("Implements the Edge interface", func() {
			e := &burrow.DeliveryEdge{
				Src: &TestNode{Val: int64(1)},
				Dst: &TestNode{Val: int64(2)},
			}

			// From() and To() should return pointers to objects
			// that implement the Node interface.
			Expect(e.From().ID()).To(Equal(int64(1)))
			Expect(e.To().ID()).To(Equal(int64(2)))

			ePrime := e.ReversedEdge()

			Expect(ePrime.From().ID()).To(Equal(int64(2)))
			Expect(ePrime.To().ID()).To(Equal(int64(1)))
		})
	})
})
