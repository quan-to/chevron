package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/slog"
	"net/http"
)

var intLog = slog.Scope("Internal Endpoint")

type InternalEndpoint struct {
	sm  etc.SMInterface
	gpg etc.PGPInterface
}

func MakeInternalEndpoint(sm etc.SMInterface, gpg etc.PGPInterface) *InternalEndpoint {
	return &InternalEndpoint{
		sm:  sm,
		gpg: gpg,
	}
}

func (ie *InternalEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/__triggerKeyUnlock", ie.triggerKeyUnlock)
	r.HandleFunc("/__getUnlockPasswords", ie.getUnlockPasswords).Methods("GET")
	r.HandleFunc("/__postEncryptedPasswords", ie.postUnlockPasswords).Methods("POST")
}

func (ie *InternalEndpoint) triggerKeyUnlock(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, intLog)
		}
	}()

	ie.sm.UnlockLocalKeys(ie.gpg)

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(intLog, r, 200, n)
}

func (ie *InternalEndpoint) getUnlockPasswords(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, intLog)
		}
	}()

	passwords := ie.sm.GetPasswords()

	bodyData, _ := json.Marshal(passwords)

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(intLog, r, 200, n)
}

func (ie *InternalEndpoint) postUnlockPasswords(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	var passwords map[string]string

	if !UnmarshalBodyOrDie(&passwords, w, r, intLog) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, intLog)
		}
	}()

	for k, v := range passwords {
		ie.sm.PutEncryptedPassword(k, v)
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(intLog, r, 200, n)
}
