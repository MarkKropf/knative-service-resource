package out

import (
	"github.com/jchesterpivotal/knative-service-resource/pkg"
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
)

type Outer interface {
	Out() (config.OutResponse, error)
}

type outer struct {
	clients *clients.Clients

	source *config.Source
	params *config.PutParams
}

func NewOuter(clients *clients.Clients, source *config.Source, params *config.PutParams) Outer {
	return &outer{
		clients: clients,
		source: source,
		params: params,
	}
}

func (o *outer) Out() (config.OutResponse, error) {
	return config.OutResponse{}, nil
}