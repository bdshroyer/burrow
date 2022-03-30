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

// Creates a new DeliveryNodes struct with a properly-initialized index. Note that the Node() and Current() methods will fail on this new struct unless Next() is called first.
func NewDeliveryNodes() *DeliveryNodes {
	return &DeliveryNodes{
		Payload:    []DeliveryNode{},
		CurrentIdx: -1,
	}
}

// Len() returns the number of nodes remaining in the iterator.
func (d DeliveryNodes) Len() int {
	if d.CurrentIdx < 0 {
		return len(d.Payload)
	}

	return len(d.Payload) - d.CurrentIdx
}

// Node() returns the current node without advancing the iterator; i.e., it works as an implementation of peek.
func (d *DeliveryNodes) Node() graph.Node {
	if d.CurrentIdx >= 0 && d.CurrentIdx < len(d.Payload) {
		return d.Payload[d.CurrentIdx]
	}

	return nil
}

/* Next() returns true if there are any nodes remaining in the iterator; it then advances the Nodes iterator to the next node, if one exists.

Note that Next() must be called to _initialize_ the iterator. If this is not done, the Current() and Node() methods will return nil.

If the current node is the last node in the iterator, Next() returns false and does not advance.

This implementation borrows heavily from the [OrderedNodes](https://github.com/gonum/gonum/blob/v0.11.0/graph/iterator/nodes.go) iterator in the Gonum graph library, as it's the most economical way to meet the interface requirements.
*/
func (d *DeliveryNodes) Next() bool {
	if d.CurrentIdx < len(d.Payload) {
		d.CurrentIdx++
		return d.CurrentIdx < len(d.Payload)
	}

	return false
}

// Reset() moves the DeliveryNodes iterator's internal reference back to the start, effectively resetting the iterator.
func (d *DeliveryNodes) Reset() {
	d.CurrentIdx = -1
}

// Current() returns the current node without advancing the iterator. This essentially the same as Node(), but the more generic name allows me to standardize the interface across some other graph iterator types.
func (d *DeliveryNodes) Current() graph.Node {
	return d.Node()
}
