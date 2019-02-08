package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/quan-to/graphql"
	"github.com/quan-to/handler"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/server/agent"
	"net/http"
)

var amLog = SLog.Scope("AgentAdmin")

type AgentAdmin struct {
	tm      etc.TokenManager
	handler *handler.Handler
	ctx     context.Context
}

func MakeAgentAdmin(tm etc.TokenManager, am etc.AuthManager) *AgentAdmin {
	schemaConfig := graphql.SchemaConfig{
		Query:    agent.RootManagementQuery,
		Mutation: agent.RootManagementMutation,
	}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		amLog.Fatal(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: false,
	})

	return &AgentAdmin{
		handler: h,
		tm:      tm,
		ctx: remote_signer.ContextWithValues(context.Background(), map[string]interface{}{
			agent.TokenManagerKey: tm,
			agent.AuthManagerKey:  am,
		}),
	}
}

func (admin *AgentAdmin) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	admin.handler.ContextHandler(admin.ctx, w, r)
}

func (admin *AgentAdmin) AddHandlers(r *mux.Router) {
	r.HandleFunc("", admin.handleGraphQL)
	r.HandleFunc("/", admin.handleGraphQL)
}
