package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/fieldcipher"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/openpgp"
	"github.com/quan-to/slog"
	"net/http"
)

var jfcLog = slog.Scope("JFC Endpoint")

type JFCEndpoint struct {
	sm  etc.SMInterface
	gpg etc.PGPInterface
}

func MakeJFCEndpoint(sm etc.SMInterface, gpg etc.PGPInterface) *JFCEndpoint {
	return &JFCEndpoint{
		sm:  sm,
		gpg: gpg,
	}
}

func (jfc *JFCEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/cipher", jfc.cipher).Methods("POST")
	r.HandleFunc("/decipher", jfc.decipher).Methods("POST")
}

func (jfc *JFCEndpoint) cipher(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	var data models.FieldCipherInput

	if !UnmarshalBodyOrDie(&data, w, r, sksLog) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, sksLog)
		}
	}()

	keys := make([]*openpgp.Entity, 0)

	for i, v := range data.Keys {
		k := jfc.gpg.GetPublicKeyEntity(v)
		if k == nil {
			NotFound(fmt.Sprintf("data.Keys[%d]", i), fmt.Sprintf("publickey for fingerPrint %s was not found", v), w, r, jfcLog)
			return
		}
		keys = append(keys, k)
	}

	if len(keys) == 0 {
		InvalidFieldData("data.Keys", "no keys specified", w, r, jfcLog)
		return
	}

	cipher := fieldcipher.MakeCipher(keys)

	packet, err := cipher.GenerateEncryptedPacket(data.JSON, data.SkipFields)

	if err != nil {
		InternalServerError(err.Error(), err, w, r, jfcLog)
		return
	}

	d, _ := json.Marshal(packet)

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(d))
	LogExit(jfcLog, r, 200, n)
}

func (jfc *JFCEndpoint) decipher(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	var data models.FieldDecipherInput

	if !UnmarshalBodyOrDie(&data, w, r, sksLog) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, sksLog)
		}
	}()

	keys := jfc.gpg.GetPrivate(data.KeyFingerprint)
	if len(keys) == 0 {
		NotFound("keyFingerprint", fmt.Sprintf("There is no such key %s or its not decrypted.", data.KeyFingerprint), w, r, jfcLog)
		return
	}

	decipher, err := fieldcipher.MakeDecipher(keys)

	if err != nil {
		jfcLog.Error(err)
		InternalServerError("Error processing your request. Please try again.", err, w, r, jfcLog)
		return
	}

	dec, err := decipher.DecipherPacket(fieldcipher.CipherPacket{
		EncryptedKey:  data.EncryptedKey,
		EncryptedJSON: data.EncryptedJSON,
	})

	if err != nil {
		jfcLog.Error(err)
		InvalidFieldData("payload", err.Error(), w, r, jfcLog)
		return
	}

	d, _ := json.Marshal(dec)

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(d))
	LogExit(jfcLog, r, 200, n)
}
