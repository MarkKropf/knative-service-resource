package clients

import (
	"github.com/knative/serving/pkg/client/clientset/versioned"
	servingtyped "github.com/knative/serving/pkg/client/clientset/versioned/typed/serving/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
	"fmt"
)

type Clients struct {
	Kube          kubernetes.Interface
	Service       servingtyped.ServiceInterface
	Configuration servingtyped.ConfigurationInterface
	Revision      servingtyped.RevisionInterface
}

const defaultNamespace = "default"
const clientVersion = "0.0.0"

// NewClients instantiates and returns several clientsets required for making request to the
// Knative Serving cluster specified by the combination of clusterName and configPath. Clients can
// make requests within namespace.
func NewClients(src *config.Source, operation string) (*Clients, error) {
	clients := &Clients{}
	cfg, err := buildClientConfig(src, operation)
	if err != nil {
		return nil, err
	}

	clients.Kube, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	cs, err := versioned.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	clients.Service = cs.ServingV1alpha1().Services(defaultNamespace)
	clients.Configuration = cs.ServingV1alpha1().Configurations(defaultNamespace)
	clients.Revision = cs.ServingV1alpha1().Revisions(defaultNamespace)

	return clients, nil
}

func buildClientConfig(src *config.Source, operation string) (*rest.Config, error) {
	caData := []byte(src.KubernetesCa)
	userAgent := buildUserAgent(operation, clientVersion)

	conf := &rest.Config{
		Host:        src.KubernetesUri,
		BearerToken: src.KubernetesToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: false,
			CAData:   caData,
		},
		UserAgent: userAgent,
	}

	return conf, nil
}

func buildUserAgent(operation string, version string) string {
	return fmt.Sprintf(
		"%s (knative-service-resource; %s; %s)",
		rest.DefaultKubernetesUserAgent(),
		operation,
		version,
	)
}
