package burrow

import (
	"time"
	"container/heap"
)

type StopFactory struct {
	Counter int64
}

func NewStopFactory() *StopFactory {
	return &StopFactory{Counter: 1}
}

func (sf *StopFactory) MakeStop(ts time.Time) *StopNode {
	output := &StopNode{Val: sf.Counter, Timestamp: ts}
	sf.Counter++

	return output
}


type StopNodeHeap struct {
	Payload []*StopNode
}

func (s StopNodeHeap) Len() int {
	return len(s.Payload)
}

func (s StopNodeHeap) Less(i, j int) bool {
	return s.Payload[i].Timestamp.Before(s.Payload[j].Timestamp)
}

func (s StopNodeHeap) Swap(i, j int) {
	s.Payload[i], s.Payload[j] = s.Payload[j], s.Payload[i]
}

func (s *StopNodeHeap) Push(v any) {
	(*s).Payload = append((*s).Payload, v.(*StopNode))
	heap.Fix(s, (*s).Len() - 1)
}

func (s *StopNodeHeap) Pop() any {
	payload := (*s).Payload[0]
	(*s).Payload = (*s).Payload[1:]
	heap.Fix(s, 0)

	return payload
}

func MakeStopDAG (newNodeCount uint, distro SampleDistribution[time.Time]) (*DeliveryNetwork, error) {
	G := &DeliveryNetwork{
		Hubs: make(map[int64]*HubNode),
		Stops: make(map[int64]*StopNode),
		DEdges: make(map[int64][]*DeliveryEdge),
	}

	stops := NewStopFactory()
	generator, err := NewSampleGenerator(distro)
	if err != nil {
		return nil, err
	}

	nodeHeap := new(StopNodeHeap)

	for sample := range generator.Sample(newNodeCount) {
		nodeHeap.Push(stops.MakeStop(sample))
	}


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
