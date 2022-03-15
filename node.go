package burrow

import (
	"time"

	"gonum.org/v1/gonum/graph"
)

// DeliveryNode an extension of the standard gonum Node interface that encompasses both delivery stops as well as the hubs from which vehicles are dispatched.
type DeliveryNode interface {
	graph.Node
	IsHub() bool
}

// HubNode represents a location from which vehicles are dispatched.
type HubNode struct {
	Val int64
}

// ID() is a Node interface implementer that returns the hub node's ID.
func (n *HubNode) ID() int64 {
	return n.Val
}

// IsHub() returns true for hub nodes.
func (n *HubNode) IsHub() bool {
	return true
}

// StopNode represents a delivery stop made by a vehicle. It is implicitly assumed that stops cannot be hubs.
type StopNode struct {
	Val       int64
	Timestamp time.Time
}

// ID() is a Node interface implementer that returns the stop node's ID.
func (s *StopNode) ID() int64 {
	return s.Val
}

// IsHub() always returns false for stop nodes, as it's assumed for algorithmic purposes that delivery stops cannot also be hubs.
func (s *StopNode) IsHub() bool {
	return false
}
