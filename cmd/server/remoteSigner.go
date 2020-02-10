package main

import (
	"context"
	_ "github.com/quan-to/chevron/cmd/server/init"
	"os"
	"os/signal"
	"syscall"

	remote_signer "github.com/quan-to/chevron"
	"github.com/quan-to/chevron/QuantoError"
	"github.com/quan-to/chevron/bootstrap"
	"github.com/quan-to/chevron/etc/magicBuilder"
	"github.com/quan-to/chevron/kubernetes"
	"github.com/quan-to/chevron/server"
	"github.com/quan-to/slog"
)

var log = slog.Scope("QRS").Tag(remote_signer.DefaultTag)

func main() {
	if os.Getenv("SHOW_LINES") == "true" {
		slog.SetShowLines(true)
	}

	QuantoError.EnableStackTrace()

	bootstrap.RunBootstraps()

	ctx := context.Background()
	sm := magicBuilder.MakeSM(log)
	gpg := magicBuilder.MakePGP(log)

	gpg.LoadKeys(ctx)

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
