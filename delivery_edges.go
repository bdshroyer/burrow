package burrow

import "gonum.org/v1/gonum/graph"

// An iterator-like structure that implements the graph.Edges and graph.WeightedEdges interfaces.
type DeliveryEdges struct {
	Payload    []*DeliveryEdge
	CurrentIdx int
}

// Constructor that returns a properly initialized DeliveryEdges iterator.
func NewDeliveryEdges() *DeliveryEdges {
	return &DeliveryEdges{
		Payload:    []*DeliveryEdge{},
		CurrentIdx: -1,
	}
}

// Returns the current edge as a graph.Edge, or a nil if the iterator is exhausted.
func (e DeliveryEdges) Edge() graph.Edge {
	if e.CurrentIdx < 0 || e.CurrentIdx >= len(e.Payload) {
		return nil
	}

	return e.Payload[e.CurrentIdx]
}

// Returns the current edge as a graph.WeightedEdge, or a nil if the iterator is exhausted.
func (e DeliveryEdges) WeightedEdge() graph.WeightedEdge {
	if e.CurrentIdx < 0 || e.CurrentIdx >= len(e.Payload) {
		return nil
	}

	return e.Payload[e.CurrentIdx]
}

// Advances the iterator to the next item if any items in the iterator are unexamined.
//
// Next() returns true if there are any unexamined items remaining after the _new_ current item. If there are not, or if the iterator has already been exhausted, it returns false.
func (e *DeliveryEdges) Next() bool {
	if e.CurrentIdx < len(e.Payload) {
		e.CurrentIdx++
		return e.CurrentIdx < len(e.Payload)
	}

	return false
}

// Returns the number of items remaining in the iterator.
func (e DeliveryEdges) Len() int {
	if e.CurrentIdx == -1 {
		return len(e.Payload)
	}

	return len(e.Payload) - e.CurrentIdx
}

// Returns the iterator's focus to the start of the collection.
func (e *DeliveryEdges) Reset() {
	e.CurrentIdx = -1
	return
}

// Returns the current edge pointed to by the iterator. Identical to WeightedEdge, but allows me to create a common interface around gonum's graph iterators.
func (e DeliveryEdges) Current() graph.WeightedEdge {
	return e.WeightedEdge()
}
