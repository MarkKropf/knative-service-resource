package config

type Source struct {
	Name            string `json:"name"`
	KubernetesUri   string `json:"kubernetes_uri"`
	KubernetesToken string `json:"kubernetes_token"`
	KubernetesCa    string `json:"kubernetes_ca"`
}
