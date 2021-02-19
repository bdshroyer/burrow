package matchers_test

import (
    . "github.com/onsi/gomega"
    . "github.com/onsi/ginkgo"

    "github.com/bdshroyer/burrow/matchers"
)


type NodeType1 struct {
	Val int64
}

func (n *NodeType1) ID() int64 {
	return n.Val
}

type NodeType2 struct {
	Val int64
}

func (n *NodeType2) ID() int64 {
	return n.Val
}


var _ = Describe("NodeMatcher", func() {
    It("Matches nodes with identical values.", func() {
        expected := &NodeType1{Val: int64(2)}
        actual := &NodeType1{Val: int64(2)}
        Expect(actual).To(matchers.MatchNode(expected))
    })

    It("Does not match nodes with different values.", func() {
        expected := &NodeType1{Val: int64(2)}
        actual := &NodeType1{Val: int64(3)}
        Expect(actual).NotTo(matchers.MatchNode(expected))
    })

    Context("When the expected input is not a node", func() {
        It("should raise an error", func() {
            expected := "oink"
            actual := &NodeType1{Val: int64(3)}

            match, err := matchers.MatchNode(expected).Match(actual)

            Expect(match).To(BeFalse())
            Expect(err).To(MatchError("NodeMatcher requires an expected input that implements the Node interface."))
        })
    })

    Context("When the actual input is not a node", func() {
        It("should raise an error", func() {
            expected := &NodeType1{Val: int64(2)}
            actual:= "quack"

            match, err := matchers.MatchNode(expected).Match(actual)

            Expect(match).To(BeFalse())
            Expect(err).To(MatchError("NodeMatcher requires an actual input that implements the Node interface."))
        })
    })

    // Fails when the actual is unassignable to the expected, but NOT the other way around.
    // This allows us to expect against a DeliveryNode interface.
    Context("When the actual type doesn't match the expected type", func() {
        It("should raise an error if the actual type can't be assigned to the expected type.", func() {
            expected := &NodeType1{Val: int64(2)}
            actual:= &NodeType2{Val: int64(2)}

            match, err := matchers.MatchNode(expected).Match(actual)

            Expect(match).To(BeFalse())
            Expect(err).To(MatchError("*matchers_test.NodeType2 cannot be assigned to *matchers_test.NodeType1."))
        })
    })
})
