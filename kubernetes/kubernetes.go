package kubernetes

import (
	"encoding/json"
	"fmt"
	"github.com/mewkiz/pkg/osutil"
	"github.com/quan-to/remote-signer/models/KubernetesModels"
	"github.com/quan-to/slog"
	"io/ioutil"
	"os"
	"path"
)

const ServiceAccountPath = "/run/secrets/kubernetes.io/serviceaccount"

var inKubernetes = false
var currentNamespace = ""
var currentHostname = ""
var currentKubeToken = ""
var kubeLog = slog.Scope("Kubernetes")
var me *KubernetesModels.Pod

func serviceURL() string {
	return fmt.Sprintf("https://kubernetes.default.svc/api/v1/namespaces/%s/services", currentNamespace)
}

func podURL() string {
	return fmt.Sprintf("https://kubernetes.default.svc/api/v1/namespaces/%s/pods", currentNamespace)
}

func init() {
	setup()
}

func setup() {
	var err error
	kubeLog.Info("Checking if running in kubernetes")
	inKubernetes = checkInKubernetes()

	if !inKubernetes {
		kubeLog.Info("Not running in kubernetes...")
		return
	}

	kubeLog.Info("In Kubernetes!")
	currentHostname, err = os.Hostname()
	if err != nil {
		kubeLog.Error("Error getting hostname: %s", err)
	}

	currentNamespace = myNamespace()
	currentKubeToken = kubeToken()
	me = MySelf()
}

func checkInKubernetes() bool {
	return osutil.Exists(path.Join(ServiceAccountPath, "token"))
}

func myNamespace() string {
	data, err := ioutil.ReadFile(path.Join(ServiceAccountPath, "namespace"))
	if err != nil {
		kubeLog.Error("Error loading namespace: %s", err)
	}

	return string(data)
}
func kubeToken() string {
	data, err := ioutil.ReadFile(path.Join(ServiceAccountPath, "token"))
	if err != nil {
		kubeLog.Error("Error loading token: %s", err)
	}

	return string(data)
}

func Hostname() string {
	return currentHostname
}

func Namespace() string {
	return currentNamespace
}

func MySelf() *KubernetesModels.Pod {
	data, err := getWithToken(fmt.Sprintf("%s/%s", podURL(), currentHostname), currentKubeToken)

	if err != nil {
		kubeLog.Error("Error fetching Myself: %s", err)
		return nil
	}

	var pod KubernetesModels.Pod

	err = json.Unmarshal([]byte(data), &pod)
	if err != nil {
		kubeLog.Error("Error deserializing: %s", err)
	}

	return &pod
}

func Pods() []KubernetesModels.Pod {
	data, err := getWithToken(podURL(), currentKubeToken)
	if err != nil {
		kubeLog.Error("Error fetching Pods: %s", err)
		return nil
	}

	var podList KubernetesModels.PodList

	err = json.Unmarshal([]byte(data), &podList)
	if err != nil {
		kubeLog.Error("Error deserializing: %s", err)
	}

	return podList.Items
}

func Services() []KubernetesModels.Service {
	data, err := getWithToken(serviceURL(), currentKubeToken)
	if err != nil {
		kubeLog.Error("Error fetching Services: %s", err)
		return nil
	}

	var pods []KubernetesModels.Service

	err = json.Unmarshal([]byte(data), &pods)
	if err != nil {
		kubeLog.Error("Error deserializing: %s", err)
	}

	return pods
}

func Me() KubernetesModels.Pod {
	if me == nil {
		return KubernetesModels.Pod{}
	}

	return *me
}

func InKubernetes() bool {
	return inKubernetes
}
