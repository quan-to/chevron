package kubernetes

import "time"

type PodStatus struct {
	Phase     string    `json:"phase"`
	HostIP    string    `json:"hostIP"`
	PodIP     string    `json:"podIP"`
	StartTime time.Time `json:"startTime"`
}
