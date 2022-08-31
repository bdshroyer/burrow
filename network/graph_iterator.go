package network

import "gonum.org/v1/gonum/graph"

type GraphIterator[T any] interface {
	graph.Iterator
	Current() T
}
