package server

import (
	"bytes"
	"crypto"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/slog"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var agentLog = slog.Scope("Agent")

type AgentProxy struct {
	gpg       etc.PGPInterface
	transport *http.Transport
	tm        etc.TokenManager
}

func MakeAgentProxy(gpg etc.PGPInterface, tm etc.TokenManager) *AgentProxy {
	return &AgentProxy{
		gpg: gpg,
		transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
		tm: tm,
	}
}

func (proxy *AgentProxy) defaultHandler(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, intLog)
		}
	}()

	h := r.Header

	targetUrl := remote_signer.AgentTargetURL

	if h.Get("serverUrl") != "" {
		targetUrl = h.Get("serverUrl")
	}

	token := ""

	if !remote_signer.AgentBypassLogin {
		if h.Get("proxyToken") == "" {
			PermissionDenied("proxyToken", "Please check if your proxyToken is valid", w, r, agentLog)
			return
		}

		token = h.Get("proxyToken")
		h.Del("proxyToken")

		err := proxy.tm.Verify(token)

		if err != nil {
			PermissionDenied("proxyToken", "Please check if your proxyToken is valid", w, r, agentLog)
			return
		}
	}

	client := &http.Client{
		Transport: proxy.transport,
	}

	fingerPrint := remote_signer.AgentKeyFingerPrint

	if !remote_signer.AgentBypassLogin {
		user := proxy.tm.GetUserData(token)
		fingerPrint = user.GetFingerPrint()
	}

	bodyData, err := ioutil.ReadAll(r.Body)

	if err != nil {
		InternalServerError("There was an error processing your request", err.Error(), w, r, agentLog)
		return
	}

	var jsondata map[string]interface{}

	err = json.Unmarshal(bodyData, &jsondata)

	if err != nil {
		InternalServerError("There was an error processing your request", err.Error(), w, r, agentLog)
		return
	}

	jsondata["_timestamp"] = time.Now().Unix() * 1000
	jsondata["_timeUniqueId"] = uuid.New().String()

	bodyData, _ = json.Marshal(jsondata)

	req, err := http.NewRequest(r.Method, targetUrl, bytes.NewBuffer(bodyData))

	if err != nil {
		InternalServerError("There was an error processing your request", err.Error(), w, r, agentLog)
		return
	}

	signature, err := proxy.gpg.SignData(fingerPrint, bodyData, crypto.SHA512)

	if err != nil {
		InternalServerError("There was an error signing your request", err.Error(), w, r, agentLog)
		return
	}

	quantoSig := remote_signer.GPG2Quanto(signature, fingerPrint, "SHA512")

	req.Header.Add("signature", quantoSig)
	req.Header.Add("X-Powered-By", "RemoteSigner Agent")

	for k, v := range r.Header {
		if len(v) > 1 {
			for _, t := range v {
				req.Header.Add(k, t)
			}
		} else {
			req.Header.Set(k, v[0])
		}
	}

	res, err := client.Do(req)

	if err != nil {
		InternalServerError("There was an error processing your request", err.Error(), w, r, agentLog)
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

	n, _ := io.Copy(w, res.Body)
	LogExit(geLog, r, res.StatusCode, int(n))
}

func (proxy *AgentProxy) AddHandlers(r *mux.Router) {
	r.HandleFunc("/", proxy.defaultHandler)
	r.HandleFunc("", proxy.defaultHandler)
}
