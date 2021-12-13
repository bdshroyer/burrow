package burrow_test

import (
	. "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/bdshroyer/burrow"
    "github.com/bdshroyer/burrow/matchers"
)

var _ = XDescribe("Nodes", func() {
    Context("DeliveryNodes", func() {
        var (
            hubs []*burrow.HubNode
            stops []*burrow.StopNode
        )

        Context("hybrid input", func() {
            BeforeEach(func() {
                hubs = []*burrow.HubNode{
                    &burrow.HubNode{Val: 4},
                }

                stops = []*burrow.StopNode{
                    &burrow.StopNode{Val: 5},
                    &burrow.StopNode{Val: 3},
                }
            })

            It("Steps through a collection of nodes as an iterator", func() {
                nodeIter := burrow.NewDeliveryNodes(hubs, stops)

                Expect(nodeIter).NotTo(BeNil())
                Expect(nodeIter.Node()).To(matchers.MatchNode(&burrow.HubNode{Val: 4}))
            })
        })
    })
})
