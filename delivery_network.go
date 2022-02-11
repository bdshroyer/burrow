/*
burrow is a collection of network structures for exploring specific network theory concepts that interest me. In particular, there is a focus on delivery networks, which are DAGs with heterogeneous node types representing the locations visited by delivery vehicles.

This package builds on and makes extensive use of the gonum graph library. For more details on the underlying interfaces, see:

	https://pkg.go.dev/gonum.org/v1/gonum/graph

All node, edge and graph structures in burrow implement the corresponding interfaces in gonum/graph. More specifically:

	* HubNode, StopNode -> gonum/graph.Node
	* DeliveryEdge -> gonum/graph.Edge
	* DeliveryNetwork -> gonum/graph.Graph
*/
package burrow

// DeliveryNodes is the DeliveryNetwork implementation of the Nodes type in gonum/graph. Specifically, it is an iterator that allows application code to traverse the delivery network nodes in a list-like fashion.
//
// Note that the DeliveryNodes iterator traverses hub nodes first, then stop nodes.
type DeliveryNodes struct {
	Payload    []DeliveryNode
	CurrentIdx int
}

// Len() returns the number of nodes covered by the iterator.
func (d DeliveryNodes) Len() int {
	return len(d.Payload)
}

// Node() returns the current node without advancing the iterator; i.e., it works as an implementation of peek.
func (d *DeliveryNodes) Node() DeliveryNode {
	var current DeliveryNode

	if d.Len() > 0 && d.CurrentIdx < d.Len() {
		current = d.Payload[d.CurrentIdx]
	}

	return current
}

/* Next() advances the Nodes iterator to the next node, if one exists

Note that Next() returns true if the iterator has nodes remaining *AFTER* the advance, and false otherwise. For instance, calling Next() to advance to the final node in the iterator will return false, as will subsequent calls to Next().

If the current node is the last node in the iterator, Next() returns false and does not advance.
*/
func (d *DeliveryNodes) Next() bool {
	if d.CurrentIdx < d.Len() {
		d.CurrentIdx++
	}

	return d.CurrentIdx < d.Len()-1
}

// Reset() moves the DeliveryNodes iterator's internal reference back to the start, effectively resetting the iterator.
func (d *DeliveryNodes) Reset() {
	d.CurrentIdx = 0
}

// DeliveryNetwork implements a two-type DAG structure for networks consisting of delivery hubs and stops for vehicles. This network implements the Graph interface from gonum/graph.
//
// The DeliveryNetwork struct stores nodes and edges internally using maps. This decision was made to make accessing structures by index fast and easy; the tradeoff is that it requires a little more work to marshal member structures into collections.
type DeliveryNetwork struct {
	Stops map[int64]*StopNode
	Hubs  map[int64]*HubNode
	Edges map[int64][]*DeliveryEdge
}

// NewDeliveryNetwork() bootstraps a new delivery network from a list of indices corresponding to hubs and stops, as well as the index-level description of the edges connecting them.
func NewDeliveryNetwork(stops []int64, hubs []int64, edges [][2]int64) *DeliveryNetwork {
	G := &DeliveryNetwork{
		Stops: make(map[int64]*StopNode),
		Hubs:  make(map[int64]*HubNode),
		Edges: make(map[int64][]*DeliveryEdge),
	}

	for _, stop_index := range stops {
		G.Stops[stop_index] = &StopNode{Val: stop_index}
	}

	for _, hub_index := range hubs {
		G.Hubs[hub_index] = &HubNode{Val: hub_index}
	}

	for _, edge_pair := range edges {
		src, dst := G.Node(edge_pair[0]), G.Node(edge_pair[1])
		G.Edges[src.ID()] = append(G.Edges[src.ID()], &DeliveryEdge{Src: src, Dst: dst})
	}

	return G
}

// Node(int) returns the node referenced by the given index, or nil if the index can't be found in the network.
//
// Note that this function does not distinguish between hub or stop nodes.
func (G *DeliveryNetwork) Node(id int64) DeliveryNode {
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

// Returns true if an edge connects the two nodes. This function is present to satisfy the Graph interface requirements, but it is NOT direction-sensitive, even though delivery networks are.
func (G *DeliveryNetwork) HasEdgeBetween(xid, yid int64) bool {
	edge_list, ok := G.Edges[xid]
	if ok {
		for _, edge := range edge_list {
			if edge.To().ID() == yid {
				return true
			}
		}
	}

	return false
}

// Returns the edge running from uid to vid, or nil if said edge doesn't exist.
func (G *DeliveryNetwork) Edge(uid, vid int64) *DeliveryEdge {
	edge_list, ok := G.Edges[uid]
	if ok {
		for _, edge := range edge_list {
			if edge.To().ID() == vid {
				return edge
			}
		}
	}

	return nil
}

// Nodes() returns an iterator of type DeliveryNodes, allowing a pass over all the nodes in this network.
func (G *DeliveryNetwork) Nodes() *DeliveryNodes {
	dn := &DeliveryNodes{Payload: make([]DeliveryNode, 0)}

	for _, v := range G.Hubs {
		dn.Payload = append(dn.Payload, v)
	}

	for _, v := range G.Stops {
		dn.Payload = append(dn.Payload, v)
	}

	return dn
}

// From() returns an iterator all nodes for which an edge exists with id as a source. If the specified node has no outbound edges, an empty list is returned.
func (G *DeliveryNetwork) From(id int64) *DeliveryNodes {
	dn := &DeliveryNodes{Payload: make([]DeliveryNode, 0)}

	reachable, ok := G.Edges[id]

	if ok {
		for _, edge := range reachable {
			dn.Payload = append(dn.Payload, edge.To().(DeliveryNode))
		}
	}

	return dn
}
