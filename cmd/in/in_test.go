package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"bytes"
	"os/exec"
	"encoding/json"
	"github.com/onsi/gomega/gexec"
	"path/filepath"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
	"github.com/onsi/gomega/ghttp"
	"net/http"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"gopkg.in/yaml.v2"
)

var _ = Describe("In", func() {
	var destDir string
	var pathToIn string
	var err error
	var server *ghttp.Server

	BeforeEach(func() {
		pathToIn, err = gexec.Build("github.com/jchesterpivotal/knative-service-resource/cmd/in")
		Expect(err).NotTo(HaveOccurred())

		file, err := os.Open(pathToIn)
		Expect(err).NotTo(HaveOccurred())
		file.Chmod(os.FileMode(os.ModePerm))

		destDir, err = ioutil.TempDir("", "in-dir")
		Expect(err).ToNot(HaveOccurred())

		server = ghttp.NewServer()
		server.RouteToHandler("GET",
			"/apis/serving.knative.dev/v1alpha1/namespaces/default/revisions/test_name",
			ghttp.RespondWithJSONEncoded(200, NewRevision()),
		)
		server.RouteToHandler("GET",
			"/apis/serving.knative.dev/v1alpha1/namespaces/default/services/test_name",
			ghttp.RespondWithJSONEncoded(http.StatusOK, NewService()),
		)
	})

	AfterEach(func() {
		Expect(os.RemoveAll(destDir)).To(Succeed())
		gexec.CleanupBuildArtifacts()
		server.Close()
	})

	JustBeforeEach(func() {
		cmd := exec.Command(pathToIn, destDir)

		payload, err := json.Marshal(config.InRequest{
			Source: config.Source{
				Name:            "test_name",
				KubernetesUri:   server.URL(),
				KubernetesToken: "token",
				KubernetesCa:    "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
			},
			Version: config.Version{ ConfigurationGeneration: "111" },
			Params:  struct{}{},
		})
		Expect(err).ToNot(HaveOccurred())

		outBuf := new(bytes.Buffer)

		cmd.Stdin = bytes.NewBuffer(payload)
		cmd.Stdout = outBuf
		cmd.Stderr = GinkgoWriter

		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred())

		//err = json.Unmarshal(outBuf.Bytes(), &in.InResponse{})
		//Expect(err).ToNot(HaveOccurred())
	})

	It("Writes service.json", func() {
		svFile, err := os.Open(filepath.Join(destDir, "service.json"))
		Expect(err).NotTo(HaveOccurred())

		svc := &v1alpha1.Service{}
		err = json.NewDecoder(svFile).Decode(svc)
		Expect(err).NotTo(HaveOccurred())

		Expect(svc.Name).To(Equal("test_name"))
		Expect(svc.Spec.RunLatest.Configuration.RevisionTemplate.Spec.Container.Image).To(Equal("https://knative-service-image-registry.test/a-repo-path"))
	})

	It("Writes service.yaml", func() {
		svFile, err := os.Open(filepath.Join(destDir, "service.yaml"))
		Expect(err).NotTo(HaveOccurred())

		svc := &v1alpha1.Service{}
		err = yaml.NewDecoder(svFile).Decode(svc)
		Expect(err).NotTo(HaveOccurred())

		Expect(svc.Name).To(Equal("test_name"))
		Expect(svc.Spec.RunLatest.Configuration.RevisionTemplate.Spec.Container.Image).To(Equal("https://knative-service-image-registry.test/a-repo-path"))
	})

	It("Writes revision/latest.json", func() {
		rvFile, err := os.Open(filepath.Join(destDir, "revision", "latest.json"))
		Expect(err).NotTo(HaveOccurred())

		rev := &v1alpha1.Revision{}
		err = json.NewDecoder(rvFile).Decode(rev)
		Expect(err).NotTo(HaveOccurred())

		Expect(rev.Name).To(Equal("test_name"))
		Expect(rev.Spec.ServiceAccountName).To(Equal("a test value"))
	})

	It("Writes revision/latest.yaml", func() {
		rvFile, err := os.Open(filepath.Join(destDir, "revision", "latest.yaml"))
		Expect(err).NotTo(HaveOccurred())

		rev := &v1alpha1.Revision{}
		err = yaml.NewDecoder(rvFile).Decode(rev)
		Expect(err).NotTo(HaveOccurred())

		Expect(rev.Name).To(Equal("test_name"))
		Expect(rev.Spec.ServiceAccountName).To(Equal("a test value"))
	})

	//It("Returns the version", func() {})
	//It("Returns metadata", func() {})

	//Context("Something goes wrong while getting the Service and Revision", func() {
	//	It("Prints the error to stdout", func() {})
	//	It("Exits with code 1", func() {})
	//})
})

func NewService() *v1alpha1.Service {
	svc := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:         "default",
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

	return svc
}

func NewRevision() *v1alpha1.Revision {
	return &v1alpha1.Revision{
		TypeMeta: metav1.TypeMeta{
			Kind: "Revision",
			APIVersion: "serving.knative.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test_name",
			Namespace: "test",
		},
		Spec: v1alpha1.RevisionSpec{
			ServiceAccountName: "a test value",
		},
	}
}
