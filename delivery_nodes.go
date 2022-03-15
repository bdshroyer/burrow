package burrow

import (
	"gonum.org/v1/gonum/graph"
)

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
func (d *DeliveryNodes) Node() graph.Node {
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