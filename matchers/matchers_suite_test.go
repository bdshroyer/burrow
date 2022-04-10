package matchers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gonum.org/v1/gonum/graph"
)

func TestMatchers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Matchers Suite")
}

/****** graph type stubs used for testing. ******/

type nodeStub1 struct {
	Val int64
}

func (n *nodeStub1) ID() int64 {
	return n.Val
}

type nodeStub2 struct {
	Val int64
}

func (n *nodeStub2) ID() int64 {
	return n.Val
}

type edgeStub1 struct {
	Src graph.Node
	Dst graph.Node
}

func (e *edgeStub1) From() graph.Node {
	return e.Src
}

func (e *edgeStub1) To() graph.Node {
	return e.Dst
}

func (e *edgeStub1) ReversedEdge() graph.Edge {
	return &edgeStub1{Src: e.Dst, Dst: e.Src}
}

type edgeStub2 struct {
	Src graph.Node
	Dst graph.Node
	Wgt float64
}

func (e *edgeStub2) From() graph.Node {
	return e.Src
}

func (e *edgeStub2) To() graph.Node {
	return e.Dst
}

func (e *edgeStub2) ReversedEdge() graph.Edge {
	return &edgeStub2{Src: e.Dst, Dst: e.Src}
}

func (e *edgeStub2) Weight() float64 {
	return e.Wgt
}
