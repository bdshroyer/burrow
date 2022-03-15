package burrow_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow"
)

const timeFormat string = "2000-01-01 00:00:00"

func dummyStop(id int64) *burrow.StopNode {
	return &burrow.StopNode{
		Val:       id,
		Timestamp: time.Now(),
	}
}

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
			stop := dummyStop(3)
			Expect(stop.ID()).To(BeEquivalentTo(3))
		})

		It("Does not identify as a hub node", func() {
			stop := dummyStop(3)
			Expect(stop.IsHub()).To(BeFalse())
		})
	})
})
