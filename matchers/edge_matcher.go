package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/types"
	"gonum.org/v1/gonum/graph"
)

type EdgeMatcher struct {
	expected interface{}
}

func MatchEdge(expected interface{}) types.GomegaMatcher {
	return &EdgeMatcher{expected: expected}
}

func coalesceErrors(errorArray ...error) error {
	for _, err := range errorArray {
		if err != nil {
			return err
		}
	}

	return nil
}

func (em *EdgeMatcher) Match(actual interface{}) (success bool, err error) {
	actualEdge, ok := actual.(graph.Edge)
	if !ok {
		return false, fmt.Errorf("EdgeMatcher requires a non-nil actual input that implements the Edge interface")
	}

	expectedEdge, ok := em.expected.(graph.Edge)
	if !ok {
		return false, fmt.Errorf("EdgeMatcher requires a non-nil expected input that implements the Edge interface")
	}

	srcMatch, srcError := MatchNode(expectedEdge.From()).Match(actualEdge.From())
	dstMatch, dstError := MatchNode(expectedEdge.To()).Match(actualEdge.To())

	assignable := reflect.TypeOf(actual).AssignableTo(reflect.TypeOf(em.expected))

	return assignable && srcMatch && dstMatch, coalesceErrors(srcError, dstError)
}

func (em *EdgeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \n\t%#v\nto have node ID and node type matches to \n\t%#v", em.expected, actual)
}

func (em *EdgeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \n\t%#v\nnot to have node ID and node type matches \n\t%#v", em.expected, actual)
}
