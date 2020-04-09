package kubernetes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/quan-to/chevron/internal/config"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const sleepInterval = 1 * 60 * 1000

func KubeRoutine(stopSig chan bool) {
	if !inKubernetes {
		kubeLog.Error("Tried to start KubeRoutine, but not in Kubernetes! Skipping...")
		return
	}

	running := true
	kubeLog.Info("Starting Kubernetes Routine")

	randomWaitTime := rand.Int31n(5)*1000 + 1000 // Milisseconds

	kubeLog.Info("Kubernetes Namespace: %s", Namespace())
	kubeLog.Info("Pod Hostname: %s", Hostname())
	kubeLog.Info("Pod ID: %s", Me().Metadata.UID)
	kubeLog.Info("To avoid concurrency on cluster starting we're waiting 1 second plus some random time")
	kubeLog.Info("The exact time is %d ms", randomWaitTime)

	time.Sleep(time.Millisecond * time.Duration(randomWaitTime))

	for running {
		select {
		case <-stopSig:
			kubeLog.Info("Stopping Kubernetes Routine")
			running = false
		default:
		}

		kubeLog.Info("Checking for other remote-signer nodes...")
		kubeFunc()
		kubeLog.Info("Sleeping for %d ms", sleepInterval)
		time.Sleep(time.Millisecond * sleepInterval)
	}

	kubeLog.Info("Kubernetes Routine Stopped")
}

func kubeFunc() {
	pods := Pods()
	myId := Me().Metadata.UID
	kubeLog.Info("There are %d pods (including me). Fetching encrypted passwords...", len(pods))
	passwordCount := 0
	for _, pod := range pods {
		if pod.Metadata.UID == myId {
			continue
		}
		if pod.Status.Phase != Running {
			continue
		}

		getUrl := fmt.Sprintf("http://%s:%d/remoteSigner/__internal/__getUnlockPasswords", pod.Status.PodIP, config.HttpPort)
		postUrl := fmt.Sprintf("http://localhost:%d/remoteSigner/__internal/__postEncryptedPasswords", config.HttpPort)

		res, err := http.Get(getUrl)
		if err != nil {
			kubeLog.Error("Error fetching unlock passwords from %s: %s", pod.Status.PodIP, err)
			continue
		}

		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			kubeLog.Error("Error fetching unlock passwords from %s: %s", pod.Status.PodIP, err)
			continue
		}

		var passwords map[string]string

		err = json.Unmarshal(data, &passwords)

		if err != nil {
			kubeLog.Error("Error fetching unlock passwords from %s: %s", pod.Status.PodIP, err)
			continue
		}
		passwordCount += len(passwords)
		kubeLog.Info("Received %d passwords from %s", len(passwords), pod.Status.PodIP)
		if len(passwords) > 0 {
			req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(data))

			if err != nil {
				kubeLog.Error("Error sending unlock passwords from %s: %s", pod.Status.PodIP, err)
				continue
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			body, _ := ioutil.ReadAll(resp.Body)
			_ = resp.Body.Close()

			if resp.StatusCode != 200 {
				kubeLog.Error("Error posting passwords: %s", string(body))
				continue
			}
		}
	}
	if passwordCount == 0 {
		kubeLog.Info("No passwords received")
		return
	}

	kubeLog.Info("Received %d passwords from %d pods. Triggering Local Unlock", passwordCount, len(pods))
	_, _ = http.Get(fmt.Sprintf("http://localhost:%d/remoteSigner/__internal/__triggerKeyUnlock", config.HttpPort))
}
