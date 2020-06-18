package kubernetes

type Pod struct {
	Metadata ItemMetadata `json:"metadata"`
	Status   PodStatus    `json:"status"`
}
