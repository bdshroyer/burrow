package benchmarks_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBenchmarks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Benchmarks Suite")
}
