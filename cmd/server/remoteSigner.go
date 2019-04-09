package main

import (
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/bootstrap"
	"github.com/quan-to/remote-signer/etc/magicBuilder"
	"github.com/quan-to/remote-signer/kubernetes"
	"github.com/quan-to/remote-signer/server"
	"github.com/quan-to/slog"
	"os"
	"os/signal"
	"syscall"
)

var log = slog.Scope("RemoteSigner")

func main() {
	QuantoError.EnableStackTrace()

	bootstrap.RunBootstraps()

	sm := magicBuilder.MakeSM()
	gpg := magicBuilder.MakePGP()

	gpg.LoadKeys()

	stop := server.RunRemoteSignerServer(log, sm, gpg)
	localStop := make(chan bool)
	kubeStop := make(chan bool)

	if kubernetes.InKubernetes() {
		go kubernetes.KubeRoutine(kubeStop)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c // Wait for SIGTERM (Ctrl + C)
		if kubernetes.InKubernetes() {
			kubeStop <- true // Send Stop signal to Kubernetes Routine
		}
		stop <- true      // Send stop signal to HTTP
		<-stop            // Wait HTTP to Cleanup
		localStop <- true // Send Local stop
	}()

	<-localStop
	log.Info("Closing Main Routine")
}
