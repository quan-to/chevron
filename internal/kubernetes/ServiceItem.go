package kubernetes

type ServiceItem struct {
	Metadata ItemMetadata    `json:"metadata"`
	Spec     ServiceItemSpec `json:"spec"`
	Status   interface{}     `json:"status"`
}
