package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"net/http"
)

func GenRemoteSignerServerMux(slog *SLog.Instance, sm *remote_signer.SecretsManager, gpg *remote_signer.PGPManager) *mux.Router {
	ge := MakeGPGEndpoint(sm, gpg)

	r := mux.NewRouter()
	AddHKPEndpoints(r.PathPrefix("/pks").Subrouter())
	ge.AttachHandlers(r.PathPrefix("/gpg").Subrouter())

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InitHTTPTimer(r)
		CatchAllRouter(w, r, slog)
	})

	return r
}

func RunRemoteSignerServer(slog *SLog.Instance, sm *remote_signer.SecretsManager, gpg *remote_signer.PGPManager) chan bool {

	r := GenRemoteSignerServerMux(slog, sm, gpg)

	listenAddr := fmt.Sprintf("0.0.0.0:%d", remote_signer.HttpPort)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r, // Pass our instance of gorilla/mux in.
	}

	stopChannel := make(chan bool)

	go func() {
		<-stopChannel
		slog.Info("Received STOP. Closing server")
		_ = srv.Close()
		stopChannel <- true
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error(err)
		}
		slog.Info("HTTP Server Closed")
	}()

	slog.Info("Remote Signer is now listening at %s", listenAddr)

	return stopChannel
}
