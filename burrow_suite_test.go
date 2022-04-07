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

func Today() time.Time {
	payload := time.Now()
	year, month, day := payload.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, payload.Location())
}
