package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/internal/server/pages"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
)

type KeyRingEndpoint struct {
	sm  interfaces.SecretsManager
	gpg interfaces.PGPManager
	log slog.Instance
	dbh DatabaseHandler
}

// MakeKeyRingEndpoint creates an instance of key ring management endpoints
func MakeKeyRingEndpoint(log slog.Instance, sm interfaces.SecretsManager, gpg interfaces.PGPManager, dbHandler DatabaseHandler) *KeyRingEndpoint {
	if log == nil {
		log = slog.Scope("KeyRing")
	} else {
		log = log.SubScope("KeyRing")
	}

	return &KeyRingEndpoint{
		sm:  sm,
		gpg: gpg,
		log: log,
		dbh: dbHandler,
	}
}

func (kre *KeyRingEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/getKey", kre.getKey).Methods("GET")
	r.HandleFunc("/cachedKeys", kre.getCachedKeys).Methods("GET")
	r.HandleFunc("/privateKeys", kre.getLoadedPrivateKeys).Methods("GET")
	r.HandleFunc("/addPrivateKey", kre.addPrivateKey).Methods("POST")
	r.HandleFunc("/addPrivateKey", pages.ServeAddPrivateKey).Methods("GET")
	r.HandleFunc("/deletePrivateKey", kre.deletePrivateKey).Methods("POST")
}

// Get GPG Key godoc
// @id kre-get-key
// @tags Key Ring
// @Summary Fetches a GPG public key
// @Produce plain
// @param fingerPrint query string true "Fingerprint of the GPG Key to be fetched"
// @Success 200 {string} result "GPG public key"
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/getKey [get]
func (kre *KeyRingEndpoint) getKey(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(kre.log, r)
	ctx = wrapContextWithDatabaseHandler(kre.dbh, ctx)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()

	fingerPrint := q.Get("fingerPrint")

	key, _ := kre.gpg.GetPublicKeyASCII(ctx, fingerPrint)

	if key == "" {
		NotFound("fingerPrint", fmt.Sprintf("Key with fingerPrint %s was not found", fingerPrint), w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(key))
	LogExit(log, r, 200, n)
}

// Get Cached Keys godoc
// @id kre-get-cached-keys
// @tags Key Ring
// @Summary Fetches a list of cached keys
// @Produce json
// @Success 200 {object} []models.KeyInfo
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/cachedKeys [get]
func (kre *KeyRingEndpoint) getCachedKeys(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(kre.log, r)
	ctx = wrapContextWithDatabaseHandler(kre.dbh, ctx)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	cachedKeys := kre.gpg.GetCachedKeys(ctx)

	bodyData, err := json.Marshal(cachedKeys)

	if err != nil {
		log.Error("Error getting cached keys: %s", err)
		InternalServerError("There was an error processing your request. Please try again.", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

// Get Loaded Private Keys godoc
// @id kre-get-loaded-private-keys
// @tags Key Ring
// @Summary Fetches a list of loaded private keys
// @Produce json
// @Success 200 {object} []models.KeyInfo
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/privateKeys [get]
func (kre *KeyRingEndpoint) getLoadedPrivateKeys(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(kre.log, r)
	ctx = wrapContextWithDatabaseHandler(kre.dbh, ctx)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	privateKeys := kre.gpg.GetLoadedPrivateKeys(ctx)

	bodyData, err := json.Marshal(privateKeys)

	if err != nil {
		InternalServerError("There was an error processing your request. Please try again.", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

// Delete Private Key godoc
// @id kre-del-private-key
// @tags Key Ring, Key Store
// @Summary Deletes a GPG Private Key
// @Accepts json
// @Produce json
// @param message body models.KeyRingDeletePrivateKeyData true "Private Key Information"
// @Success 200 {object} models.GPGDeletePrivateKeyReturn
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/deletePrivateKey [post]
func (kre *KeyRingEndpoint) deletePrivateKey(w http.ResponseWriter, r *http.Request) {
	var data models.KeyRingDeletePrivateKeyData
	ctx := wrapContextWithRequestID(r)
	ctx = wrapContextWithDatabaseHandler(kre.dbh, ctx)
	log := wrapLogWithRequestID(kre.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	err := kre.gpg.DeleteKey(ctx, data.FingerPrint)
	if err != nil {
		log.Error("Error deleting key: %s", err)
		InternalServerError("There was an error deleting your key from the disk.", data, w, r, log)
		return
	}

	ret := models.GPGDeletePrivateKeyReturn{
		Status: "OK",
	}

	d, _ := json.Marshal(ret)

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write(d)
	LogExit(log, r, 200, n)
}

// Add Private Key godoc
// @id kre-add-private-key
// @tags Key Ring, Key Store
// @Summary Adds a GPG Private Key
// @Accepts json
// @Produce json
// @param message body models.KeyRingAddPrivateKeyData true "Private Key Information"
// @Success 200 {object} models.GPGAddPrivateKeyReturn
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/addPrivateKey [post]
func (kre *KeyRingEndpoint) addPrivateKey(w http.ResponseWriter, r *http.Request) {
	var data models.KeyRingAddPrivateKeyData
	ctx := wrapContextWithRequestID(r)
	ctx = wrapContextWithDatabaseHandler(kre.dbh, ctx)
	log := wrapLogWithRequestID(kre.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	fp, err := tools.GetFingerPrintFromKey(data.EncryptedPrivateKey)

	if err != nil {
		InvalidFieldData("EncryptedPrivateKey", "Invalid key provided. Check if its in ASCII Armored Format. Cannot read fingerprint", w, r, log)
		return
	}

	n, _ := kre.gpg.LoadKey(ctx, data.EncryptedPrivateKey) // Error never happens here due GetFingerPrintFromKey

	if n == 0 {
		NotFound("EncryptedPrivateKey", "No private keys found at specified payload", w, r, log)
		return
	}

	fingerPrint, _ := tools.GetFingerPrintFromKey(data.EncryptedPrivateKey)

	if data.Password != nil {
		pass := data.Password.(string)
		err := kre.gpg.UnlockKey(ctx, fp, pass)
		if err != nil {
			InvalidFieldData("Password", "Invalid password for the key provided", w, r, log)
			return
		}
	}

	pubKey, _ := kre.gpg.GetPublicKeyASCII(ctx, fp)

	log.Info("Adding public key for %s on PKS", fp)
	res := keymagic.PKSAdd(ctx, pubKey)
	log.Info("PKS Add Key: %s", res)

	if data.SaveToDisk {
		err = kre.gpg.SaveKey(fingerPrint, data.EncryptedPrivateKey, data.Password)
		if err != nil {
			log.Error("Error saving key: %s", err)
			InternalServerError("There was an error saving your key to disk.", data, w, r, log)
			return
		}
	}

	ret := models.GPGAddPrivateKeyReturn{
		FingerPrint: fp,
		PublicKey:   pubKey,
	}

	d, _ := json.Marshal(ret)

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ = w.Write(d)
	LogExit(log, r, 200, n)
}
