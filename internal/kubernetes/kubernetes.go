package kubernetes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mewkiz/pkg/osutil"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/slog"
)

const ServiceAccountPath = "/run/secrets/kubernetes.io/serviceaccount"

var inKubernetes = false
var currentNamespace = ""
var currentHostname = ""
var currentKubeToken = ""
var kubeLog = slog.Scope("Kubernetes").Tag(tools.DefaultTag)
var me *Pod

func podURL() string {
	return fmt.Sprintf("https://kubernetes.default.svc/api/v1/namespaces/%s/pods", currentNamespace)
}

func init() {
	setup()
}

func setup() {
	var err error
	kubeLog.Await("Checking if running in kubernetes")
	inKubernetes = checkInKubernetes()

	if !inKubernetes {
		kubeLog.Done("Not running in kubernetes...")
		return
	}

	kubeLog.Done("In Kubernetes!")
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

// Hostname returns the hostname of the kubernetes pod
func Hostname() string {
	return currentHostname
}

// Namespace returns the namespace of the kubernetes pod
func Namespace() string {
	return currentNamespace
}

// MySelf returns the data of the kubernetes pod
func MySelf() *Pod {
	data, err := getWithToken(fmt.Sprintf("%s/%s", podURL(), currentHostname), currentKubeToken)

	if err != nil {
		kubeLog.Error("Error fetching Myself: %s", err)
		return nil
	}

	var pod Pod

	err = json.Unmarshal([]byte(data), &pod)
	if err != nil {
		kubeLog.Error("Error deserializing: %s", err)
	}

	return &pod
}

// Pods returns a list of the pods in the current deployment
func Pods() []Pod {
	data, err := getWithToken(podURL(), currentKubeToken)
	if err != nil {
		kubeLog.Error("Error fetching Pods: %s", err)
		return nil
	}

	var podList PodList

	err = json.Unmarshal([]byte(data), &podList)
	if err != nil {
		kubeLog.Error("Error deserializing: %s", err)
	}

	return podList.Items
}

// Me returns the pod that is currently running this instance
func Me() Pod {
	if me == nil {
		return Pod{}
	}

	return *me
}

// InKubernetes returns true when running inside a kubernetes instance
func InKubernetes() bool {
	return inKubernetes
}
