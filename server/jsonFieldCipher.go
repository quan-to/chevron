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

type JFCEndpoint struct {
	sm  etc.SMInterface
	gpg etc.PGPInterface
	log slog.Instance
}

func MakeJFCEndpoint(log slog.Instance, sm etc.SMInterface, gpg etc.PGPInterface) *JFCEndpoint {
	if log == nil {
		log = slog.Scope("JFC")
	} else {
		log = log.SubScope("JFC")
	}

	return &JFCEndpoint{
		sm:  sm,
		gpg: gpg,
		log: log,
	}
}

func (jfc *JFCEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/cipher", jfc.cipher).Methods("POST")
	r.HandleFunc("/decipher", jfc.decipher).Methods("POST")
}

func (jfc *JFCEndpoint) cipher(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(jfc.log, r)
	InitHTTPTimer(log, r)

	var data models.FieldCipherInput

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	keys := make([]*openpgp.Entity, 0)

	for i, v := range data.Keys {
		k := jfc.gpg.GetPublicKeyEntity(v)
		if k == nil {
			NotFound(fmt.Sprintf("data.Keys[%d]", i), fmt.Sprintf("publickey for fingerPrint %s was not found", v), w, r, log)
			return
		}
		keys = append(keys, k)
	}

	if len(keys) == 0 {
		InvalidFieldData("data.Keys", "no keys specified", w, r, log)
		return
	}

	cipher := fieldcipher.MakeCipher(keys)

	packet, err := cipher.GenerateEncryptedPacket(data.JSON, data.SkipFields)

	if err != nil {
		InternalServerError(err.Error(), err, w, r, log)
		return
	}

	d, _ := json.Marshal(packet)

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(d))
	LogExit(log, r, 200, n)
}

func (jfc *JFCEndpoint) decipher(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(jfc.log, r)
	InitHTTPTimer(log, r)

	var data models.FieldDecipherInput

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	keys := jfc.gpg.GetPrivate(data.KeyFingerprint)
	if len(keys) == 0 {
		NotFound("keyFingerprint", fmt.Sprintf("There is no such key %s or its not decrypted.", data.KeyFingerprint), w, r, log)
		return
	}

	decipher, err := fieldcipher.MakeDecipher(keys)

	if err != nil {
		log.Error(err)
		InternalServerError("Error processing your request. Please try again.", err, w, r, log)
		return
	}

	dec, err := decipher.DecipherPacket(fieldcipher.CipherPacket{
		EncryptedKey:  data.EncryptedKey,
		EncryptedJSON: data.EncryptedJSON,
	})

	if err != nil {
		log.Error(err)
		InvalidFieldData("payload", err.Error(), w, r, log)
		return
	}

	d, _ := json.Marshal(dec)

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(d))
	LogExit(log, r, 200, n)
}
