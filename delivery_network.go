/*
burrow is a collection of network structures for exploring specific network theory concepts that interest me. In particular, there is a focus on delivery networks, which are DAGs with heterogeneous node types representing the locations visited by delivery vehicles.

This package builds on and makes extensive use of the gonum graph library. For more details on the underlying interfaces, see:

	https://pkg.go.dev/gonum.org/v1/gonum/graph

All node, edge and graph structures in burrow implement the corresponding interfaces in gonum/graph. More specifically:

	* HubNode, StopNode -> gonum/graph.Node
	* DeliveryNodes -> gonum/graph.Nodes
	* DeliveryEdge -> gonum/graph.{Edge, WeightedEdge}
	* DeliveryEdges -> gonum/graph.{Edges, WeightedEdges}
	* DeliveryNetwork -> gonum/graph.{Graph, Directed, Weighted}
*/
package burrow

import (
	"gonum.org/v1/gonum/graph"
)

// DeliveryNetwork implements a two-type DAG structure for networks consisting of delivery hubs and stops for vehicles. This network implements the Graph interface from gonum/graph.
//
// The DeliveryNetwork struct stores nodes and edges internally using maps. This decision was made to make accessing structures by index fast and easy; the tradeoff is that it requires a little more work to marshal member structures into collections.
type DeliveryNetwork struct {
	Stops  map[int64]*StopNode
	Hubs   map[int64]*HubNode
	DEdges map[int64][]*DeliveryEdge
}

// NewDeliveryNetwork() bootstraps a new delivery network from a list of indices corresponding to hubs and stops, as well as the index-level description of the edges connecting them.
func NewDeliveryNetwork(stops []int64, hubs []int64, edges [][2]int64) *DeliveryNetwork {
	G := &DeliveryNetwork{
		Stops:  make(map[int64]*StopNode),
		Hubs:   make(map[int64]*HubNode),
		DEdges: make(map[int64][]*DeliveryEdge),
	}

	for _, stop_index := range stops {
		G.Stops[stop_index] = &StopNode{Val: stop_index}
	}

	for _, hub_index := range hubs {
		G.Hubs[hub_index] = &HubNode{Val: hub_index}
	}

	for _, edge_pair := range edges {
		src, dst := G.Node(edge_pair[0]), G.Node(edge_pair[1])

		G.DEdges[src.ID()] = append(
			G.DEdges[src.ID()],
			&DeliveryEdge{
				Src: src.(DeliveryNode),
				Dst: dst.(DeliveryNode),
			},
		)
	}

	return G
}

// Node(int) returns the node referenced by the given index, or nil if the index can't be found in the network.
//
// Note that this function does not distinguish between hub or stop nodes.
func (G *DeliveryNetwork) Node(id int64) graph.Node {
	var node DeliveryNode

	node, ok := G.Stops[id]
	if ok {
		return node
	}

	node, ok = G.Hubs[id]
	if ok {
		return node
	}

	return nil
}

// HasEdgeBetween returns true if an edge connects the two nodes. This function is present to satisfy the Graph interface requirements, but it is NOT direction-sensitive, even though delivery networks are.
func (G *DeliveryNetwork) HasEdgeBetween(xid, yid int64) bool {
	edgeList, ok := G.DEdges[xid]
	if ok {
		for _, edge := range edgeList {
			if edge.To().ID() == yid {
				return true
			}
		}
	}

	// HasEdgeBetween doesn't care about directionality; check for an edge running in the opposite direction.
	edgeList, ok = G.DEdges[yid]
	if ok {
		for _, edge := range edgeList {
			if edge.To().ID() == xid {
				return true
			}
		}
	}

	return false
}

// HasEdgeFromTo is a directional version of HasEdgeBetween(), returning true if an edge exists with uid as a source and vid as a destination, and false otherwise.
func (G *DeliveryNetwork) HasEdgeFromTo(uid, vid int64) bool {
	edges, ok := G.DEdges[uid]

	if ok {
		for _, edge := range edges {
			if edge.To().ID() == vid {
				return true
			}
		}
	}

	return false
}

// Returns the edge running from uid to vid, or nil if said edge doesn't exist.
func (G *DeliveryNetwork) Edge(uid, vid int64) graph.Edge {
	edgeList, ok := G.DEdges[uid]
	if ok {
		for _, edge := range edgeList {
			if edge.To().ID() == vid {
				return edge
			}
		}
	}

	return nil
}

// Nodes() returns an iterator of type DeliveryNodes, allowing a pass over all the nodes in this network. If the network has no nodes, an empty list is returned.
func (G *DeliveryNetwork) Nodes() graph.Nodes {
	dn := NewDeliveryNodes()

	for _, v := range G.Hubs {
		dn.Payload = append(dn.Payload, v)
	}

	for _, v := range G.Stops {
		dn.Payload = append(dn.Payload, v)
	}

	return dn
}

// From() returns an iterator over all nodes reached by id's outbound edges. If the specified node has no outbound edges, an empty list is returned.
func (G *DeliveryNetwork) From(id int64) graph.Nodes {
	dn := NewDeliveryNodes()

	reachable, ok := G.DEdges[id]

	if ok {
		for _, edge := range reachable {
			dn.Payload = append(dn.Payload, edge.To().(DeliveryNode))
		}
	}

	return dn
}

// Returns an iterator over all nodes with a direct hop to the node specified by id. If the specified node has no inbound edges, an empty list is returned.
func (G *DeliveryNetwork) To(id int64) graph.Nodes {
	dn := NewDeliveryNodes()

	dst := G.Node(id)

	if dst == nil {
		return dn
	}

	for _, src := range G.Hubs {
		srcEdges, ok := G.DEdges[src.ID()]

		if ok {
			for _, e := range srcEdges {
				if e.To().ID() == dst.ID() {
					dn.Payload = append(dn.Payload, e.From().(DeliveryNode))
				}
			}
		}
	}

	for _, src := range G.Stops {
		srcEdges, ok := G.DEdges[src.ID()]

		if ok {
			for _, e := range srcEdges {
				if e.To().ID() == dst.ID() {
					dn.Payload = append(dn.Payload, e.From().(DeliveryNode))
				}
			}
		}
	}

	return dn
}

// Returns the weighted edge specified by the two vertex IDs uid, vid. Returns nil if no such edge exists.
func (G *DeliveryNetwork) WeightedEdge(uid, vid int64) graph.WeightedEdge {
	edges, ok := G.DEdges[uid]
	if !ok {
		return nil
	}

	for _, e := range edges {
		if e.To().ID() == vid {
			return G.Edge(uid, vid).(graph.WeightedEdge)
		}
	}

	return nil
}

// Returns an iterator over all the edges in the network.
func (G *DeliveryNetwork) Edges() graph.WeightedEdges {
	edges := NewDeliveryEdges()

	for _, es := range G.DEdges {
		edges.Payload = append(edges.Payload, es...)
	}

	return edges
}

// Returns the weight of the edge specified, as well a hash-style ok variable. If no edge exists between the specified vertices, then it returns a 0 value for the edge weight, as well as a success value of false.
//
// Note that Weight()'s ok return value will be set to true if uid == vid, even though the weight return will be the default option.
func (G *DeliveryNetwork) Weight(uid, vid int64) (float64, bool) {
	if uid == vid {
		return 0.0, true
	}

	eList, ok := G.DEdges[uid]
	if !ok {
		return 0.0, false
	}

	for _, e := range eList {
		if e.To().ID() == vid {
			return e.Weight(), true
		}
	}

	return 0.0, false
}
