package burrow_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBurrow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Burrow Suite")
}

func today() time.Time {
	base := time.Now()
	return time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, base.Location())
}
