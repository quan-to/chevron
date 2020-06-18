package kubernetes

type ServiceItemSpec struct {
	Ports           []ServiceItemSpecPort `json:"ports"`
	Selector        map[string]string     `json:"selector"`
	ClusterIP       string                `json:"clusterIP"`
	SessionAffinity string                `json:"sessionAffinity"`
}
