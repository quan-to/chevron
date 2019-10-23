package server

import (
	"bytes"
	"crypto"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	remote_signer "github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/slog"
)

type AgentProxy struct {
	gpg       etc.PGPInterface
	transport *http.Transport
	tm        etc.TokenManager
	log       slog.Instance
}

// MakeAgentProxy creates an instance of agent proxy endpoint
func MakeAgentProxy(log slog.Instance, gpg etc.PGPInterface, tm etc.TokenManager) *AgentProxy {
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

func (proxy *AgentProxy) defaultHandler(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(proxy.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	h := r.Header

	if r.Method == http.MethodOptions {
		targetUrl := remote_signer.AgentTargetURL

		if h.Get("serverUrl") != "" {
			targetUrl = h.Get("serverUrl")
		}

		log = log.WithFields(map[string]interface{}{
			"targetUrl": targetUrl,
		})

		client := &http.Client{
			Transport: proxy.transport,
		}

		req, err := http.NewRequest(r.Method, targetUrl, nil)

		if err != nil {
			InternalServerError("There was an error processing your request", err.Error(), w, r, log)
			return
		}

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

		log.Await("Sending request to %s", targetUrl)
		res, err := client.Do(req)
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
		n, _ := io.Copy(w, res.Body)
		LogExit(log, r, res.StatusCode, int(n))
	} else {
		targetUrl := remote_signer.AgentTargetURL

		if h.Get("serverUrl") != "" {
			targetUrl = h.Get("serverUrl")
		}

		log = log.WithFields(map[string]interface{}{
			"targetUrl": targetUrl,
		})

		token := ""

		if !remote_signer.AgentBypassLogin {
			if h.Get("proxyToken") == "" {
				PermissionDenied("proxyToken", "Please check if your proxyToken is valid", w, r, log)
				return
			}

			token = h.Get("proxyToken")
			h.Del("proxyToken")

			log.Await("Verifying user token")
			err := proxy.tm.Verify(token)
			log.Done("Token verified")

			if err != nil {
				PermissionDenied("proxyToken", "Please check if your proxyToken is valid", w, r, log)
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

		jsondata["_timestamp"] = time.Now().Unix() * 1000
		jsondata["_timeUniqueId"] = uuid.New().String()

		bodyData, _ = json.Marshal(jsondata)

		req, err := http.NewRequest(r.Method, targetUrl, bytes.NewBuffer(bodyData))

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

		log.Await("Sending request to %s", targetUrl)
		res, err := client.Do(req)
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
		n, _ := io.Copy(w, res.Body)
		LogExit(log, r, res.StatusCode, int(n))
	}
}

func (proxy *AgentProxy) AddHandlers(r *mux.Router) {
	r.HandleFunc("/", proxy.defaultHandler)
	r.HandleFunc("", proxy.defaultHandler)
}
