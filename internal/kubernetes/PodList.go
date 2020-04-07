package kubernetes

type PodList struct {
	Metadata map[string]string `json:"metadata"`
	Items    []Pod             `json:"items"`
}
