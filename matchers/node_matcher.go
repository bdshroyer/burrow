package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/types"
	"gonum.org/v1/gonum/graph"
)

type NodeMatcher struct {
	expected interface{}
}

func MatchNode(expected interface{}) types.GomegaMatcher {
	return &NodeMatcher{expected: expected}
}

func (nm *NodeMatcher) assignable(actual interface{}) bool {
	return reflect.TypeOf(actual).AssignableTo(reflect.TypeOf(nm.expected))
}

func (nm *NodeMatcher) Match(actual interface{}) (success bool, err error) {
	actualNode, ok := actual.(graph.Node)
	if !ok || actual == nil {
		return false, fmt.Errorf("NodeMatcher requires a non-nil actual input that implements the Node interface")
	}

	expectedNode, ok := nm.expected.(graph.Node)
	if !ok || nm.expected == nil {
		return false, fmt.Errorf("NodeMatcher requires a non-nil expected input that implements the Node interface")
	}

	return nm.assignable(actual) && expectedNode.ID() == actualNode.ID(), nil
}

func (nm *NodeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \n\t%#v\nto be assignable to and be an ID match for %#v.\n", actual, nm.expected)
}

func (nm *NodeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \n\t%#v\nto be unassignable to or be an ID mismatch for %#v.\n", actual, nm.expected)
}
