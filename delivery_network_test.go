package burrow_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gonum.org/v1/gonum/graph"

	"github.com/bdshroyer/burrow"
	"github.com/bdshroyer/burrow/matchers"
)

func MakeTestDeliveryNetwork(stops []int64, hubs []int64, edges [][2]int64) *burrow.DeliveryNetwork {
	G := &burrow.DeliveryNetwork{
		Stops: make(map[int64]*burrow.StopNode),
		Hubs:  make(map[int64]*burrow.HubNode),
		Edges: make(map[int64][]*burrow.DeliveryEdge),
	}

	for _, stop_index := range stops {
		G.Stops[stop_index] = &burrow.StopNode{Val: stop_index}
	}

	for _, hub_index := range hubs {
		G.Hubs[hub_index] = &burrow.HubNode{Val: hub_index}
	}

	for _, edge_pair := range edges {
		src, dst := G.Node(edge_pair[0]), G.Node(edge_pair[1])
		srcDN, dstDN := src.(burrow.DeliveryNode), dst.(burrow.DeliveryNode)

		G.Edges[src.ID()] = append(G.Edges[src.ID()], &burrow.DeliveryEdge{Src: srcDN, Dst: dstDN})
	}

	return G
}

func edgeFromPair(src, dst int) [2]int64 {
	return [2]int64{int64(src), int64(dst)}
}

func hubToStop(src, dst int) *burrow.DeliveryEdge {
	return &burrow.DeliveryEdge{
		Src: &burrow.HubNode{int64(src)},
		Dst: &burrow.StopNode{int64(dst)},
	}
}

func stopToStop(src, dst int) *burrow.DeliveryEdge {
	return &burrow.DeliveryEdge{
		Src: &burrow.StopNode{int64(src)},
		Dst: &burrow.StopNode{int64(dst)},
	}
}

func collectNodes(iter graph.Nodes) []interface{} {
	lst := make([]interface{}, 0, iter.Len())

	if iter.Len() == 0 {
		return lst
	}

	lst = append(lst, iter.Node())

	for {
		proceed := iter.Next()
		lst = append(lst, iter.Node())

		if !proceed {
			break
		}
	}

	return lst
}

var _ = Describe("DeliveryNetwork functionality", func() {
	Describe("NewDeliveryNetwork", func() {
		It("Returns a delivery network containing the specified indices", func() {
			G := burrow.NewDeliveryNetwork(
				[]int64{2, 4},
				[]int64{1},
				[][2]int64{edgeFromPair(2, 4), edgeFromPair(1, 4)},
			)

			Expect(G.Stops).To(HaveLen(2))
			Expect(G.Hubs).To(HaveLen(1))
			Expect(G.Edges).To(HaveLen(2))

			Expect(G.Stops).To(HaveKeyWithValue(
				BeEquivalentTo(2),
				matchers.MatchNode(&burrow.StopNode{Val: 2}),
			))

			Expect(G.Stops).To(HaveKeyWithValue(
				BeEquivalentTo(4),
				matchers.MatchNode(&burrow.StopNode{Val: 4}),
			))

			Expect(G.Hubs).To(HaveKeyWithValue(
				BeEquivalentTo(1),
				matchers.MatchNode(&burrow.HubNode{Val: 1}),
			))

			Expect(G.Edges).To(HaveKeyWithValue(
				BeEquivalentTo(1),
				ContainElements(matchers.MatchEdge(hubToStop(1, 4))),
			))

			Expect(G.Edges).To(HaveKeyWithValue(
				BeEquivalentTo(2),
				ContainElements(matchers.MatchEdge(stopToStop(2, 4))),
			))
		})

		It("Returns an empty network when passed no nodes", func() {
			G := burrow.NewDeliveryNetwork([]int64{}, []int64{}, [][2]int64{})

			Expect(G.Stops).To(HaveLen(0))
			Expect(G.Hubs).To(HaveLen(0))
			Expect(G.Edges).To(HaveLen(0))
		})
	})

	Describe("DeliveryNetwork", func() {
		Context("as a graph", func() {
			var (
				G *burrow.DeliveryNetwork
			)

			BeforeEach(func() {
				G = MakeTestDeliveryNetwork([]int64{4}, []int64{1}, [][2]int64{edgeFromPair(1, 4)})
			})

			It("Implements the required DiGraph interface", func() {
				var _ graph.Directed = (*burrow.DeliveryNetwork)(nil)
				_, ok := interface{}(G).(graph.Directed)
				Expect(ok).To(BeTrue())
			})

			Describe("Node", func() {
				It("Allows nodes to be accessed by ID (Graph interface)", func() {
					// Returns hub nodes and stop nodes seamlessly.
					Expect(G.Node(1).ID()).To(BeEquivalentTo(1))
					Expect(G.Node(4).ID()).To(BeEquivalentTo(4))
				})

				It("Returns nothing if node is not present in graph", func() {
					Expect(G.Node(3)).To(BeNil())
				})
			})

			Describe("HasEdgeBetween", func() {
				It("Indicates whether or not an edge exists between two nodes, irrespective of direction", func() {
					Expect(G.HasEdgeBetween(1, 4)).To(BeTrue())
					Expect(G.HasEdgeBetween(4, 1)).To(BeTrue())
				})

				It("returns false if queried about nodes that do not exist", func() {
					Expect(G.HasEdgeBetween(3, 5)).To(BeFalse())
					Expect(G.HasEdgeBetween(1, 5)).To(BeFalse())
					Expect(G.HasEdgeBetween(3, 4)).To(BeFalse())
				})
			})

			Describe("Edge", func() {
				It("Returns the specified edge if it exists", func() {
					edge := G.Edge(1, 4)
					Expect(edge.From().ID()).To(BeEquivalentTo(1))
					Expect(edge.To().ID()).To(BeEquivalentTo(4))
				})

				It("Returns nil if the given nodes exist but the specified edge doesn't", func() {
					Expect(G.Edge(4, 1)).To(BeNil())
				})

				It("Returns nil if one or more nodes in the edge do not exist", func() {
					Expect(G.Edge(3, 5)).To(BeNil())
					Expect(G.Edge(1, 5)).To(BeNil())
					Expect(G.Edge(3, 4)).To(BeNil())
				})
			})

			Describe("Nodes", func() {
				It("Returns an iterable collection of the graph's nodes", func() {
					G = MakeTestDeliveryNetwork(
						[]int64{4, 3},
						[]int64{1},
						[][2]int64{edgeFromPair(1, 4), edgeFromPair(1, 3)},
					)

					payload := G.Nodes()
					Expect(payload.Len()).To(Equal(3))

					Expect(payload.Node()).To(matchers.MatchNode(&burrow.HubNode{Val: 1}))

					// Test iteration
					Expect(payload.Next()).To(BeTrue())

					middleNode := payload.Node()
					Expect(middleNode).To(Or(
						matchers.MatchNode(&burrow.StopNode{Val: 4}),
						matchers.MatchNode(&burrow.StopNode{Val: 3}),
					))

					Expect(payload.Next()).To(BeFalse())

					lastNode := payload.Node()

					Expect(lastNode).To(Or(
						matchers.MatchNode(&burrow.StopNode{Val: 4}),
						matchers.MatchNode(&burrow.StopNode{Val: 3}),
					))

					Expect(lastNode).NotTo(matchers.MatchNode(middleNode))

					// Iterator returns nil after underlying data is exhausted
					Expect(payload.Next()).To(BeFalse())
					Expect(payload.Node()).To(BeNil())

					//Test reset
					payload.Reset()
					Expect(payload.Node()).To(matchers.MatchNode(&burrow.HubNode{Val: 1}))
				})

				It("Returns an empty collection when the graph is empty", func() {
					G := MakeTestDeliveryNetwork([]int64{}, []int64{}, [][2]int64{})
					Expect(G.Nodes().Len()).To(Equal(0))
					Expect(G.Nodes().Node()).To(BeNil())
				})
			})

			Context("Edge functions", func() {
				var G *burrow.DeliveryNetwork

				BeforeEach(func() {
					G = MakeTestDeliveryNetwork(
						[]int64{4, 3},
						[]int64{1},
						[][2]int64{edgeFromPair(1, 4), edgeFromPair(1, 3), edgeFromPair(3, 4)},
					)
				})

				Describe("From", func() {
					It("Returns an iterator over nodes reachable from the target", func() {
						nodes := collectNodes(G.From(1))
						Expect(len(nodes)).To(Equal(2))
					})

					It("Returns an empty list if target has no outbound edges", func() {
						nodes := collectNodes(G.From(4))
						Expect(nodes).To(BeEmpty())
					})

					It("Returns an empty list if the targeted node does not exist", func() {
						nodes := collectNodes(G.To(2))
						Expect(nodes).To(BeEmpty())
					})
				})

				Describe("To", func() {
					It("Returns an iterable collection of nodes that directly connect to the target", func() {
						nodes := collectNodes(G.To(4))
						Expect(len(nodes)).To(Equal(2))

						Expect(nodes).To(ContainElement(matchers.MatchNode(&burrow.HubNode{Val: 1})))
						Expect(nodes).To(ContainElement(matchers.MatchNode(&burrow.StopNode{Val: 3})))
					})

					It("Returns an empty list if the target has no inbound edges", func() {
						nodes := collectNodes(G.To(1))
						Expect(nodes).To(BeEmpty())
					})

					It("Returns an empty list if the targeted node does not exist", func() {
						nodes := collectNodes(G.To(2))
						Expect(nodes).To(BeEmpty())
					})
				})
			})

			Describe("HasEdgeFromTo", func() {
				var G *burrow.DeliveryNetwork

				BeforeEach(func() {
					G = MakeTestDeliveryNetwork(
						[]int64{4, 3},
						[]int64{1},
						[][2]int64{
							edgeFromPair(1, 4),
							edgeFromPair(1, 3),
							edgeFromPair(3, 4),
							edgeFromPair(3, 1),
						},
					)
				})

				It("returns true if the specified directional edge exists", func() {
					Expect(G.HasEdgeFromTo(1, 4)).To(BeTrue())
				})

				It("returns false if the edge doesn't exist", func() {
					Expect(G.HasEdgeFromTo(4, 1)).To(BeFalse())

					Expect(G.HasEdgeFromTo(1, 2)).To(BeFalse())
					Expect(G.HasEdgeFromTo(2, 1)).To(BeFalse())
					Expect(G.HasEdgeFromTo(2, 5)).To(BeFalse())
				})
			})
		})
	})
})
