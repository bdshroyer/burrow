package burrow

import "gonum.org/v1/gonum/graph"

type DeliveryEdge struct {
	Src DeliveryNode
	Dst DeliveryNode
}

func (e *DeliveryEdge) From() graph.Node {
	return e.Src
}

func (e *DeliveryEdge) To() graph.Node {
	return e.Dst
}

func (e *DeliveryEdge) ReversedEdge() graph.Edge {
	return &DeliveryEdge{
		Src: e.Dst,
		Dst: e.Src,
	}
}
