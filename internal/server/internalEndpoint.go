package server

import (
	"encoding/json"
	"github.com/quan-to/chevron/internal/models"
	"github.com/quan-to/chevron/pkg/interfaces"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quan-to/slog"
)

type InternalEndpoint struct {
	sm  interfaces.SecretsManager
	gpg interfaces.PGPManager
	log slog.Instance
}

// MakeInternalEndpoint creates an instance to handle internal control endpoints such as key password data
func MakeInternalEndpoint(log slog.Instance, sm interfaces.SecretsManager, gpg interfaces.PGPManager) *InternalEndpoint {
	if log == nil {
		log = slog.Scope("Internal")
	} else {
		log = log.SubScope("Internal")
	}

	return &InternalEndpoint{
		sm:  sm,
		gpg: gpg,
		log: log,
	}
}

func (ie *InternalEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/__triggerKeyUnlock", ie.triggerKeyUnlock)
	r.HandleFunc("/__getUnlockPasswords", ie.getUnlockPasswords).Methods("GET")
	r.HandleFunc("/__postEncryptedPasswords", ie.postUnlockPasswords).Methods("POST")
}

func (ie *InternalEndpoint) triggerKeyUnlock(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ie.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	ie.sm.UnlockLocalKeys(ctx, ie.gpg)

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(log, r, 200, n)
}

func (ie *InternalEndpoint) getUnlockPasswords(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ie.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	passwords := ie.sm.GetPasswords(ctx)

	bodyData, _ := json.Marshal(passwords)

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

func (ie *InternalEndpoint) postUnlockPasswords(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ie.log, r)
	InitHTTPTimer(log, r)

	var passwords map[string]string

	if !UnmarshalBodyOrDie(&passwords, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	for k, v := range passwords {
		ie.sm.PutEncryptedPassword(ctx, k, v)
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(log, r, 200, n)
}
