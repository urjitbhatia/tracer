package tracer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTracer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tracer Suite")
}
