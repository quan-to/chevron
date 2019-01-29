package KubernetesModels

type Service struct {
	Kind       string        `json:"kind"`
	APIVersion string        `json:"apiVersion"`
	Metadata   interface{}   `json:"metadata"`
	Items      []ServiceItem `json:"items"`
}
