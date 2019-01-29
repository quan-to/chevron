package KubernetesModels

type Pod struct {
	Metadata ItemMetadata `json:"metadata"`
	Status   PodStatus    `json:"status"`
}
