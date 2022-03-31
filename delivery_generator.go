package burrow

import (
	"time"
	"container/heap"
)

// Manages creation variables for delivery nodes.
type NodeFactory struct {
	Counter int64
}

// NewNodeFactory initializes a new stop factory where the counter is set above 0. The start of the counter is not especially significant, but it at least means that node.ID() > 0 is a decent smoke test. 
func NewNodeFactory() *NodeFactory {
	return &NodeFactory{Counter: 1}
}

// Produces a new stop node with timestamp ts. This operation increases the factory's counter.
func (sf *NodeFactory) MakeStop(ts time.Time) *StopNode {
	output := &StopNode{Val: sf.Counter, Timestamp: ts}
	sf.Counter++

	return output
}

// StopNodeHeap exists to implement the container/heap interface in Golang. Note that the specification includes very little in the way of error checking; users are kind of on the honor system not to do anything that might cause the heap code to panic.
type StopNodeHeap struct {
	Payload []*StopNode
}

// Returns the length of the heap.
func (s StopNodeHeap) Len() int {
	return len(s.Payload)
}

// Returns true if timestamp i is earlier than timestamp j. Imposes an earliest-node-first ordering on the heap.
func (s StopNodeHeap) Less(i, j int) bool {
	return s.Payload[i].Timestamp.Before(s.Payload[j].Timestamp)
}

// Swaps two node positions in the heap.
func (s StopNodeHeap) Swap(i, j int) {
	s.Payload[i], s.Payload[j] = s.Payload[j], s.Payload[i]
}

// Pushes a new node onto the heap. Includes a call to Fix(), the container/heap equivalent of heapify.
func (s *StopNodeHeap) Push(v any) {
	(*s).Payload = append((*s).Payload, v.(*StopNode))
	heap.Fix(s, (*s).Len() - 1)
}

// Removes the first node in the heap. According to the given implementation of Less(), this should be the node with the earliest timestamp. Includes a call to Fix(), the container/heap equivalent of heapify.
func (s *StopNodeHeap) Pop() any {
	payload := (*s).Payload[0]
	(*s).Payload = (*s).Payload[1:]
	heap.Fix(s, 0)

	return payload
}

// MakeStopDAG produces a directed acyclic DeliveryNetwork comprised of StopNodes whose timestamps are generated randomly according to the given distribution. Errors out if passed a bad sample distribution.
func MakeStopDAG (newNodeCount uint, distro SampleDistribution[time.Time]) (*DeliveryNetwork, error) {
	G := &DeliveryNetwork{
		Hubs: make(map[int64]*HubNode),
		Stops: make(map[int64]*StopNode),
		DEdges: make(map[int64][]*DeliveryEdge),
	}

	generator, err := NewSampleGenerator(distro)
	if err != nil {
		return nil, err
	}

	nodeHeap := new(StopNodeHeap)
	stops := NewNodeFactory()

	// Generate new stop nodes and store them on a sorted min-heap.
	for sample := range generator.Sample(newNodeCount) {
		nodeHeap.Push(stops.MakeStop(sample))
	}


	// Extract nodes in order from the heap, connect its predecessors to it, and store it in the graph.
	// Since each node stored in the graph prior to the given node is an earlier stop (due to the min-heap property), a new edge should be drawn from each node in the graph to the new node.
	// The exception to this rule is if two nodes share the exact same timestamp.
	for nodeHeap.Len() > 0 {
		newStop := nodeHeap.Pop().(*StopNode)
		for _, prevStop := range G.Stops {
			if prevStop.Timestamp.Before(newStop.Timestamp) {
				edge := &DeliveryEdge{
					Src: prevStop,
					Dst: newStop,
					Wgt: float64(newStop.Timestamp.Sub(prevStop.Timestamp)),
				}
				G.DEdges[prevStop.ID()] = append(G.DEdges[prevStop.ID()], edge)
			}
		}

		G.Stops[newStop.ID()] = newStop
	}

	return G, err
}
