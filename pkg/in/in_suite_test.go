package in_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "In Suite")
}
