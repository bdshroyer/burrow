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

func (nm *NodeMatcher) Match(actual interface{}) (success bool, err error) {
        actualNode, ok := actual.(graph.Node)
        if !ok {
            return false, fmt.Errorf("NodeMatcher requires an actual input that implements the Node interface.")
        }

        expectedNode, ok := nm.expected.(graph.Node)
        if !ok {
            return false, fmt.Errorf("NodeMatcher requires an expected input that implements the Node interface.")
        }


        assignable := reflect.TypeOf(actual).AssignableTo(reflect.TypeOf(nm.expected))
        if !assignable {
            return false, fmt.Errorf("%T cannot be assigned to %T.", actual, nm.expected)
        }

        return expectedNode.ID() == actualNode.ID(), nil
}

func (nm *NodeMatcher) FailureMessage(actual interface{}) (message string) {
    return fmt.Sprintf("Expected \n\t%#v\nto have an ID matching that of \n\t%#v", nm.expected, actual)
}

func (nm *NodeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
    return fmt.Sprintf("Expected \n\t%#v\nnot to have an ID matching that of \n\t%#v", nm.expected, actual)
}


