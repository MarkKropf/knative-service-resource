package concourse_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jchesterpivotal/knative-service-resource/pkg/concourse"
)

var _ = Describe("VersionMetadata", func() {
	var verMeta *concourse.VersionMetadata
	var err error

	BeforeEach(func() {
		rawJson := `{
          "kubernetes_cluster_name": "test_cluster_name",
          "kubernetes_creation_timestamp": "1234567890",
          "kubernetes_resource_version": "999",
          "kubernetes_uid": "AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE"
        }`

		verMeta, err = concourse.VersionMetadataFromInput(rawJson)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Parsing the JSON input by Concourse", func() {
		It("Extracts the K8s cluster name", func() {
			Expect(verMeta.KubernetesClusterName).To(Equal("test_cluster_name"))
		})
		It("Extracts the K8s creation timestamp", func() {
			Expect(verMeta.KubernetesCreationTimestamp).To(Equal("1234567890"))
		})
		It("Extracts the K8s resource version", func() {
			Expect(verMeta.KubernetesResourceVersion).To(Equal("999"))
		})
		It("Extracts the K8s UID", func() {
			Expect(verMeta.KubernetesUid).To(Equal("AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE"))
		})
	})
})
