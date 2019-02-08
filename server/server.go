package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/agent"
	"github.com/quan-to/remote-signer/etc"
	"net/http"
)

func GenRemoteSignerServerMux(slog *SLog.Instance, sm etc.SMInterface, gpg etc.PGPInterface) *mux.Router {
	ge := MakeGPGEndpoint(sm, gpg)
	ie := MakeInternalEndpoint(sm, gpg)
	te := MakeTestsEndpoint()
	kre := MakeKeyRingEndpoint(sm, gpg)
	sks := MakeSKSEndpoint(sm, gpg)
	tm := agent.MakeTokenManager()
	am := agent.MakeAuthManager()
	ap := MakeAgentProxy(gpg, tm)
	agentAdmin := MakeAgentAdmin(tm, am)

	r := mux.NewRouter()
	// Add for /
	AddHKPEndpoints(r.PathPrefix("/pks").Subrouter())
	ge.AttachHandlers(r.PathPrefix("/gpg").Subrouter())
	ie.AttachHandlers(r.PathPrefix("/__internal").Subrouter())
	te.AttachHandlers(r.PathPrefix("/tests").Subrouter())
	kre.AttachHandlers(r.PathPrefix("/keyRing").Subrouter())
	sks.AttachHandlers(r.PathPrefix("/sks").Subrouter())

	// Add for /remoteSigner
	AddHKPEndpoints(r.PathPrefix("/remoteSigner/pks").Subrouter())
	ge.AttachHandlers(r.PathPrefix("/remoteSigner/gpg").Subrouter())
	ie.AttachHandlers(r.PathPrefix("/remoteSigner/__internal").Subrouter())
	te.AttachHandlers(r.PathPrefix("/remoteSigner/tests").Subrouter())
	kre.AttachHandlers(r.PathPrefix("/remoteSigner/keyRing").Subrouter())
	sks.AttachHandlers(r.PathPrefix("/remoteSigner/sks").Subrouter())

	// Agent
	ap.AddHandlers(r.PathPrefix("/agent").Subrouter())

	// Agent Admin
	agentAdmin.AddHandlers(r.PathPrefix("/agentAdmin").Subrouter())

	// Catch All for unhandled endpoints
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InitHTTPTimer(r)
		CatchAllRouter(w, r, slog)
	})

	return r
}

func RunRemoteSignerServer(slog *SLog.Instance, sm etc.SMInterface, gpg etc.PGPInterface) chan bool {

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
