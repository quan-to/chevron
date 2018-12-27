package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/SLog"
	"log"
	"net/http"
)

var slog = SLog.Scope("RemoteSigner")

func main() {
	QuantoError.EnableStackTrace()
	sm := remote_signer.MakeSecretsManager()
	gpg := remote_signer.MakePGPManager()

	gpg.LoadKeys()

	ge := remote_signer.MakeGPGEndpoint(sm, gpg)

	// UnlockKey("0016A9CA870AFA59", "I think you will never guess")

	r := mux.NewRouter()
	remote_signer.AddHKPEndpoints(r.PathPrefix("/pks").Subrouter())
	ge.AttachHandlers(r.PathPrefix("/gpg").Subrouter())

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remote_signer.InitHTTPTimer(r)
		remote_signer.CatchAllRouter(w, r, slog)
	})

	listenAddr := fmt.Sprintf("0.0.0.0:%d", remote_signer.HttpPort)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r, // Pass our instance of gorilla/mux in.
	}

	slog.Info("Remote Signer is now listening at %s", listenAddr)

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
