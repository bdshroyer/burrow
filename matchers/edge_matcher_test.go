package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gonum.org/v1/gonum/graph"

	"github.com/bdshroyer/burrow/matchers"
)

var _ = Describe("EdgeMatcher", func() {
	var (
		srcExpected, dstExpected graph.Node
		edgeExpected             graph.Edge
	)

	BeforeEach(func() {
		srcExpected = &nodeStub1{Val: int64(2)}
		dstExpected = &nodeStub2{Val: int64(3)}
		edgeExpected = &edgeStub1{Src: srcExpected, Dst: dstExpected}
	})

	It("Matches edges with identical values and node types", func() {
		srcActual := &nodeStub1{Val: int64(2)}
		dstActual := &nodeStub2{Val: int64(3)}
		edgeActual := &edgeStub1{Src: srcActual, Dst: dstActual}

		Expect(edgeActual).To(matchers.MatchEdge(edgeExpected))
	})

	It("Does not match edges with different source or destination values", func() {
		srcActual1 := &nodeStub1{Val: int64(2)}
		dstActual1 := &nodeStub2{Val: int64(4)}
		edgeActual1 := &edgeStub1{Src: srcActual1, Dst: dstActual1}

		srcActual2 := &nodeStub1{Val: int64(4)}
		dstActual2 := &nodeStub2{Val: int64(3)}
		edgeActual2 := &edgeStub1{Src: srcActual2, Dst: dstActual2}

		Expect(edgeActual1).NotTo(matchers.MatchEdge(edgeExpected))
		Expect(edgeActual2).NotTo(matchers.MatchEdge(edgeExpected))
	})

	It("Does not match edges with mismatched node types", func() {
		srcActual1 := &nodeStub1{Val: int64(2)}
		dstActual1 := &nodeStub1{Val: int64(3)}
		edgeActual1 := &edgeStub1{Src: srcActual1, Dst: dstActual1}

		srcActual2 := &nodeStub2{Val: int64(2)}
		dstActual2 := &nodeStub2{Val: int64(3)}
		edgeActual2 := &edgeStub1{Src: srcActual2, Dst: dstActual2}

		matcher := matchers.MatchEdge(edgeExpected)
		match, err := matcher.Match(edgeActual1)

		Expect(match).To(BeFalse())
		Expect(err).NotTo(HaveOccurred())

		Expect(edgeActual1).NotTo(matchers.MatchEdge(edgeExpected))
		Expect(matcher.FailureMessage(edgeActual1)).To(ContainSubstring("to have node ID and node type matches to"))

		match, err = matcher.Match(edgeActual2)
		Expect(edgeActual2).NotTo(matchers.MatchEdge(edgeExpected))
	})

	Context("When the expected input is not a node", func() {
		var (
			srcActual1, dstActual1 graph.Node
			edgeActual1            graph.Edge
		)

		BeforeEach(func() {
			srcActual1 = &nodeStub1{Val: int64(2)}
			dstActual1 = &nodeStub1{Val: int64(3)}
			edgeActual1 = &edgeStub1{Src: srcActual1, Dst: dstActual1}
		})

		It("should raise an error on bad inputs", func() {
			match, err := matchers.MatchEdge(edgeExpected).Match("oink")
			Expect(match).To(BeFalse())
			Expect(err).To(MatchError("EdgeMatcher requires a non-nil actual input that implements the Edge interface"))

			match, err = matchers.MatchEdge("quack").Match(edgeActual1)
			Expect(match).To(BeFalse())
			Expect(err).To(MatchError("EdgeMatcher requires a non-nil expected input that implements the Edge interface"))
		})

		It("should raise an error on nil inputs", func() {
			match, err := matchers.MatchEdge(edgeExpected).Match(nil)
			Expect(match).To(BeFalse())
			Expect(err).To(MatchError("EdgeMatcher requires a non-nil actual input that implements the Edge interface"))

			match, err = matchers.MatchEdge(nil).Match(edgeActual1)
			Expect(match).To(BeFalse())
			Expect(err).To(MatchError("EdgeMatcher requires a non-nil expected input that implements the Edge interface"))
		})
	})

	Context("When the actual type doesn't match the expected type", func() {
		It("should fail without raising an error", func() {
			srcActual1 := &nodeStub1{Val: int64(2)}
			dstActual1 := &nodeStub2{Val: int64(3)}
			edgeActual1 := &edgeStub2{Src: srcActual1, Dst: dstActual1}

			matcher := matchers.MatchEdge(edgeExpected)
			match, err := matcher.Match(edgeActual1)

			Expect(match).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
			Expect(matcher.FailureMessage(edgeActual1)).To(ContainSubstring("to have node ID and node type matches to"))
		})
	})
})
