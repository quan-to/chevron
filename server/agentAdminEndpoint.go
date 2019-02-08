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

type graphIntercept struct {
	originalHandler http.ResponseWriter
	WrittenBytes    int
	StatusCode      int
}

func (gi *graphIntercept) Header() http.Header {
	return gi.originalHandler.Header()
}

func (gi *graphIntercept) Write(data []byte) (int, error) {
	n, err := gi.originalHandler.Write(data)
	gi.WrittenBytes += n
	return n, err
}

func (gi *graphIntercept) WriteHeader(statusCode int) {
	gi.StatusCode = statusCode
	gi.originalHandler.WriteHeader(statusCode)
}

func (admin *AgentAdmin) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, amLog)
		}
	}()

	gi := graphIntercept{originalHandler: w, StatusCode: http.StatusOK}
	ctx := context.WithValue(admin.ctx, agent.HTTPRequestKey, r)

	token := r.Header.Get("proxyToken")

	if token != "" {
		err := admin.tm.Verify(token)
		if err != nil {
			InvalidFieldData("proxyToken", "The specified proxyToken is either invalid or expired.", w, r, amLog)
			return
		}

		user := admin.tm.GetUserData(token)
		ctx = context.WithValue(ctx, agent.LoggedUserKey, user)
	}

	admin.handler.ContextHandler(ctx, &gi, r)
	LogExit(amLog, r, gi.StatusCode, gi.WrittenBytes)
}

func (admin *AgentAdmin) AddHandlers(r *mux.Router) {
	r.HandleFunc("", admin.handleGraphQL)
	r.HandleFunc("/", admin.handleGraphQL)
}
