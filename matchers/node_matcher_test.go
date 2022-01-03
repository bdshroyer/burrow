package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bdshroyer/burrow/matchers"
)

var _ = Describe("NodeMatcher", func() {
	It("Matches nodes with identical values.", func() {
		expected := &nodeStub1{Val: int64(2)}
		actual := &nodeStub1{Val: int64(2)}
		Expect(actual).To(matchers.MatchNode(expected))
	})

	It("Does not match nodes with different values.", func() {
		expected := &nodeStub1{Val: int64(2)}
		actual := &nodeStub1{Val: int64(3)}
		Expect(actual).NotTo(matchers.MatchNode(expected))
	})

	Context("When the expected input is not a node", func() {
		It("should raise an error on a bad type", func() {
			expected := "oink"
			actual := &nodeStub1{Val: int64(3)}

			match, err := matchers.MatchNode(expected).Match(actual)

			Expect(match).To(BeFalse())
			Expect(err).To(MatchError("NodeMatcher requires a non-nil expected input that implements the Node interface"))
		})

		It("should raise an error on a nil", func() {
			actual := &nodeStub1{Val: int64(3)}

			match, err := matchers.MatchNode(nil).Match(actual)

			Expect(match).To(BeFalse())
			Expect(err).To(MatchError("NodeMatcher requires a non-nil expected input that implements the Node interface"))
		})
	})

	Context("When the actual input is not a node", func() {
		It("should raise an error on a bad type", func() {
			expected := &nodeStub1{Val: int64(2)}
			actual := "quack"

			match, err := matchers.MatchNode(expected).Match(actual)

			Expect(match).To(BeFalse())
			Expect(err).To(MatchError("NodeMatcher requires a non-nil actual input that implements the Node interface"))
		})

		It("Should raise a specific error on a nil", func() {
			expected := &nodeStub1{Val: int64(2)}

			match, err := matchers.MatchNode(expected).Match(nil)

			Expect(match).To(BeFalse())
			Expect(err).To(MatchError("NodeMatcher requires a non-nil actual input that implements the Node interface"))
		})
	})

	// Fails when the actual is unassignable to the expected.
	// This allows us to expect against a DeliveryNode interface.
	Context("When the actual type is incompatible with the expected type", func() {
		It("should fail without raising an error", func() {
			expected := &nodeStub1{Val: int64(2)}
			actual := &nodeStub2{Val: int64(2)}

			matcher := matchers.MatchNode(expected)
			match, err := matcher.Match(actual)

			Expect(match).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
			Expect(matcher.FailureMessage(actual)).To(ContainSubstring("to be assignable to and be an ID match for"))
		})
	})
})
