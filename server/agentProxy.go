package server

import (
	"bytes"
	"crypto"
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/etc"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var agentLog = SLog.Scope("Agent")

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

	q := r.URL.Query()

	targetUrl := remote_signer.AgentTargetURL

	if q.Get("serverUrl") != "" {
		targetUrl = q.Get("serverUrl")
	}

	token := ""

	if !remote_signer.AgentBypassLogin {
		if q.Get("proxyToken") == "" {
			PermissionDenied("proxyToken", "Please check if your proxyToken is valid", w, r, agentLog)
			return
		}

		token = q.Get("proxyToken")
		q.Del("proxyToken")

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
