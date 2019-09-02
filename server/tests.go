package server

import (
	"github.com/gorilla/mux"
	remote_signer "github.com/quan-to/chevron"
	"github.com/quan-to/chevron/database"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/vaultManager"
	"net/http"
)

type TestsEndpoint struct{}

func MakeTestsEndpoint() *TestsEndpoint {
	return &TestsEndpoint{}
}

func (ge *TestsEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/ping", ge.ping)
}

func (ge *TestsEndpoint) checkExternal() bool {
	isHealthy := true

	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()

		if !conn.IsConnected() {
			isHealthy = false
		}
	}

	if remote_signer.VaultStorage {
		vm := vaultManager.MakeVaultManager(remote_signer.KeyPrefix)
		health, err := vm.HealthStatus()

		if err != nil {
			isHealthy = false
		}

		if !health.Initialized || health.Sealed {
			isHealthy = false
		}
	}

	return isHealthy
}

func (ge *TestsEndpoint) ping(w http.ResponseWriter, r *http.Request) {
	isHealthy := ge.checkExternal()

	// Do not log here. This call will flood the log
	w.Header().Set("Content-Type", models.MimeText)

	if isHealthy {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("OK"))
	} else {
		w.WriteHeader(503)
		_, _ = w.Write([]byte("Service Unavailable"))
	}
}
