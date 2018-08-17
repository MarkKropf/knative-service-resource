package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
	"os"

	"github.com/onsi/gomega/ghttp"
	"net/http"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"os/exec"
	"bytes"
	"encoding/json"
	"io/ioutil"
)

var _ = Describe("Out", func() {
	var pathToOut string
	var destDir string
	var server *ghttp.Server
	var response config.OutResponse
	var session *Session
	var err error

	BeforeSuite(func() {
		pathToOut, err = Build("github.com/jchesterpivotal/knative-service-resource/cmd/out")
		Expect(err).NotTo(HaveOccurred())

		file, err := os.Open(pathToOut)
		Expect(err).NotTo(HaveOccurred())
		file.Chmod(os.FileMode(os.ModePerm))

		destDir, err = ioutil.TempDir("", "out-dir")
		Expect(err).ToNot(HaveOccurred())

		server = ghttp.NewServer()
		server.RouteToHandler("POST",
			"/apis/serving.knative.dev/v1alpha1/namespaces/default/services/test_name",
			ghttp.RespondWithJSONEncoded(http.StatusOK, NewService()),
		)
	})

	AfterSuite(func() {
		Expect(os.RemoveAll(destDir)).To(Succeed())
		CleanupBuildArtifacts()
		server.Close()
	})

	Context("Updating the Configuration is successful", func() {
		BeforeEach(func() {
			payload, err := json.Marshal(config.OutRequest{
				Source: config.Source{
					Name:            "test_name",
					KubernetesUri:   server.URL(),
					KubernetesToken: "token",
					KubernetesCa:    "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
				},
				Params:  config.PutParams{},
			})
			Expect(err).ToNot(HaveOccurred())

			cmd := exec.Command(pathToOut, destDir)
			cmd.Stdin = bytes.NewBuffer(payload)
			session, err = Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Returns the newly-created version", func() {
			Expect(response.Version.ConfigurationGeneration).To(Equal("111"))
		})

		It("Returns metadata", func() {
			Expect(response.Metadata).To(ContainElement(config.VersionMetadataField{Name: "kubernetes_uid", Value: "test-uid-value"}))
		})
	})

	Context("Something goes wrong while updating the Configuration", func() {
		BeforeEach(func() {
			payload, err := json.Marshal(config.OutRequest{
				Source: config.Source{
					Name:            "test_name",
					KubernetesUri:   server.URL(),
					KubernetesToken: "token",
					KubernetesCa:    "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
				},
				Params:  config.PutParams{},
			})
			Expect(err).ToNot(HaveOccurred())

			cmd := exec.Command(pathToOut)
			cmd.Stdin = bytes.NewBuffer(payload)
			session, err = Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Prints the error to stderr", func() {
			Eventually(session.Err).Should(Say("flargle"))
		})

		It("Exits with code 1", func() {
			Eventually(session).Should(Exit(1))
		})
	})
})

func NewService() *v1alpha1.Service {
	svc := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test_name",
			UID:       "test-uid-value",
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
