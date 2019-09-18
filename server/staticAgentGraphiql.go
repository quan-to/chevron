package server

import (
	"github.com/gorilla/mux"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/server/agent"
	"github.com/quan-to/slog"
	"net/http"
	"path"
	"strings"
)

type StaticGraphiQL struct {
	log slog.Instance
}

func MakeStaticGraphiQL(log slog.Instance) *StaticGraphiQL {
	if log == nil {
		log = slog.Scope("GraphiQL")
	} else {
		log = log.SubScope("GraphiQL")
	}

	return &StaticGraphiQL{
		log: log,
	}
}

func (gql *StaticGraphiQL) displayFile(filename string, w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(gql.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	fileData, err := agent.Asset("bundle" + filename)

	if strings.Index(filename, "index.htm") > -1 {
		// Add server URL
		f := string(fileData)
		f = strings.Replace(f, "{SERVER_URL}", remote_signer.AgentTargetURL, -1)
		f = strings.Replace(f, "{AGENT_URL}", remote_signer.AgentExternalURL, -1)
		f = strings.Replace(f, "{AGENT_ADMIN_URL}", remote_signer.AgentAdminExternalURL, -1)
		fileData = []byte(f)
	}

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, log)
		return
	}

	w.WriteHeader(200)
	n, _ := w.Write([]byte(fileData))
	LogExit(log, r, 200, n)
}

func (gql *StaticGraphiQL) AttachHandlers(r *mux.Router) {
	files, _ := agent.AssetDir("bundle")

	for _, v := range files {
		filePath := path.Join("/", v)
		gql.log.Debug("Attaching %s", filePath)
		r.HandleFunc(filePath, func(w http.ResponseWriter, r *http.Request) {
			gql.displayFile(filePath, w, r)
		})
	}

	gql.log.Debug("Attaching /")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gql.displayFile("/index.html", w, r)
	})
	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		gql.displayFile("/index.html", w, r)
	})
}
