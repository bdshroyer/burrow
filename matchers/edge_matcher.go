package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega"
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

func coalesceBooleans(boolArray ...bool) bool {
	for _, b := range boolArray {
		if !b {
			return false
		}
	}

	return true
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

	matches := make([]bool, 0, 3)
	errors := make([]error, 0, 2)

	match, err := MatchNode(expectedEdge.From()).Match(actualEdge.From())
	matches = append(matches, match)
	errors = append(errors, err)

	match, err = MatchNode(expectedEdge.To()).Match(actualEdge.To())
	matches = append(matches, match)
	errors = append(errors, err)

	assignable := reflect.TypeOf(actual).AssignableTo(reflect.TypeOf(em.expected))
	matches = append(matches, assignable)

	if assignable {
		expectedWeighted, expectOK := expectedEdge.(graph.WeightedEdge)
		actualWeighted, actualOK := actualEdge.(graph.WeightedEdge)

		if expectOK && actualOK {
			match, err = gomega.BeNumerically("~", expectedWeighted.Weight()).Match(actualWeighted.Weight())
			matches = append(matches, match)
			errors = append(errors, err)
		}
	}

	return coalesceBooleans(matches...), coalesceErrors(errors...)
}

func (em *EdgeMatcher) FailureMessage(actual interface{}) (message string) {
	errMsg := fmt.Sprintf("Expected \n\t%#v\nto have node ID and node type matches to \n\t%#v", em.expected, actual)

	_, weighted := em.expected.(graph.WeightedEdge)
	if weighted {
		errMsg = fmt.Sprintf("Expected \n\t%#v\nto have node ID, node type, and weight matches to \n\t%#v", em.expected, actual)
	}

	return errMsg
}

func (em *EdgeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	errMsg := fmt.Sprintf("Expected \n\t%#v\nnot to have node ID and node type matches \n\t%#v", em.expected, actual)
	_, weighted := em.expected.(graph.WeightedEdge)
	if weighted {
		errMsg = fmt.Sprintf("Expected \n\t%#v\nnot to have node ID, node type and weight matches to \n\t%#v", em.expected, actual)
	}

	return errMsg
}
