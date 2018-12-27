package remote_signer

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/models"
	"net/http"
)

var geLog = SLog.Scope("GPG Endpoint")

type GPGEndpoint struct {
	sm  *SecretsManager
	gpg *PGPManager
}

func MakeGPGEndpoint(sm *SecretsManager, gpg *PGPManager) *GPGEndpoint {
	return &GPGEndpoint{
		sm:  sm,
		gpg: gpg,
	}
}

func (ge *GPGEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/generateKey", ge.generateKey).Methods("POST")
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

	if data.Bits < minKeyBits {
		InvalidFieldData("Bits", fmt.Sprintf("The key should be at least %d bits length.", minKeyBits), w, r, geLog)
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
