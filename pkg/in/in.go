package in

import (
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
	"github.com/jchesterpivotal/knative-service-resource/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"

	"fmt"
)

type Inner interface {
	In() (config.InResponse, v1alpha1.Service, v1alpha1.Revision, error)
}

type inner struct {
	clients *clients.Clients

	source  *config.Source
	version *config.Version
}

func (i *inner) In() (config.InResponse, v1alpha1.Service, v1alpha1.Revision, error) {
	svc, err := i.getService()
	if err != nil {
		return config.InResponse{},
		v1alpha1.Service{},
		v1alpha1.Revision{},
		fmt.Errorf("could not find Knative service '%s' in Kubernetes: %s", i.source.Name, err)
	}

	rev, err := i.getRevision()
	if err != nil {
		return config.InResponse{},
			v1alpha1.Service{},
			v1alpha1.Revision{},
			fmt.Errorf("could not find Knative revision for '%s' in Kubernetes: %s", i.source.Name, err)
	}

	output := config.InResponse{
		Version: *i.version,
		Metadata: []config.VersionMetadataField{
			{Name: "kubernetes_cluster_name", Value: svc.ClusterName},
			{Name: "kubernetes_creation_timestamp", Value: svc.CreationTimestamp.String()},
			{Name: "kubernetes_resource_version", Value: svc.ResourceVersion},
			{Name: "kubernetes_uid", Value: string(svc.UID)},
		},
	}

	return output, *svc, *rev, nil
}

func NewInner(clients *clients.Clients, source *config.Source, version *config.Version) Inner {
	return &inner{
		clients: clients,
		source:  source,
		version: version,
	}
}

func (i *inner) getService() (*v1alpha1.Service, error) {
	serviceName := i.source.Name
	service, err := i.clients.Service.Get(serviceName, v1.GetOptions{IncludeUninitialized: false})
	if err != nil {
		return nil, err
	}

	return service.DeepCopy(), nil
}

func (i *inner) getRevision() (*v1alpha1.Revision, error) {
	serviceName := i.source.Name
	revision, err := i.clients.Revision.Get(serviceName, v1.GetOptions{IncludeUninitialized: false})
	if err != nil {
		return nil, err
	}

	return revision.DeepCopy(), nil
}
