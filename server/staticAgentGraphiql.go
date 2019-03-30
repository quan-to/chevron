package server

import (
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/server/agent"
	"net/http"
	"path"
	"strings"
)

var sgLog = SLog.Scope("GraphiQL")

type StaticGraphiQL struct{}

func MakeStaticGraphiQL() *StaticGraphiQL {
	return &StaticGraphiQL{}
}

func (gql *StaticGraphiQL) displayFile(filename string, w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, sgLog)
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
		InternalServerError("There was an internal server error. Please try again", nil, w, r, sksLog)
		return
	}

	w.WriteHeader(200)
	n, _ := w.Write([]byte(fileData))
	LogExit(sgLog, r, 200, n)
}

func (gql *StaticGraphiQL) AttachHandlers(r *mux.Router) {
	files, _ := agent.AssetDir("bundle")

	for _, v := range files {
		filePath := path.Join("/", v)
		sgLog.Debug("Attaching %s", filePath)
		r.HandleFunc(filePath, func(w http.ResponseWriter, r *http.Request) {
			gql.displayFile(filePath, w, r)
		})
	}

	sgLog.Debug("Attaching /")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gql.displayFile("/index.html", w, r)
	})
	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		gql.displayFile("/index.html", w, r)
	})
}
