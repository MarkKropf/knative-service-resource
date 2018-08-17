package out_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//. "github.com/onsi/gomega/gbytes"
	//. "github.com/onsi/gomega/gexec"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jchesterpivotal/knative-service-resource/pkg"

	kf "k8s.io/client-go/kubernetes/fake"
	sf "github.com/knative/serving/pkg/client/clientset/versioned/fake"
	"github.com/jchesterpivotal/knative-service-resource/pkg/out"
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
)

var _ = Describe("Out", func() {
	var svc *v1alpha1.Service
	var fakedClients *clients.Clients
	var outer out.Outer
	var response config.OutResponse
	var err error

	Context("Updating the Configuration is successful", func() {
		BeforeEach(func() {
			svc = NewService()

			fakeKubeClient := kf.NewSimpleClientset()
			fakeServiceClient := sf.NewSimpleClientset(svc).ServingV1alpha1().Services("test")
			fakeConfigClient := sf.NewSimpleClientset().ServingV1alpha1().Configurations("test")
			fakeRevClient := sf.NewSimpleClientset().ServingV1alpha1().Revisions("test")

			fakedClients = &clients.Clients{
				Kube:          fakeKubeClient,
				Service:       fakeServiceClient,
				Configuration: fakeConfigClient,
				Revision:      fakeRevClient,
			}

			outer = out.NewOuter(fakedClients, &config.Source{}, &config.PutParams{})
			response, err = outer.Out()
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("injecting fully resolved container images", func() {
			It("recognises SHA digests and injects them into the Configuration", func() {
				// expect fake service client to have received some stuff



				//Expect(outSvc.Name).To(Equal(svc.Name))
			})
			PIt("recognises and resolves tags into SHA digests before injecting into the Configuration", func() {
				// expect fake registry to have received a tag query?
			})
		})

		Describe("injecting the Concourse build", func() {
			PIt("provides a link back to the build that made the deployment")
		})

		Context("When the Service doesn't already exist", func() {
			PIt("Creates a new Service when create_if_not_exists is true")
			PIt("Errors when create_if_not_exists is true")
			PIt("Errors when create_if_not_exists is nil")
		})
	})

	Context("when an image tag is provided but does not exist in the registry", func() {
		PIt("does not deploy")
	})
})

func NewService() *v1alpha1.Service {
	return &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:         "test",
			Name:              "test_name",
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
}

