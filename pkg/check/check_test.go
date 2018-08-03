package check_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jchesterpivotal/knative-service-resource/pkg"
	kf "k8s.io/client-go/kubernetes/fake"
	sf "github.com/knative/serving/pkg/client/clientset/versioned/fake"
	"github.com/jchesterpivotal/knative-service-resource/pkg/check"
	"github.com/jchesterpivotal/knative-service-resource/pkg/concourse"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
)

var _ = Describe("Check", func() {
	var checker check.Checker
	fakeKubeClient := kf.NewSimpleClientset()
	fakeServiceClient := sf.NewSimpleClientset().ServingV1alpha1().Services("test")
	fakeConfigClient := sf.NewSimpleClientset().ServingV1alpha1().Configurations("test")
	fakeRevClient := sf.NewSimpleClientset().ServingV1alpha1().Revisions("test")

	fakedClients := &clients.Clients{
		Kube:          fakeKubeClient,
		Service:       fakeServiceClient,
		Configuration: fakeConfigClient,
		Revision:      fakeRevClient,
	}
	source := &concourse.Source{
		Name:            "test_name",
		KubernetesUri:   "https://kubernetes.test",
		KubernetesToken: "tokentokentoken",
		KubernetesCa:    "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
	}

	Describe("NewChecker()", func() {
		BeforeEach(func() {
			version := &concourse.Version{ConfigurationGeneration: "999"}

			checker = check.NewChecker(fakedClients, source, version)
		})

		It("Makes a Checker", func() {
			Expect(checker).NotTo(BeNil())
		})
	})

	Describe("Check()", func() {
		Context("Kubernetes version is ahead of Concourse version", func() {
			It("Returns versions that have not been seen by Concourse", func() {
				svc1 := NewServiceWithGeneration(1)
				rev1 := NewRevisionWithGeneration(1)
				rev2 := NewRevisionWithGeneration(2)
				rev3 := NewRevisionWithGeneration(3)

				concourseVersion := &concourse.Version{ConfigurationGeneration: "1"}
				fakedClients.Service = sf.NewSimpleClientset(svc1, rev1).ServingV1alpha1().Services("test")

				svc1.Status.ObservedGeneration = 3
				fakedClients.Service.Update(svc1)
				fakedClients.Revision.Create(rev2)
				fakedClients.Revision.Create(rev3)
				checker = check.NewChecker(fakedClients, source, concourseVersion)

				out, err := checker.Check()

				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(ConsistOf(
					concourse.Version{ConfigurationGeneration: "2"},
					concourse.Version{ConfigurationGeneration: "3"},
				))
			})
		})

		Context("Kubernetes and Concourse have the same version", func() {
			It("Returns the version found in both", func() {
				knativeVersion := NewServiceWithGeneration(22)
				concourseVersion := &concourse.Version{ConfigurationGeneration: "22"}
				fakedClients.Service = sf.NewSimpleClientset(knativeVersion).ServingV1alpha1().Services("test")
				checker = check.NewChecker(fakedClients, source, concourseVersion)

				out, err := checker.Check()

				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(ConsistOf(
					concourse.Version{ConfigurationGeneration: "22"},
				))
			})
		})

		Context("Concourse version is ahead of Kubernetes version", func() {
			It("Returns an error", func() {
				svc1 := NewServiceWithGeneration(1)
				rev1 := NewRevisionWithGeneration(1)
				rev2 := NewRevisionWithGeneration(2)

				concourseVersion := &concourse.Version{ConfigurationGeneration: "3"}
				fakedClients.Service = sf.NewSimpleClientset(svc1, rev1).ServingV1alpha1().Services("test")

				svc1.Status.ObservedGeneration = 2
				fakedClients.Service.Update(svc1)
				fakedClients.Revision.Create(rev2)
				checker = check.NewChecker(fakedClients, source, concourseVersion)

				out, err := checker.Check()

				Expect(err).To(MatchError("version known to Concourse (3) was ahead of version known to Kubernetes (2)"))
				Expect(out).To(BeNil())
			})
		})

		Context("Concourse has a version, but not Kubernetes", func() {
			// error
		})

		Context("Kubernetes has a version, but not Concourse", func() {
			It("Returns the latest version in Kubernetes", func() {
				svc := NewServiceWithGeneration(111)
				fakedClients.Service = sf.NewSimpleClientset(svc).ServingV1alpha1().Services("test")
				checker = check.NewChecker(fakedClients, source, nil)

				out, err := checker.Check()

				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(ConsistOf(
					concourse.Version{ConfigurationGeneration: "111"},
				))
			})
		})

		Context("There is no version in either of Kubernetes or Concourse", func() {
			// error
		})
	})
})

func NewServiceWithGeneration(generation int64) *v1alpha1.Service {
	svc := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "test_name",
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

	svc.Spec.RunLatest.Configuration.RevisionTemplate.Labels = map[string]string{
		"serving.knative.dev/configuration": "test_name",
	}

	svc.Status.ObservedGeneration = generation

	return svc
}

func NewRevisionWithGeneration(generation int64) *v1alpha1.Revision {
	revisionName := fmt.Sprintf("test_name-%04d", generation)

	return &v1alpha1.Revision{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      revisionName,
			Labels: map[string]string{
				"serving.knative.dev/configuration": "test_name",
			},
		},
		Spec: v1alpha1.RevisionSpec{
			Generation: generation,
		},
	}
}
