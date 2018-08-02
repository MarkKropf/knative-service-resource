package concourse

import "encoding/json"

type VersionMetadata struct {
	KubernetesClusterName       string `json:"kubernetes_cluster_name"`
	KubernetesCreationTimestamp string `json:"kubernetes_creation_timestamp"`
	KubernetesResourceVersion   string `json:"kubernetes_resource_version"`
	KubernetesUid               string `json:"kubernetes_uid"`
}

func VersionMetadataFromInput(rawJson string) (*VersionMetadata, error) {
	verMeta := &VersionMetadata{}

	err := json.Unmarshal([]byte(rawJson), verMeta)
	if err != nil {
		return nil, err
	}

	return verMeta, nil
}