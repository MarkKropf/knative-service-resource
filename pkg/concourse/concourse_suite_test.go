package concourse_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConcourseTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Concourse Types Suite")
}
