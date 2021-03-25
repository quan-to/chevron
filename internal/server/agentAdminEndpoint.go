package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/handler"
	"github.com/quan-to/chevron/internal/server/agent"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
)

type AgentAdmin struct {
	tm      interfaces.TokenManager
	handler *handler.Handler
	ctx     context.Context
	log     slog.Instance
}

// MakeAgentAdmin creates an instance of Agent Administration endpoint
func MakeAgentAdmin(log slog.Instance, tm interfaces.TokenManager, am interfaces.AuthManager) *AgentAdmin {
	if log == nil {
		log = slog.Scope("AgentAdmin")
	} else {
		log = log.SubScope("AgentAdmin")
	}

	schemaConfig := graphql.SchemaConfig{
		Query:    agent.RootManagementQuery,
		Mutation: agent.RootManagementMutation,
	}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: false,
		FormatErrorFn: func(err error) gqlerrors.FormattedError {
			switch err := err.(type) {
			case *gqlerrors.Error:
				log.Error("%+v", err.OriginalError)
				return gqlerrors.FormatError(err)
			case gqlerrors.ExtendedError:
				log.Error("%+v", err.Error())
				return gqlerrors.FormatError(err)
			case *QuantoError.ErrorObject:
				return err.ToFormattedError()
			default:
				log.Error("%+v", err.Error())
				return gqlerrors.FormatError(err)
			}
		},
	})

	return &AgentAdmin{
		handler: h,
		tm:      tm,
		ctx: tools.ContextWithValues(context.Background(), map[string]interface{}{
			agent.TokenManagerKey: tm,
			agent.AuthManagerKey:  am,
		}),
		log: log,
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

// Agent Admin GraphQL godoc
// @id agent-proxy-admin
// @tags Agent
// @Summary Administration of the Agent proxy tokens.
// @Accept json
// @Produce json
// @param proxyToken header string true "Proxy Token of the admin user. It is required for all calls besides the login"
// @param message body string true "The JSON content of the graphql query"
// @Success 200 {string} result "result of the query"
// @Failure default {object} QuantoError.ErrorObject
// @Router /agentAdmin [post]
func (admin *AgentAdmin) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(admin.log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	gi := graphIntercept{originalHandler: w, StatusCode: http.StatusOK}
	ctx := context.WithValue(admin.ctx, agent.HTTPRequestKey, r)

	token := r.Header.Get("proxyToken")

	if token != "" {
		err := admin.tm.Verify(token)
		if err != nil {
			InvalidFieldData("proxyToken", "The specified proxyToken is either invalid or expired.", w, r, log)
			return
		}

		user := admin.tm.GetUserData(token)
		ctx = context.WithValue(ctx, agent.LoggedUserKey, user)
	}
	admin.handler.ContextHandler(ctx, &gi, r)
}

func (admin *AgentAdmin) AddHandlers(r *mux.Router) {
	r.HandleFunc("", admin.handleGraphQL)
	r.HandleFunc("/", admin.handleGraphQL)
}
