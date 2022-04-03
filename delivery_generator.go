package burrow

import (
	"fmt"
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
func (nf *NodeFactory) MakeStop(ts time.Time) *StopNode {
	output := &StopNode{Val: nf.Counter, Timestamp: ts}
	nf.Counter++

	return output
}

// Produces a new hub node. The operation increases the factory's counter.
func (nf *NodeFactory) MakeHub() *HubNode {
	output := &HubNode{Val: nf.Counter}
	nf.Counter++

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
	if distro == nil {
		return nil, fmt.Errorf("Must receive a non-null sample distribution.")
	}

	G := &DeliveryNetwork{
		Hubs: make(map[int64]*HubNode),
		Stops: make(map[int64]*StopNode),
		DEdges: make(map[int64][]*DeliveryEdge),
	}

	nodeHeap := new(StopNodeHeap)
	stops := NewNodeFactory()

	// Generate new stop nodes and store them on a sorted min-heap.
	for i := 0; uint(i) < newNodeCount; i++ {
		nodeHeap.Push(stops.MakeStop(distro()))
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

	return G, nil
}

func MakeDeliveryNetwork(nHubNodes, nStopNodes uint, distro SampleDistribution[time.Time]) (*DeliveryNetwork, error) {
	if distro == nil {
		return nil, fmt.Errorf("Must receive a non-null sample distribution.")
	}

	G := &DeliveryNetwork{
		Hubs: make(map[int64]*HubNode),
		Stops: make(map[int64]*StopNode),
		DEdges: make(map[int64][]*DeliveryEdge),
	}

	nFactory := NewNodeFactory()

	for i := 0; uint(i) < nHubNodes; i++ {
		newHub := nFactory.MakeHub()
		G.Hubs[newHub.ID()] = newHub
	}

	nodeHeap := new(StopNodeHeap)

	// Generate new stop nodes and store them on a sorted min-heap.
	for i := 0; uint(i) < nStopNodes; i++ {
		newStop := nFactory.MakeStop(distro())
		nodeHeap.Push(nFactory.MakeStop(distro()))

		// Add edge nodes linking each hub node to each stop node in both directions.
		for _, hub := range G.Hubs {
			edge := &DeliveryEdge{
				Src: hub,
				Dst: newStop,
				Wgt: 1.0,
			}

			G.DEdges[hub.ID()] = append(G.DEdges[hub.ID()], edge)
			G.DEdges[newStop.ID()] = append(G.DEdges[newStop.ID()], edge.ReversedEdge().(*DeliveryEdge))
		}
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

	return G, nil
}

