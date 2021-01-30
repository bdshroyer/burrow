package burrow

import "gonum.org/v1/gonum/graph"

type DeliveryNode interface {
	graph.Node
	IsHub() bool
}

type HubNode struct {
	Val int64
}

func (n *HubNode) ID() int64 {
	return n.Val
}

func (n *HubNode) IsHub() bool {
	return true
}

type StopNode struct {
	Val int64
}

func (s *StopNode) ID() int64 {
	return s.Val
}

func (s *StopNode) IsHub() bool {
	return false
}
