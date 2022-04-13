package burrow

import (
	"fmt"
	"time"
	"sort"
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

// StopNodeSortable exists to implement the container/heap interface in Golang. Note that the specification includes very little in the way of error checking; users are kind of on the honor system not to do anything that might cause the heap code to panic.
type StopNodeSortable struct {
	Payload []*StopNode
}

// Returns the length of the heap.
func (s StopNodeSortable) Len() int {
	return len(s.Payload)
}

// Returns true if timestamp j is earlier than timestamp i. Imposes an earliest-node-first ordering on the heap.
func (s StopNodeSortable) Less(i, j int) bool {
	return s.Payload[i].Timestamp.Before(s.Payload[j].Timestamp)
}

// Swaps two node positions in the heap.
func (s StopNodeSortable) Swap(i, j int) {
	s.Payload[i], s.Payload[j] = s.Payload[j], s.Payload[i]
}

func SortInPlace(elts []*StopNode) {
	sns := &StopNodeSortable{Payload: elts}
	sort.Sort(sns)
}

// Creates a delivery network with the specified number of hubs and stops  using the provided distribution.
// Returns an error if distro is not a valid sample distribution.
func MakeDeliveryNetwork(nHubNodes, nStopNodes uint, distro SampleDistribution[time.Time]) (*DeliveryNetwork, error) {
	if distro == nil {
		return nil, fmt.Errorf("Must receive a non-null sample distribution.")
	}

	G := &DeliveryNetwork{
		Hubs: make(map[int64]*HubNode, nHubNodes),
		Stops: make(map[int64]*StopNode, nStopNodes),
		DEdges: make(map[int64][]*DeliveryEdge, nHubNodes + nStopNodes),
	}

	nFactory := NewNodeFactory()

	for i := 0; uint(i) < nHubNodes; i++ {
		newHub := nFactory.MakeHub()
		G.Hubs[newHub.ID()] = newHub

		// Allocation hint based on the assumption that most stops are reachable by all hubs
		G.DEdges[newHub.ID()] = make([]*DeliveryEdge, 0, nStopNodes)
	}

	nodeList := make([]*StopNode, 0, nStopNodes)

	// Generate new stop nodes and store them on a sorted min-heap.
	for i := 0; uint(i) < nStopNodes; i++ {
		newStop := nFactory.MakeStop(distro())
		nodeList = append(nodeList, newStop)

		// Allocation hint based on the assumption that most nodes will have an edge leading back to each hub
		G.DEdges[newStop.ID()] = make([]*DeliveryEdge, 0, nHubNodes + (nStopNodes - uint(i) + 1))

		// Add edge nodes linking each hub node to each stop node in both directions.
		for _, hub := range G.Hubs {
			edge := &DeliveryEdge{
				Src: hub,
				Dst: newStop,
				Wgt: float64(1 * time.Hour),
			}

			G.DEdges[hub.ID()] = append(G.DEdges[hub.ID()], edge)
			G.DEdges[newStop.ID()] = append(G.DEdges[newStop.ID()], edge.ReversedEdge().(*DeliveryEdge))
		}
	}

	SortInPlace(nodeList)

	// Extract nodes in order from the heap, connect its predecessors to it, and store it in the graph.
	// Since each node stored in the graph prior to the given node is an earlier stop (due to the min-heap property), a new edge should be drawn from each node in the graph to the new node.
	// The exception to this rule is if two nodes share the exact same timestamp.
	for i := 0; i < len(nodeList); i++ {
		nextStop := nodeList[i]

		for _, prevStop := range G.Stops {
			weight := float64(nextStop.Timestamp.Sub(prevStop.Timestamp))

			if weight > 0.0 {
				edge := &DeliveryEdge{
					Src: prevStop,
					Dst: nextStop,
					Wgt: weight,
				}

				G.DEdges[prevStop.ID()] = append(G.DEdges[prevStop.ID()], edge)
			}
		}

		G.Stops[nextStop.ID()] = nextStop
	}

	return G, nil
}

