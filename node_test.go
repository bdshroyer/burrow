package burrow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow"
)

var _ = Describe("Node", func() {
	Context("HubNode", func() {
		It("Implements Node interface", func() {
			hub := &burrow.HubNode{4}
			Expect(hub.ID()).To(BeEquivalentTo(4))
		})

		It("Identifies as a hub node", func() {
			hub := &burrow.HubNode{4}
			Expect(hub.IsHub()).To(BeTrue())
		})
	})

	Context("StopNode", func() {
		It("Implements Node interface", func() {
			stop := &burrow.StopNode{3}
			Expect(stop.ID()).To(BeEquivalentTo(3))
		})

		It("Does not identify as a hub node", func() {
			stop := &burrow.StopNode{3}
			Expect(stop.IsHub()).To(BeFalse())
		})
	})
})
