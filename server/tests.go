package server

import (
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer/models"
	"net/http"
)

type TestsEndpoint struct{}

func MakeTestsEndpoint() *TestsEndpoint {
	return &TestsEndpoint{}
}

func (ge *TestsEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/ping", ge.ping)
}

func (ge *TestsEndpoint) ping(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)
	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(geLog, r, 200, n)
}
