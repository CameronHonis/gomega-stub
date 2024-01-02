package stub_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGomegaStub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GomegaStub Suite")
}
