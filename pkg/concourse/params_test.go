package concourse_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jchesterpivotal/knative-service-resource/pkg/concourse"
)

var _ = Describe("PutParams", func() {
	var params *concourse.PutParams
	var err error

	BeforeEach(func() {
		rawJson := `{
          "image_repository": "https://registry.test/repositorypath",
          "image_digest": "abc123def456"
        }`

		params, err = concourse.ParamsFromInput(rawJson)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Parsing the JSON input by Concourse", func() {
		It("Extracts the image repository", func() {
			Expect(params.ImageRepository).To(Equal("https://registry.test/repositorypath"))
		})
		It("Extracts the image digest", func() {
			Expect(params.ImageDigest).To(Equal("abc123def456"))
		})
	})
})
