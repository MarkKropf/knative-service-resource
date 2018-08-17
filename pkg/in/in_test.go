package in_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jchesterpivotal/knative-service-resource/pkg/in"
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"github.com/jchesterpivotal/knative-service-resource/pkg"
	kf "k8s.io/client-go/kubernetes/fake"
	sf "github.com/knative/serving/pkg/client/clientset/versioned/fake"
	"time"
)

var _ = Describe("In", func() {
	var inner in.Inner

	svc := NewServiceWithMetadata(
		"test_cluster_name",
		metav1.NewTime(time.Unix(0, 0).UTC()),
		"111",
		types.UID("test-uid-string"),
	)

	rev := &v1alpha1.Revision{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name: "test_name",
			UID: "test-revision-uid",
		},
	}

	fakeKubeClient := kf.NewSimpleClientset()
	fakeServiceClient := sf.NewSimpleClientset(svc).ServingV1alpha1().Services("test")
	fakeConfigClient := sf.NewSimpleClientset().ServingV1alpha1().Configurations("test")
	fakeRevClient := sf.NewSimpleClientset(rev).ServingV1alpha1().Revisions("test")

	fakedClients := &clients.Clients{
		Kube:          fakeKubeClient,
		Service:       fakeServiceClient,
		Configuration: fakeConfigClient,
		Revision:      fakeRevClient,
	}

	source := &config.Source{
		Name:            "test_name",
		KubernetesUri:   "https://kubernetes.test",
		KubernetesToken: "tokentokentoken",
		KubernetesCa:    "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
	}

	Context("Kubernetes knows about the requested Service", func() {
		var out config.InResponse
		var outSvc v1alpha1.Service
		var outRev v1alpha1.Revision
		var err error

		BeforeEach(func() {
			inner := in.NewInner(fakedClients, source, &config.Version{ConfigurationGeneration: "111"})
			out, outSvc, outRev, err = inner.In()
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("when returning a version", func() {
			It("Returns the version passed in by Concourse", func() {
				Expect(out.Version).To(Equal(config.Version{ConfigurationGeneration: "111"}))
			})
		})

		Describe("when returning version metadata", func() {
			var metadata []config.VersionMetadataField

			BeforeEach(func() {
				inner = in.NewInner(fakedClients, source, &config.Version{ConfigurationGeneration: "1"})
				out, _, _, err = inner.In()
				Expect(err).NotTo(HaveOccurred())
				metadata = out.Metadata
			})

			It("includes kubernetes_cluster_name", func() {
				Expect(metadata).To(ContainElement(config.VersionMetadataField{Name: "kubernetes_cluster_name", Value: "test_cluster_name"}))
			})
			It("includes kubernetes_creation_timestamp", func() {
				Expect(metadata).To(ContainElement(config.VersionMetadataField{Name: "kubernetes_creation_timestamp", Value: "1970-01-01 00:00:00 +0000 UTC"}))
			})
			It("includes kubernetes_resource_version", func() {
				Expect(metadata).To(ContainElement(config.VersionMetadataField{Name: "kubernetes_resource_version", Value: "111"}))
			})
			It("includes kubernetes_uid", func() {
				Expect(metadata).To(ContainElement(config.VersionMetadataField{Name: "kubernetes_uid", Value: "test-uid-string"}))
			})
		})

		Describe("when returning the Service", func() {
			It("returns the version provided by Kubernetes", func() {
				Expect(outSvc.Name).To(Equal(svc.Name))
				Expect(outSvc.UID).To(Equal(svc.UID))
			})
		})

		Describe("when returning the latest Revision", func() {
			It("returns the version provided by Kubernetes", func() {
				Expect(outRev.Name).To(Equal(rev.Name))
				Expect(outRev.UID).To(Equal(rev.UID))
			})
		})
	})

	Context("Kubernetes does not know about the requested Service", func() {
		It("Returns an error", func() {
			fakedClients.Service = sf.NewSimpleClientset().ServingV1alpha1().Services("test")
			inner := in.NewInner(fakedClients, source, &config.Version{ConfigurationGeneration: "999"})

			out, _, _, err := inner.In()

			Expect(out).To(Equal(config.InResponse{}))
			Expect(err).To(MatchError("could not find Knative service 'test_name' in Kubernetes: services.serving.knative.dev \"test_name\" not found"))
		})
	})
})

func NewServiceWithMetadata(clusterName string, creationTimestamp metav1.Time, resourceVer string, uid types.UID) *v1alpha1.Service {
	svc := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:         "test",
			Name:              "test_name",
			ClusterName:       clusterName,
			CreationTimestamp: creationTimestamp,
			ResourceVersion:   resourceVer,
			UID:               uid,
		},
		Spec: v1alpha1.ServiceSpec{
			RunLatest: &v1alpha1.RunLatestType{
				Configuration: v1alpha1.ConfigurationSpec{
					RevisionTemplate: v1alpha1.RevisionTemplateSpec{
						Spec: v1alpha1.RevisionSpec{
							Container: corev1.Container{
								Image: "https://knative-service-image-registry.test/a-repo-path",
							},
						},
					},
				},
			},
		},
	}

	return svc
}
