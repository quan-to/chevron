package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/quan-to/chevron/pkg/fieldcipher"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"

	"github.com/gorilla/mux"
	"github.com/quan-to/chevron/pkg/openpgp"
	"github.com/quan-to/slog"
)

type JFCEndpoint struct {
	sm  interfaces.SecretsManager
	gpg interfaces.PGPManager
	log slog.Instance
}

// MakeJFCEndpoint creates a handler for Json Field Cipher Endpoints
func MakeJFCEndpoint(log slog.Instance, sm interfaces.SecretsManager, gpg interfaces.PGPManager) *JFCEndpoint {
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

// Field Cipher godoc
// @id field-cipher-cipher
// @tags Field Cipher
// @Summary Encrypts JSON fields to specified GPG keys
// @Accept json
// @Produce json
// @param message body models.FieldCipherInput true "The encryption parameters"
// @Success 200 {object} fieldcipher.CipherPacket
// @Failure default {object} QuantoError.ErrorObject
// @Router /fieldCipher/cipher [post]
func (jfc *JFCEndpoint) cipher(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(jfc.log, r)

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
		k := jfc.gpg.GetPublicKeyEntity(ctx, v)
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
	w.Write([]byte(d))
}

// Field Decipher godoc
// @id field-cipher-decipher
// @tags Field Cipher
// @Summary Decrypts JSON fields from specified GPG keys.
// @Accept json
// @Produce json
// @param message body models.FieldDecipherInput true "The decryption parameters"
// @Success 200 {object} fieldcipher.DecipherPacket
// @Failure default {object} QuantoError.ErrorObject
// @Router /fieldCipher/decipher [post]
func (jfc *JFCEndpoint) decipher(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(jfc.log, r)

	var data models.FieldDecipherInput

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	keys := jfc.gpg.GetPrivate(ctx, data.KeyFingerprint)
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
	w.Write([]byte(d))
}
