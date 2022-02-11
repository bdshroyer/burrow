package burrow

import "gonum.org/v1/gonum/graph"

// DeliveryEdge is a directional edge connecting two DeliveryNodes. It implements the standard gonum Edge interface.
type DeliveryEdge struct {
	Src DeliveryNode
	Dst DeliveryNode
}

// From() returns the edge's source node.
func (e *DeliveryEdge) From() graph.Node {
	return e.Src
}

// To() returns the edge's destination node.
func (e *DeliveryEdge) To() graph.Node {
	return e.Dst
}

// ReversedEdge() returns a DeliveryEdge struct with the same source and destination as the receiver, but reversed.
func (e *DeliveryEdge) ReversedEdge() graph.Edge {
	return &DeliveryEdge{
		Src: e.Dst,
		Dst: e.Src,
	}
}
