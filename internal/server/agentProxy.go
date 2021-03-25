package server

import (
	"bytes"
	"crypto"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/uuid"
	"github.com/quan-to/slog"
)

type AgentProxy struct {
	gpg       interfaces.PGPManager
	transport *http.Transport
	tm        interfaces.TokenManager
	log       slog.Instance
}

// MakeAgentProxy creates an instance of agent proxy endpoint
func MakeAgentProxy(log slog.Instance, gpg interfaces.PGPManager, tm interfaces.TokenManager) *AgentProxy {
	if log == nil {
		log = slog.Scope("Agent")
	} else {
		log = log.SubScope("Agent")
	}

	return &AgentProxy{
		gpg: gpg,
		transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
		tm:  tm,
		log: log,
	}
}

func injectUniquenessFields(log slog.Instance, json map[string]interface{}) error {
	uniqueString := uuid.EnsureUUID(log)

	json["_timeUniqueId"] = uniqueString
	json["_timestamp"] = time.Now().UnixNano() / 1e6
	log.DebugNote("Request UUID: %q - RequestTimestamp: %d", json["_timeUniqueId"], json["_timestamp"])

	return nil
}

// Agent Proxy Call godoc
// @id agent-proxy-call
// @tags Agent
// @Summary Signs the request with GPG key specified by the token
// @Accept json
// @Produce json
// @param proxyToken header string false "Proxy Token generated with agentAdmin. It is required if running with authentication enabled"
// @param serverUrl header string false "Target server URL. Defaults to environment variable AGENT_TARGET_URL"
// @param message body string true "POST Content to send signed to the target server. The message body will be sent to the target server with it's signature in a header field named 'signature'."
// @Success 200 {string} result "result of the query"
// @Failure default {object} QuantoError.ErrorObject
// @Router /agent [post]
func (proxy *AgentProxy) defaultHandler(w http.ResponseWriter, r *http.Request) {
	var res *http.Response
	var req *http.Request
	var err error

	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(proxy.log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	h := r.Header

	targetURL := config.AgentTargetURL

	if h.Get("serverUrl") != "" {
		targetURL = h.Get("serverUrl")
	}

	log = log.WithFields(map[string]interface{}{
		"targetURL": targetURL,
	})

	client := &http.Client{
		Transport: proxy.transport,
	}

	if r.Method == http.MethodOptions {
		req, err = http.NewRequest(r.Method, targetURL, nil)

		if err != nil {
			InternalServerError("There was an error processing your request", err.Error(), w, r, log)
			return
		}

		req.Header.Add("X-Powered-By", "RemoteSigner Agent")
	} else {
		token := ""

		if !config.AgentBypassLogin {
			if h.Get("proxyToken") == "" {
				PermissionDenied("proxyToken", "Please check if your proxyToken is valid", w, r, log)
				return
			}

			token = h.Get("proxyToken")
			h.Del("proxyToken")

			log.Await("Verifying user token")
			err = proxy.tm.Verify(token)
			log.Done("Token verified")

			if err != nil {
				PermissionDenied("proxyToken", "Please check if your proxyToken is valid", w, r, log)
				return
			}
		}

		fingerPrint := config.AgentKeyFingerPrint

		if !config.AgentBypassLogin {
			user := proxy.tm.GetUserData(token)
			fingerPrint = user.GetFingerPrint()
		}

		log.DebugAwait("Reading body")
		bodyData, err := ioutil.ReadAll(r.Body)
		log.DebugDone("Body read")

		if err != nil {
			InternalServerError("There was an error processing your request", err.Error(), w, r, log)
			return
		}

		var jsondata map[string]interface{}

		err = json.Unmarshal(bodyData, &jsondata)

		if err != nil {
			InternalServerError("There was an error processing your request", err.Error(), w, r, log)
			return
		}

		err = injectUniquenessFields(log, jsondata)

		if err != nil {
			InternalServerError("There was an error ensuring request uniqueness", err.Error(), w, r, log)
			return
		}

		bodyData, _ = json.Marshal(jsondata)

		req, err = http.NewRequest(r.Method, targetURL, bytes.NewBuffer(bodyData))

		if err != nil {
			InternalServerError("There was an error processing your request", err.Error(), w, r, log)
			return
		}

		log.Await("Signing data with %s", fingerPrint)
		signature, err := proxy.gpg.SignData(ctx, fingerPrint, bodyData, crypto.SHA512)
		log.Done("Data signed")

		if err != nil {
			InternalServerError("There was an error signing your request", err.Error(), w, r, log)
			return
		}

		quantoSig := tools.GPG2Quanto(signature, fingerPrint, "SHA512")

		req.Header.Add("signature", quantoSig)
		req.Header.Add("X-Powered-By", "RemoteSigner Agent")
	}

	for k, v := range r.Header {
		if len(v) > 1 {
			for _, t := range v {
				req.Header.Add(k, t)
			}
		} else {
			req.Header.Set(k, v[0])
		}
	}

	log.Await("Sending request to %s", targetURL)
	res, err = client.Do(req)
	log.Done("Received response")

	if err != nil {
		InternalServerError("There was an error processing your request", err.Error(), w, r, log)
		return
	}

	for k, v := range res.Header {
		if len(v) > 1 {
			for _, t := range v {
				w.Header().Add(k, t)
			}
		} else {
			w.Header().Set(k, v[0])
		}
	}

	log.Info("Sending response")
	w.WriteHeader(res.StatusCode)
	_, _ = io.Copy(w, res.Body)
}

func (proxy *AgentProxy) AddHandlers(r *mux.Router) {
	r.HandleFunc("/", proxy.defaultHandler)
	r.HandleFunc("", proxy.defaultHandler)
}
