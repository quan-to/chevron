package main

import "C"
import (
	"github.com/quan-to/chevron/chevronlib"
	"strings"
)

// LoadKey loads a private or public key into the memory keyring
//export LoadKey
func LoadKey(keyData *C.char, result *C.char, resultLen C.int) (err C.int, loadedPrivateKeys C.int) {
	goKeyData := C.GoString(keyData)
	rLen := int(resultLen)

	l, e := chevronlib.LoadKey(goKeyData)
	loadedPrivateKeys = C.int(l)
	err = OK
	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR, C.int(0)
	}
	return
}

// UnlockKey unlocks a private key to be used
//export UnlockKey
func UnlockKey(fingerprint, password *C.char, result *C.char, resultLen C.int) C.int {
	goFingerprint := C.GoString(fingerprint)
	goPassword := C.GoString(password)
	e := chevronlib.UnlockKey(goFingerprint, goPassword)
	rLen := int(resultLen)

	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}
	return OK
}

// VerifySignature verifies a signature using a already loaded public key
//export VerifySignature
func VerifySignature(data *C.char, dataLen C.int, signature *C.char, result *C.char, resultLen C.int) C.int {
	goData := make([]byte, int(dataLen))
	copyFromCToGo(goData, data, int(dataLen))
	goSignature := C.GoString(signature)

	rLen := int(resultLen)

	r, e := chevronlib.VerifySignature(goData, goSignature)

	if e != nil { // Error
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	if r { // Signature Valid
		return TRUE
	}

	// Signature Invalid
	return FALSE
}

// QuantoVerifySignature verifies a signature in Quanto Signature Format using a already loaded public key
//export QuantoVerifySignature
func QuantoVerifySignature(data *C.char, dataLen C.int, signature *C.char, result *C.char, resultLen C.int) C.int {
	goData := make([]byte, int(dataLen))
	copyFromCToGo(goData, data, int(dataLen))
	goSignature := C.GoString(signature)

	rLen := int(resultLen)

	r, e := chevronlib.QuantoVerifySignature(goData, goSignature)

	if e != nil { // Error
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	if r { // Signature Valid
		return TRUE
	}

	// Signature Invalid
	return FALSE
}

// VerifyBase64DataSignature verifies a signature using a already loaded public key. The b64data is a raw binary data encoded in base64 string
//export VerifyBase64DataSignature
func VerifyBase64DataSignature(b64data, signature *C.char, result *C.char, resultLen C.int) C.int {
	goB64Data := C.GoString(b64data)
	goSignature := C.GoString(signature)
	rLen := int(resultLen)
	r, e := chevronlib.VerifyBase64DataSignature(goB64Data, goSignature)

	if e != nil { // Error
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	if r { // Signature Valid
		return TRUE
	}

	// Signature Invalid
	return FALSE
}

// QuantoVerifyBase64DataSignature verifies a signature in Quanto Signature Format using a already loaded public key.
// The b64data is a raw binary data encoded in base64 string
//export QuantoVerifyBase64DataSignature
func QuantoVerifyBase64DataSignature(b64data, signature *C.char, result *C.char, resultLen C.int) C.int {
	goB64Data := C.GoString(b64data)
	goSignature := C.GoString(signature)
	rLen := int(resultLen)
	r, e := chevronlib.QuantoVerifyBase64DataSignature(goB64Data, goSignature)

	if e != nil { // Error
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	if r { // Signature Valid
		return TRUE
	}

	// Signature Invalid
	return FALSE
}

// SignData signs data using a already loaded and unlocked private key
//export SignData
func SignData(data *C.char, dataLen C.int, fingerprint *C.char, result *C.char, resultLen C.int) C.int {
	goData := make([]byte, int(dataLen))
	copyFromCToGo(goData, data, int(dataLen))
	goFingerprint := C.GoString(fingerprint)

	rLen := int(resultLen)

	r, e := chevronlib.SignData(goData, goFingerprint)
	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	copyStringToC(result, []byte(r), rLen)

	return OK
}

// QuantoSignData signs data using a already loaded and unlocked private key and returns in Quanto Signature Format
//export QuantoSignData
func QuantoSignData(data *C.char, dataLen C.int, fingerprint *C.char, result *C.char, resultLen C.int) C.int {
	goData := make([]byte, int(dataLen))
	copyFromCToGo(goData, data, int(dataLen))
	goFingerprint := C.GoString(fingerprint)

	rLen := int(resultLen)

	r, e := chevronlib.QuantoSignData(goData, goFingerprint)
	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	copyStringToC(result, []byte(r), rLen)

	return OK
}

// SignBase64Data signs data using a already loaded and unlocked private key.
// The b64data is a raw binary data encoded in base64 string
//export SignBase64Data
func SignBase64Data(b64data, fingerprint *C.char, result *C.char, resultLen C.int) C.int {
	goB64Data := C.GoString(b64data)
	goFingerprint := C.GoString(fingerprint)
	r, e := chevronlib.SignBase64Data(goB64Data, goFingerprint)

	rLen := int(resultLen)

	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	copyStringToC(result, []byte(r), rLen)

	return OK
}

// SignBase64Data signs data using a already loaded and unlocked private key. Returns in Quanto Signature Format
// The b64data is a raw binary data encoded in base64 string
//export QuantoSignBase64Data
func QuantoSignBase64Data(b64data, fingerprint *C.char, result *C.char, resultLen C.int) C.int {
	goB64Data := C.GoString(b64data)
	goFingerprint := C.GoString(fingerprint)
	r, e := chevronlib.QuantoSignBase64Data(goB64Data, goFingerprint)

	rLen := int(resultLen)

	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	copyStringToC(result, []byte(r), rLen)

	return OK
}

// GetKeyFingerprints returns all fingerprints in CSV format from a ASCII Armored PGP Keychain
//export GetKeyFingerprints
func GetKeyFingerprints(keyData *C.char, result *C.char, resultLen C.int) C.int {
	goKeyData := C.GoString(keyData)
	r, e := chevronlib.GetKeyFingerprints(goKeyData)

	rLen := int(resultLen)

	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	resultString := strings.Join(r, ",")

	copyStringToC(result, []byte(resultString), rLen)

	return OK
}

// ChangeKeyPassword re-encrypts the input key using newPassword
//export ChangeKeyPassword
func ChangeKeyPassword(keyData, currentPassword, newPassword *C.char, result *C.char, resultLen C.int) C.int {
	goKeyData := C.GoString(keyData)
	goCurrentPassword := C.GoString(currentPassword)
	goNewPassword := C.GoString(newPassword)

	r, e := chevronlib.ChangeKeyPassword(goKeyData, goCurrentPassword, goNewPassword)

	rLen := int(resultLen)

	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	copyStringToC(result, []byte(r), rLen)

	return OK
}

// GetPublicKey returns the cached public key from the specified fingerprint
//export GetPublicKey
func GetPublicKey(fingerprint *C.char, result *C.char, resultLen C.int) C.int {
	goFingerprint := C.GoString(fingerprint)
	r, e := chevronlib.GetPublicKey(goFingerprint)

	rLen := int(resultLen)

	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	copyStringToC(result, []byte(r), rLen)

	return OK
}

// GenerateKey generates a new key using specified bits and identifier and encrypts it using the specified password
//export GenerateKey
func GenerateKey(password, identifier *C.char, bits C.int, result *C.char, resultLen C.int) C.int {
	goPassword := C.GoString(password)
	goIdentifier := C.GoString(identifier)
	goInt := int(bits)

	r, e := chevronlib.GenerateKey(goPassword, goIdentifier, goInt)

	rLen := int(resultLen)

	if e != nil {
		copyStringToC(result, []byte(e.Error()), rLen)
		return ERROR
	}

	copyStringToC(result, []byte(r), rLen)

	return OK
}
