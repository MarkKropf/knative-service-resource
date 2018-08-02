package concourse

import (
	"encoding/json"
)

type Source struct {
	Name            string `json:"name"`
	KubernetesUri   string `json:"kubernetes_uri"`
	KubernetesToken string `json:"kubernetes_token"`
	KubernetesCa    string `json:"kubernetes_ca"`
}

func SourceFromInput(rawJson string) (*Source, error) {
	newSource := &Source{}

	err := json.Unmarshal([]byte(rawJson), newSource)
	if err != nil {
		return nil, err
	}

	return newSource, nil
}
