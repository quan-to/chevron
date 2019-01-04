package main

import (
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/server"
	"os"
	"os/signal"
	"syscall"
)

var slog = SLog.Scope("RemoteSigner")

func main() {
	QuantoError.EnableStackTrace()
	sm := remote_signer.MakeSecretsManager()
	gpg := remote_signer.MakePGPManager()

	gpg.LoadKeys()

	stop := server.RunRemoteSignerServer(slog, sm, gpg)
	localStop := make(chan bool)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c               // Wait for SIGTERM (Ctrl + C)
		stop <- true      // Send stop signal to HTTP
		<-stop            // Wait HTTP to Cleanup
		localStop <- true // Send Local stop
	}()

	<-localStop
	slog.Info("Closing Main Routine")
}
