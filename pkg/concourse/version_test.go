package concourse_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jchesterpivotal/knative-service-resource/pkg/concourse"
)

var _ = Describe("Version", func() {
	var version *concourse.Version
	var err error

	BeforeEach(func() {
		rawJson := `{
          "configuration_generation": "999"
        }`

		version, err = concourse.VersionFromInput(rawJson)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Parsing the JSON input by Concourse", func() {
		It("Extracts the configuration generation", func() {
			Expect(version.ConfigurationGeneration).To(Equal("999"))
		})
	})
})
