package server

import (
	"github.com/gorilla/mux"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/database"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/vaultManager"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"net/http"
)

type TestsEndpoint struct {
	log slog.Instance
	vm *vaultManager.VaultManager
}

// MakeTestsEndpoint creates an instance of healthcheck tests endpoint
func MakeTestsEndpoint(log slog.Instance, vm *vaultManager.VaultManager) *TestsEndpoint {
	if log == nil {
		log = slog.Scope("Tests")
	} else {
		log = log.SubScope("Tests")
	}

	return &TestsEndpoint{
		log: log,
		vm: vm,
	}
}

func (ge *TestsEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/ping", ge.ping)
}

func (ge *TestsEndpoint) checkExternal() bool {
	isHealthy := true

	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()

		_, err := r.Expr(1).Run(conn)
		if err != nil {
			ge.log.Error(err)
			isHealthy = false
		}
	}

	if remote_signer.VaultStorage {
		health, err := ge.vm.HealthStatus()

		if err != nil {
			ge.log.Error(err)
			return false
		}

		if health != nil && !health.Initialized || health.Sealed {
			ge.log.Info("Vault initialized? %t, is sealed? %t", health.Initialized, health.Sealed)
			return false
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
