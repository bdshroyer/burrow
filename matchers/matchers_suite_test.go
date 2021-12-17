package matchers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
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

type edgeStub struct {
	Src graph.Node
	Dst graph.Node
}

func (e *edgeStub) From() graph.Node {
	return e.Src
}

func (e *edgeStub) To() graph.Node {
	return e.Dst
}

func (e *edgeStub) ReversedEdge() graph.Edge {
	return &edgeStub{Src: e.Dst, Dst: e.Src}
}
