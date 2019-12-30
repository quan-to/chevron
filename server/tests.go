package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	remote_signer "github.com/quan-to/chevron"
	"github.com/quan-to/chevron/database"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/vaultManager"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type TestsEndpoint struct {
	log slog.Instance
	vm  *vaultManager.VaultManager
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
		vm:  vm,
	}
}

func (ge *TestsEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/ping", ge.ping)
}

func (ge *TestsEndpoint) checkExternal(ctx context.Context) bool {
	requestID := remote_signer.GetRequestIDFromContext(ctx)
	log := ge.log.Tag(requestID)
	isHealthy := true

	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()

		_, err := r.Expr(1).Run(conn)
		if err != nil {
			log.Error(err)
			isHealthy = false
		}
	}

	if ge.vm != nil {
		health, err := ge.vm.HealthStatus()

		if err != nil {
			log.Error(err)
			return false
		}

		if health != nil && !health.Initialized || health.Sealed {
			log.Info("Vault initialized? %t, is sealed? %t", health.Initialized, health.Sealed)
			return false
		}
	}

	return isHealthy
}

func (ge *TestsEndpoint) ping(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	isHealthy := ge.checkExternal(ctx)

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
