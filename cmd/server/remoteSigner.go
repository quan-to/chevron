package main

import (
	"context"
	_ "github.com/quan-to/chevron/cmd/server/init"
	"github.com/quan-to/chevron/internal/agent"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/etc/magicbuilder"
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
	var stop chan bool
	var err error

	if os.Getenv("SHOW_LINES") == "true" {
		slog.SetShowLines(true)
	}

	ctx := context.Background()

	dbh, err := agent.MakeDatabaseHandler(log)
	if err != nil {
		slog.Fatal("Error initializing selected database: %s", err)
	}
	ctx = context.WithValue(ctx, tools.CtxDatabaseHandler, dbh)

	sm := magicbuilder.MakeSM(log, dbh)
	gpg := magicbuilder.MakePGP(log, dbh)

	gpg.LoadKeys(ctx)

	if config.SingleKeyMode {
		stop, err = server.RunRemoteSignerServerSingleKey(log, sm, gpg, dbh)
		if err != nil {
			log.Fatal("Error starting in single-key mode: %s", err)
		}
	} else {
		stop = server.RunRemoteSignerServer(log, sm, gpg, dbh)
	}

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
