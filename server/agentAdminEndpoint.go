package server

import (
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer/etc"
	"net/http"
)

type AgentAdmin struct {
	tm etc.TokenManager
}

func (admin *AgentAdmin) login(w http.ResponseWriter, r *http.Request) {

}

func (admin *AgentAdmin) AddHandlers(r *mux.Router) {
	//r.HandleFunc("/", admin.defaultHandler)
}
