package burrow

type DeliveryNodes struct {
	Payload    []DeliveryNode
	CurrentIdx int
}

func (d DeliveryNodes) Len() int {
	return len(d.Payload)
}

func (d *DeliveryNodes) Node() DeliveryNode {
	var current DeliveryNode

	if d.Len() > 0 && d.CurrentIdx < d.Len() {
		current = d.Payload[d.CurrentIdx]
	}

	return current
}

func (d *DeliveryNodes) Next() bool {
	if d.CurrentIdx < d.Len() {
		d.CurrentIdx++
	}

	return d.CurrentIdx < d.Len()-1
}

func (d *DeliveryNodes) Reset() {
	d.CurrentIdx = 0
}

type DeliveryNetwork struct {
	Stops map[int64]*StopNode
	Hubs  map[int64]*HubNode
	Edges map[int64][]*DeliveryEdge
}

func NewDeliveryNetwork(stops []int64, hubs []int64, edges [][2]int64) *DeliveryNetwork {
	G := &DeliveryNetwork{
		Stops: make(map[int64]*StopNode),
		Hubs:  make(map[int64]*HubNode),
		Edges: make(map[int64][]*DeliveryEdge),
	}

	for _, stop_index := range stops {
		G.Stops[stop_index] = &StopNode{Val: stop_index}
	}

	for _, hub_index := range hubs {
		G.Hubs[hub_index] = &HubNode{Val: hub_index}
	}

	for _, edge_pair := range edges {
		src, dst := G.Node(edge_pair[0]), G.Node(edge_pair[1])
		G.Edges[src.ID()] = append(G.Edges[src.ID()], &DeliveryEdge{Src: src, Dst: dst})
	}

	return G
}

func (G *DeliveryNetwork) Node(id int64) DeliveryNode {
	var node DeliveryNode

	node, ok := G.Stops[id]
	if ok {
		return node
	}

	node, ok = G.Hubs[id]
	if ok {
		return node
	}

	return nil
}

func (G *DeliveryNetwork) HasEdgeBetween(xid, yid int64) bool {
	edge_list, ok := G.Edges[xid]
	if ok {
		for _, edge := range edge_list {
			if edge.To().ID() == yid {
				return true
			}
		}
	}

	return false
}

func (G *DeliveryNetwork) Edge(uid, vid int64) *DeliveryEdge {
	edge_list, ok := G.Edges[uid]
	if ok {
		for _, edge := range edge_list {
			if edge.To().ID() == vid {
				return edge
			}
		}
	}

	return nil
}

func (G *DeliveryNetwork) Nodes() *DeliveryNodes {
	dn := &DeliveryNodes{Payload: make([]DeliveryNode, 0)}

	for _, v := range G.Hubs {
		dn.Payload = append(dn.Payload, v)
	}

	for _, v := range G.Stops {
		dn.Payload = append(dn.Payload, v)
	}

	return dn
}

func (G *DeliveryNetwork) From(id int64) *DeliveryNodes {
	dn := &DeliveryNodes{Payload: make([]DeliveryNode, 0)}

	reachable, ok := G.Edges[id]

	if ok {
		for _, edge := range reachable {
			dn.Payload = append(dn.Payload, edge.To().(DeliveryNode))
		}
	}

	return dn
}
