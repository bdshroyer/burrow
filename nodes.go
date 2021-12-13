package burrow

import (
    "gonum.org/v1/gonum/graph"
)

type DeliveryNodes struct {
    Hubs []*HubNode
    Stops []*StopNode
}

func NewDeliveryNodes(hubs []*HubNode, stops []*StopNode) *DeliveryNodes {
    return &DeliveryNodes{Hubs: hubs, Stops: stops}
}

func (dn *DeliveryNodes) Node() graph.Node {
    return &HubNode{}
}
