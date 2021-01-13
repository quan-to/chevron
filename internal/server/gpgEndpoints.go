package server

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"

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

// Decrypt godoc
// @id gpg-data-decrypt
// @tags GPG Operations
// @Summary Decrypts data using the specified GPG Key. The private key should be previously loaded.
// @Accept json
// @Produce json
// @Param message body models.GPGDecryptData true "Information to decrypt"
// @Success 200 {object} models.GPGDecryptedData
// @Failure default {object} QuantoError.ErrorObject
// @Router /gpg/decrypt [post]
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
	n, _ := w.Write(d)
	LogExit(log, r, 200, n)
}

// Encrypt godoc
// @id gpg-data-encrypt
// @tags GPG Operations
// @Summary Encrypts data for the specified GPG Public Key
// @Accept json
// @Produce json
// @Param message body models.GPGEncryptData true "Information to encrypt to public key"
// @Success 200 {string} Encrypted Data "wcDMA8HPMfuMKotZAQwADzmQgwJiz3p5suaYpPwCbOluqvu2O5kVitJNO86KfkSYgbR0y67c+fGk5nO+Zm66qeolXLqVBHUvSnpZf9jMupRZLRmSZ0JmmvXoJIdiahj+NLwF6NVBvmoJ8BkMEQkr5oCNkKBveaCYXdQ7Gba2buICwxxwEmq3LV6/D0Zg4AmKX/k2N1kjRGJaUeHH3oU1YEjPo3A3bo9EZLGLI+J5VSlxkydxXUkF2TISKCr2rkhUmH5E7CUFu6H2nOofxk9tJDoSfjACkEjFKdg3BbTqNlYeuNmdJHwLfHDI+WcbL3/Hsl5MVnyHGeztsj0jn2bAIcT9FHfw1W3LUpaTNlemfrn52la7zN3r2588JDRbSaqLQ/d5+3hHWyE7RsRL0jdpEj/HM3ue2mi6GfyxDZy1DxdZsy7kqoYbBIwbtCdqZetU+bH6hWk92BY89AJUpV7xPCzRozw5WvCTsPYsu10JDvvPvj1c47BA9KlJ1wTcB2lYhmoX39T3ymjMKJ+6NAOF0uAB5PToGBs3BjE4MsxQHMLchK3hTuXg+uAY4fVU4I3jFyDPs8zYKsfgCOIHYBV84Obhm9rgqOAh4Ifi+klQeOCf4+p0IGeF6b6+4IPiAtTxRuB+5KnAAWAlBpwJWAqwNJ68HIjiN9UOgeGU+wA="
// @Failure default {object} QuantoError.ErrorObject
// @Router /gpg/encrypt [post]
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

// VerifySignature godoc
// @id gpg-data-verify
// @tags GPG Operations
// @Summary Verifies a signature in the standard GPG format
// @Accept json
// @Produce json
// @Param message body models.GPGVerifySignatureDataNonQuanto true "Information to verify a signature in GPG format"
// @Success 200 {string} Returns OK
// @Failure default {object} QuantoError.ErrorObject
// @Router /gpg/verifySignature [post]
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

// VerifySignatureQuanto godoc
// @id gpg-data-verify-quanto
// @tags GPG Operations
// @Summary Verifies a signature in Quanto's signature format
// @Accept json
// @Produce json
// @Param message body models.GPGVerifySignatureData true "Information to verify a signature in quanto format"
// @Success 200 {string} Returns OK
// @Failure default {object} QuantoError.ErrorObject
// @Router /gpg/verifySignatureQuanto [post]
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

// Sign godoc
// @id gpg-data-sign
// @tags GPG Operations
// @Summary Signs a payload with a standard GPG signature format
// @Description Signs a payload using the specified GPG key and returns the signature in GPG Format
// @Accept json
// @Produce plain
// @Param message body models.GPGSignData true "Data to sign"
// @Success 200 {string} Signature "-----BEGIN PGP SIGNATURE-----\n\nwsDcBAABCgAQBQJf+LriCRAFUfRSq+RjpAAAuL0MAGGrSJfK/tnMkwZ2Rkh3JcvF\nE8WU8jwc8quz+0p9gMDscby0jShJ2G2XXMm3WAYXW88J6h8u2E/lTb6l3oBq/FPb\n15gTM5Ie0p0kHBUlgP5bkV9EF9+VQif40fhVX7OPrS27jWtVNP374ARzSIgKMLa6\nKBZhV1eQecLIlEYXahUP9jyt4cR4A4d9P+YJS/L6d/tQT4g9DBo66hYt5lu4sagG\nDHsW2HK9I7fizCBaE8azLtQd3RRFTWZshln7OGVypwcdbzWbYr5uEhituxAnZKS4\nSWwI0hgj1OkZeOhKwaydtITnaeH+nmlLBzhGKQWjCiLlsDNkkp3/4FKOuYJkYXeZ\nm61GV6G5ZpW/gFVJXXyPz6ElNfWCorZQvxLbY4YWTBLdLyblHnp9kshav6dnexN1\nwQyBDk8jxucmKNE8kCu591dPj/g/H38/zpGZQhj8Firb0rCFumqsAwxFeyTEFjVI\ncyDHa5K+ytmSrITIdQUUsp1M4UQiRH63c1HYOLQurw==\n=BRZt\n-----END PGP SIGNATURE-----"
// @Failure default {object} QuantoError.ErrorObject
// @Router /gpg/sign [post]
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

// SignQuanto godoc
// @id gpg-data-sign-quanto
// @tags GPG Operations
// @Summary Signs a payload with a Quanto's signature format
// @Description Signs a payload using the specified GPG key and returns the signature in Quanto Format
// @Accept json
// @Produce plain
// @Param message body models.GPGSignData true "Data to sign"
// @Success 200 {string} Signature 0551F452ABE463A4_SHA512_wsDcBAABCgAQBQJf+LYnCRAFUfRSq+RjpAAA7/oMACHJPMtQs4rr0uxX4AMZ8akb+x2p5ZYL+uRug+zctp82sJEJmL76HG++UyzDmMUCagJ+LBWp2RcCQvfsIhX5MqD7lPkEdtl0uNCIU40apvzn1+0kndl7LnFtzyHMWrHrRqEFGJ0E2APPqv7g1pehVKeusMOkTNUmmsJNgZBYrluZxHnai/Rudoe9jBxihY4ALF0eOyTCHbtWy0z6fll3Bo/iPe777kplDXmTBzCEM8uD3/VZmY6pGn6oXUov/z8Dcrg2x5qT4i5DgdF8OSLbsxVW2OIV8DwCicQCT2tK95fctBqJ22vfmhNlxI3KzI9ShxeV6Eci5p5Zydgoh77pDiWDysrq1dOZ+o7T+ij72K3s63w3loERFVoDxDuKG3jS3+fj+ggqqtpUpm957+9+4QlnJqZk0v9TKT661HnoH4MfZR3muBir8/dgF4mNtuQLSswOxdVs1sHSC3ssTIzzpQqeI2iy3m8Svgl5unAdv2QE81EM/wT5brc2R/abSRz52A===J34T
// @Failure default {object} QuantoError.ErrorObject
// @Router /gpg/signQuanto [post]
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

// UnlockKey godoc
// @id gpg-key-unlock
// @tags GPG Operations
// @Summary Unlocks a pre-loaded GPG Private Key
// @Description Unlocks a locked pre-loaded key inside remote signer
// @Accept json
// @Produce plain
// @Param message body models.GPGUnlockKeyData true "Unlock Data"
// @Success 200 {string} Result Returns OK on success
// @Failure default {object} QuantoError.ErrorObject
// @Router /gpg/unlockKey [post]
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

// GenerateKey godoc
// @id gpg-key-generate
// @tags GPG Operations
// @Summary Generates a new GPG Key pair
// @Description Generates a new GPG Key by specifying the Identifier, Bits and Password
// @Accept json
// @Produce json
// @Param message body models.GPGGenerateKeyData true "Information to generate the key. The minimum acceptable bits is 2048."
// @Success 200 {string} Generated GPG Key
// @Failure default {object} QuantoError.ErrorObject
// @Router /gpg/generateKey [post]
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
