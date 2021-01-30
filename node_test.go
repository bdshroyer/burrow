package burrow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow"
)

var _ = Describe("Node", func() {
	Describe("HubNode", func() {
		It("Implements Node interface", func() {
			hub := &burrow.HubNode{4}
			Expect(hub.ID()).To(Equal(int64(4)))
		})

		It("Identifies as a hub node", func() {
			hub := &burrow.HubNode{4}
			Expect(hub.IsHub()).To(BeTrue())
		})
	})

	Describe("StopNode", func() {
		It("Implements Node interface", func() {
			stop := &burrow.StopNode{3}
			Expect(stop.ID()).To(Equal(int64(3)))
		})

		It("Does not identify as a hub node", func() {
			stop := &burrow.StopNode{3}
			Expect(stop.IsHub()).To(BeFalse())
		})
	})
})
