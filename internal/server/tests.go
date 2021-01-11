package server

import (
	"context"
	"net/http"

	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/internal/vaultManager"
	"github.com/quan-to/chevron/pkg/models"

	"github.com/gorilla/mux"
	"github.com/quan-to/slog"
)

type HealthCheckHandler interface {
	HealthCheck() error
}

type TestsEndpoint struct {
	log slog.Instance
	vm  *vaultManager.VaultManager
	db  HealthCheckHandler
}

// MakeTestsEndpoint creates an instance of healthcheck tests endpoint
func MakeTestsEndpoint(log slog.Instance, vm *vaultManager.VaultManager, dbHandler HealthCheckHandler) *TestsEndpoint {
	if log == nil {
		log = slog.Scope("Tests")
	} else {
		log = log.SubScope("Tests")
	}

	return &TestsEndpoint{
		log: log,
		vm:  vm,
		db:  dbHandler,
	}
}

func (ge *TestsEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/ping", ge.ping)
}

func (ge *TestsEndpoint) checkExternal(ctx context.Context) bool {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := ge.log.Tag(requestID)
	isHealthy := true

	if ge.db != nil {
		err := ge.db.HealthCheck()
		if err != nil {
			log.Error("Database Health Check error: %s", err)
			isHealthy = false
		}
	}

	if ge.vm != nil {
		health, err := ge.vm.HealthStatus()

		if err != nil {
			log.Error("Vault Health Check error: %s", err)
			return false
		}

		if health != nil && !health.Initialized || health.Sealed {
			log.Info("Vault initialized? %t, is sealed? %t", health.Initialized, health.Sealed)
			return false
		}
	}

	return isHealthy
}

// Health Check godoc
// @id tests-hc-ping
// @tags Tests
// @Summary Checks if Remote-Signer and all its dependencies are working
// @Produce plain
// @Success 200 {string} result "OK"
// @Failure 503 {string} result "Service Unavailable"
// @Router /tests/ping [get]
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
