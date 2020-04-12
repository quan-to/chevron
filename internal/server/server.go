package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/chevron/internal/agent"
	remote_signer "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/server/pages"
	"github.com/quan-to/chevron/internal/vaultManager"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
	"net/http"
)

// GenRemoteSignerServerMux generates a remote signer HTTP Router
func GenRemoteSignerServerMux(slog slog.Instance, sm interfaces.SMInterface, gpg interfaces.PGPInterface) *mux.Router {
	var vm *vaultManager.VaultManager
	log := slog.Scope("MUX")

	if remote_signer.VaultStorage {
		vm = vaultManager.MakeVaultManager(log, remote_signer.KeyPrefix)
	}

	ge := MakeGPGEndpoint(log, sm, gpg)
	ie := MakeInternalEndpoint(log, sm, gpg)
	te := MakeTestsEndpoint(log, vm)
	kre := MakeKeyRingEndpoint(log, sm, gpg)
	sks := MakeSKSEndpoint(log, sm, gpg)
	tm := agent.MakeTokenManager(log)
	am := agent.MakeAuthManager(log)
	ap := MakeAgentProxy(log, gpg, tm)
	sGql := MakeStaticGraphiQL(log)
	agentAdmin := MakeAgentAdmin(log, tm, am)
	jfc := MakeJFCEndpoint(log, sm, gpg)

	if ge == nil || ie == nil || te == nil || kre == nil || sks == nil || tm == nil || am == nil || ap == nil || agentAdmin == nil {
		slog.Error("One or more services has not been initialized.")
		slog.Error("    GPG Endpoint: %p", ge)
		slog.Error("    Internal Endpoint: %p", ie)
		slog.Error("    Tests Endpoint: %p", te)
		slog.Error("    KeyRing Endpoint: %p", kre)
		slog.Error("    SKS Endpoint: %p", sks)
		slog.Error("    Token Manager: %p", tm)
		slog.Error("    Auth Manager: %p", am)
		slog.Error("    Agent Proxy: %p", ap)
		slog.Error("    Agent Admin: %p", agentAdmin)
		slog.Fatal("Please check if the settings are correct.")
	}

	r := mux.NewRouter()
	// Add for /
	AddHKPEndpoints(log, r.PathPrefix("/pks").Subrouter())
	ge.AttachHandlers(r.PathPrefix("/gpg").Subrouter())
	ie.AttachHandlers(r.PathPrefix("/__internal").Subrouter())
	te.AttachHandlers(r.PathPrefix("/tests").Subrouter())
	kre.AttachHandlers(r.PathPrefix("/keyRing").Subrouter())
	sks.AttachHandlers(r.PathPrefix("/sks").Subrouter())
	jfc.AttachHandlers(r.PathPrefix("/fieldCipher").Subrouter())

	// Add for /remoteSigner
	AddHKPEndpoints(log, r.PathPrefix("/remoteSigner/pks").Subrouter())
	ge.AttachHandlers(r.PathPrefix("/remoteSigner/gpg").Subrouter())
	ie.AttachHandlers(r.PathPrefix("/remoteSigner/__internal").Subrouter())
	te.AttachHandlers(r.PathPrefix("/remoteSigner/tests").Subrouter())
	kre.AttachHandlers(r.PathPrefix("/remoteSigner/keyRing").Subrouter())
	sks.AttachHandlers(r.PathPrefix("/remoteSigner/sks").Subrouter())
	jfc.AttachHandlers(r.PathPrefix("/remoteSigner/fieldCipher").Subrouter())

	// Agent
	ap.AddHandlers(r.PathPrefix("/agent").Subrouter())

	// Agent Admin
	agentAdmin.AddHandlers(r.PathPrefix("/agentAdmin").Subrouter())

	// Static GraphiQL
	sGql.AttachHandlers(r.PathPrefix("/graphiql").Subrouter())

	pages.AddHandlers(r.PathPrefix("/assets").Subrouter())

	// Catch All for unhandled endpoints
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InitHTTPTimer(log, r)
		CatchAllRouter(w, r, log)
	})

	return r
}

// RunRemoteSignerServer runs a remote signer server asynchronously and returns a stop channel
func RunRemoteSignerServer(slog slog.Instance, sm interfaces.SMInterface, gpg interfaces.PGPInterface) chan bool {

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
