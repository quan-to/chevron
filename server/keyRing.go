package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/models"
	"net/http"
)

var kreLog = SLog.Scope("KeyRing Endpoint")

type KeyRingEndpoint struct {
	sm  etc.SMInterface
	gpg etc.PGPInterface
}

func MakeKeyRingEndpoint(sm etc.SMInterface, gpg etc.PGPInterface) *KeyRingEndpoint {
	return &KeyRingEndpoint{
		sm:  sm,
		gpg: gpg,
	}
}

func (kre *KeyRingEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/getKey", kre.getKey).Methods("GET")
	r.HandleFunc("/cachedKeys", kre.getCachedKeys).Methods("GET")
	r.HandleFunc("/privateKeys", kre.getLoadedPrivateKeys).Methods("GET")
	r.HandleFunc("/addPrivateKey", kre.addPrivateKey).Methods("POST")
}

func (kre *KeyRingEndpoint) getKey(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	q := r.URL.Query()

	fingerPrint := q.Get("fingerPrint")

	key := kre.gpg.GetPublicKeyAscii(fingerPrint)

	if key == "" {
		NotFound("fingerPrint", fmt.Sprintf("Key with fingerPrint %s was not found", fingerPrint), w, r, kreLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(key))
	LogExit(geLog, r, 200, n)
}

func (kre *KeyRingEndpoint) getCachedKeys(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	cachedKeys := kre.gpg.GetCachedKeys()

	bodyData, err := json.Marshal(cachedKeys)

	if err != nil {
		InternalServerError("There was an error processing your request. Please try again.", nil, w, r, kreLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(geLog, r, 200, n)
}

func (kre *KeyRingEndpoint) getLoadedPrivateKeys(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	privateKeys := kre.gpg.GetLoadedPrivateKeys()

	bodyData, err := json.Marshal(privateKeys)

	if err != nil {
		InternalServerError("There was an error processing your request. Please try again.", nil, w, r, kreLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(geLog, r, 200, n)
}

func (kre *KeyRingEndpoint) addPrivateKey(w http.ResponseWriter, r *http.Request) {
	var data models.KeyRingAddPrivateKeyData
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	if !UnmarshalBodyOrDie(&data, w, r, geLog) {
		return
	}

	fp, err := remote_signer.GetFingerPrintFromKey(data.EncryptedPrivateKey)

	if err != nil {
		InvalidFieldData("EncryptedPrivateKey", "Invalid key provided. Check if its in ASCII Armored Format. Cannot read fingerprint", w, r, kreLog)
		return
	}

	err, n := kre.gpg.LoadKey(data.EncryptedPrivateKey)

	if err != nil {
		InvalidFieldData("EncryptedPrivateKey", "Invalid key provided. Check if its in ASCII Armored Format", w, r, kreLog)
		return
	}

	if n == 0 {
		NotFound("EncryptedPrivateKey", "No private keys found at specified payload", w, r, kreLog)
		return
	}

	fingerPrint, _ := remote_signer.GetFingerPrintFromKey(data.EncryptedPrivateKey)

	if data.Password != nil {
		pass := data.Password.(string)
		err := kre.gpg.UnlockKey(fp, pass)
		if err != nil {
			InvalidFieldData("Password", "Invalid password for the key provided", w, r, kreLog)
			return
		}
	}

	if data.SaveToDisk {
		err = kre.gpg.SavePrivateKey(fingerPrint, data.EncryptedPrivateKey, data.Password)
		if err != nil {
			InternalServerError("There was an error saving your key to disk.", data, w, r, kreLog)
		}
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ = w.Write([]byte("OK"))
	LogExit(geLog, r, 200, n)
}
