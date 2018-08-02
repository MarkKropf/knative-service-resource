package concourse_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jchesterpivotal/knative-service-resource/pkg/concourse"
)

var _ = Describe("Source", func() {
	var src *concourse.Source
	var err error

	BeforeEach(func() {
		rawJson := `{
          "name": "test_name",
          "kubernetes_uri": "https://kubernetes.test",
          "kubernetes_token": "tokentokentokentoken",
          "kubernetes_ca": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----"
        }`

		src, err = concourse.SourceFromInput(rawJson)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Parsing the JSON input by Concourse", func() {
		It("Extracts the name", func() {
			Expect(src.Name).To(Equal("test_name"))
		})
		It("Extracts the Kubernetes API URI", func() {
			Expect(src.KubernetesUri).To(Equal("https://kubernetes.test"))

		})
		It("Extracts the Kubernetes authentication token", func() {
			Expect(src.KubernetesToken).To(Equal("tokentokentokentoken"))

		})
		It("Extracts the Kubernetes CA certificate", func() {
			Expect(src.KubernetesCa).To(Equal("-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----"))
		})
	})
})
