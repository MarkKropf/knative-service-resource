package config

type Source struct {
	Name            string `json:"name"`
	KubernetesUri   string `json:"kubernetes_uri"`
	KubernetesToken string `json:"kubernetes_token"`
	KubernetesCa    string `json:"kubernetes_ca"`
}

type Version struct {
	ConfigurationGeneration string `json:"configuration_generation"`
}

type VersionMetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PutParams struct {
	ImageRepository string `json:"image_repository"`
	ImageDigest     string `json:"image_digest"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version,omitempty"`
}

type CheckResponse []Version

type InRequest struct {
	Source  Source   `json:"source"`
	Version Version  `json:"version"`
	Params  struct{} `json:"params"`
}

type InResponse struct {
	Version  Version                `json:"version"`
	Metadata []VersionMetadataField `json:"metadata"`
}

