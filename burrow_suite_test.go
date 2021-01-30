package burrow_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBurrow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Burrow Suite")
}
