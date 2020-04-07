package main

import (
	"context"
	_ "github.com/quan-to/chevron/cmd/server/init"
	"github.com/quan-to/chevron/internal/bootstrap"
	"github.com/quan-to/chevron/internal/etc/magicBuilder"
	"github.com/quan-to/chevron/internal/kubernetes"
	"github.com/quan-to/chevron/internal/server"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/slog"
	"os"
	"os/signal"
	"syscall"
)

var log = slog.Scope("QRS").Tag(tools.DefaultTag)

func main() {
	if os.Getenv("SHOW_LINES") == "true" {
		slog.SetShowLines(true)
	}

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
