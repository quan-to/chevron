package server

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/quan-to/chevron/internal/models"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/interfaces"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/quan-to/slog"
)

type GPGEndpoint struct {
	sm  interfaces.SecretsManager
	gpg interfaces.PGPManager
	log slog.Instance
}

// MakeGPGEndpoint Creates an instance of an endpoint that handles GPG Calls
func MakeGPGEndpoint(log slog.Instance, sm interfaces.SecretsManager, gpg interfaces.PGPManager) *GPGEndpoint {
	if log == nil {
		log = slog.Scope("GPG (HTTP)")
	} else {
		log = log.SubScope("GPG (HTTP)")
	}

	return &GPGEndpoint{
		sm:  sm,
		gpg: gpg,
		log: log,
	}
}

func (ge *GPGEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/generateKey", ge.generateKey).Methods("POST")
	r.HandleFunc("/unlockKey", ge.unlockKey).Methods("POST")
	r.HandleFunc("/sign", ge.sign).Methods("POST")
	r.HandleFunc("/signQuanto", ge.signQuanto).Methods("POST")
	r.HandleFunc("/verifySignature", ge.verifySignature).Methods("POST")
	r.HandleFunc("/verifySignatureQuanto", ge.verifySignatureQuanto).Methods("POST")
	r.HandleFunc("/encrypt", ge.encrypt).Methods("POST")
	r.HandleFunc("/decrypt", ge.decrypt).Methods("POST")
}

func (ge *GPGEndpoint) decrypt(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ge.log, r)
	InitHTTPTimer(log, r)

	var data models.GPGDecryptData
	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	decrypted, err := ge.gpg.Decrypt(ctx, data.AsciiArmoredData, data.DataOnly)

	if err != nil {
		InvalidFieldData("Decryption", fmt.Sprintf("Error decrypting data: %s", err.Error()), w, r, log)
		return
	}

	d, _ := json.Marshal(*decrypted)

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(d))
	LogExit(log, r, 200, n)
}

func (ge *GPGEndpoint) encrypt(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ge.log, r)
	InitHTTPTimer(log, r)
	var data models.GPGEncryptData

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	bytes, err := base64.StdEncoding.DecodeString(data.Base64Data)

	if err != nil {
		InvalidFieldData("Base64Data", err.Error(), w, r, log)
		return
	}

	encrypted, err := ge.gpg.Encrypt(ctx, data.Filename, data.FingerPrint, bytes, data.DataOnly)

	if err != nil {
		InvalidFieldData("Encryption", fmt.Sprintf("Error encrypting data: %s", err.Error()), w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(encrypted))
	LogExit(log, r, 200, n)
}

func (ge *GPGEndpoint) verifySignature(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ge.log, r)
	InitHTTPTimer(log, r)
	var data models.GPGVerifySignatureData

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	bytes, err := base64.StdEncoding.DecodeString(data.Base64Data)

	if err != nil {
		InvalidFieldData("Base64Data", err.Error(), w, r, log)
		return
	}

	valid, err := ge.gpg.VerifySignature(ctx, bytes, data.Signature)

	if err != nil {
		InvalidFieldData("Signature", err.Error(), w, r, log)
		return
	}

	if !valid {
		InvalidFieldData("Signature", "The provided signature is invalid", w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(log, r, 200, n)
}

func (ge *GPGEndpoint) verifySignatureQuanto(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ge.log, r)
	InitHTTPTimer(log, r)
	var data models.GPGVerifySignatureData

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	bytes, err := base64.StdEncoding.DecodeString(data.Base64Data)

	if err != nil {
		InvalidFieldData("Base64Data", err.Error(), w, r, log)
		return
	}

	signature := tools.Quanto2GPG(data.Signature)
	valid, err := ge.gpg.VerifySignature(ctx, bytes, signature)

	if err != nil {
		if strings.Contains(err.Error(), "cannot find public key to verify signature") {
			NotFound("publicKey", err.Error(), w, r, log)
			return
		}
		InvalidFieldData("Signature", err.Error(), w, r, log)
		return
	}

	if !valid {
		InvalidFieldData("Signature", "The provided signature is invalid", w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(log, r, 200, n)
}

func (ge *GPGEndpoint) sign(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ge.log, r)
	InitHTTPTimer(log, r)
	var data models.GPGSignData

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	bytes, err := base64.StdEncoding.DecodeString(data.Base64Data)

	if err != nil {
		InvalidFieldData("Base64Data", err.Error(), w, r, log)
		return
	}

	signature, err := ge.gpg.SignData(ctx, data.FingerPrint, bytes, crypto.SHA512)

	if err != nil {
		InvalidFieldData("Key", fmt.Sprintf("There was an error signing your data: %s", err.Error()), w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(signature))
	LogExit(log, r, 200, n)
}

func (ge *GPGEndpoint) signQuanto(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ge.log, r)
	InitHTTPTimer(log, r)
	var data models.GPGSignData

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	bytes, err := base64.StdEncoding.DecodeString(data.Base64Data)

	if err != nil {
		InvalidFieldData("Base64Data", err.Error(), w, r, log)
		return
	}

	signature, err := ge.gpg.SignData(ctx, data.FingerPrint, bytes, crypto.SHA512)

	if err != nil {
		InvalidFieldData("Key", fmt.Sprintf("There was an error signing your data: %s", err.Error()), w, r, log)
		return
	}

	quantoSig := tools.GPG2Quanto(signature, data.FingerPrint, "SHA512")

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(quantoSig))
	LogExit(log, r, 200, n)
}

func (ge *GPGEndpoint) unlockKey(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ge.log, r)
	InitHTTPTimer(log, r)
	var data models.GPGUnlockKeyData

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	err := ge.gpg.UnlockKey(ctx, data.FingerPrint, data.Password)

	if err != nil {
		InvalidFieldData("Password/Key", fmt.Sprintf("There is no such key %s or the password is invalid.", data.FingerPrint), w, r, log)
		return
	}

	fp := ge.gpg.FixFingerPrint(data.FingerPrint)

	ge.sm.PutKeyPassword(ctx, fp, data.Password)

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(log, r, 200, n)
}

func (ge *GPGEndpoint) generateKey(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(ge.log, r)
	InitHTTPTimer(log, r)
	var data models.GPGGenerateKeyData

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	if data.Bits < ge.gpg.MinKeyBits() {
		InvalidFieldData("Bits", fmt.Sprintf("The key should be at least %d bits length.", ge.gpg.MinKeyBits()), w, r, log)
		return
	}

	if len(data.Password) == 0 {
		InvalidFieldData("Password", "You should provide a password.", w, r, log)
		return
	}

	key, err := ge.gpg.GeneratePGPKey(ctx, data.Identifier, data.Password, data.Bits)

	if err != nil {
		InternalServerError("There was an error generating your key. Please try again.", err.Error(), w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(key))
	LogExit(log, r, 200, n)
}
