package server

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/models"
	"net/http"
)

var geLog = SLog.Scope("GPG Endpoint")

type GPGEndpoint struct {
	sm  *remote_signer.SecretsManager
	gpg *remote_signer.PGPManager
}

func MakeGPGEndpoint(sm *remote_signer.SecretsManager, gpg *remote_signer.PGPManager) *GPGEndpoint {
	return &GPGEndpoint{
		sm:  sm,
		gpg: gpg,
	}
}

func (ge *GPGEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/generateKey", ge.generateKey).Methods("POST")
	r.HandleFunc("/unlockKey", ge.unlockKey).Methods("POST")
	r.HandleFunc("/sign", ge.sign).Methods("POST")
	r.HandleFunc("/signQuanto", ge.signQuanto).Methods("POST")
}

func (ge *GPGEndpoint) sign(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)
	var data models.GPGSignData

	if !UnmarshalBodyOrDie(&data, w, r, geLog) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	bytes, err := base64.StdEncoding.DecodeString(data.Base64Data)

	if err != nil {
		InvalidFieldData("Base64Data", err.Error(), w, r, geLog)
		return
	}

	signature, err := ge.gpg.SignData(data.FingerPrint, bytes, crypto.SHA512)

	if err != nil {
		InvalidFieldData("Key", fmt.Sprintf("There was an error signing your data: %s", err.Error()), w, r, geLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(signature))
	LogExit(geLog, r, 200, n)
}

func (ge *GPGEndpoint) signQuanto(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)
	var data models.GPGSignData

	if !UnmarshalBodyOrDie(&data, w, r, geLog) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	bytes, err := base64.StdEncoding.DecodeString(data.Base64Data)

	if err != nil {
		InvalidFieldData("Base64Data", err.Error(), w, r, geLog)
		return
	}

	signature, err := ge.gpg.SignData(data.FingerPrint, bytes, crypto.SHA512)

	if err != nil {
		InvalidFieldData("Key", fmt.Sprintf("There was an error signing your data: %s", err.Error()), w, r, geLog)
		return
	}

	quantoSig := remote_signer.GPG2Quanto(signature, data.FingerPrint, "SHA512")

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(quantoSig))
	LogExit(geLog, r, 200, n)
}

func (ge *GPGEndpoint) unlockKey(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)
	var data models.GPGUnlockKeyData

	if !UnmarshalBodyOrDie(&data, w, r, geLog) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	err := ge.gpg.UnlockKey(data.FingerPrint, data.Password)

	if err != nil {
		InvalidFieldData("Password/Key", fmt.Sprintf("There is no such key %s or the password is invalid.", data.FingerPrint), w, r, geLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(geLog, r, 200, n)
}

func (ge *GPGEndpoint) generateKey(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)
	var data models.GPGGenerateKeyData

	if !UnmarshalBodyOrDie(&data, w, r, geLog) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	if data.Bits < remote_signer.MinKeyBits {
		InvalidFieldData("Bits", fmt.Sprintf("The key should be at least %d bits length.", remote_signer.MinKeyBits), w, r, geLog)
		return
	}

	if len(data.Password) == 0 {
		InvalidFieldData("Password", "You should provide a password.", w, r, geLog)
		return
	}

	key, err := ge.gpg.GeneratePGPKey(data.Identifier, data.Password, data.Bits)

	if err != nil {
		InternalServerError("There was an error generating your key. Please try again.", err, w, r, geLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(key))
	LogExit(geLog, r, 200, n)
}
