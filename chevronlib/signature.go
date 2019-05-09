package chevronlib

import (
	"crypto"
	"encoding/base64"
)

// LoadKey loads a private or public key into the memory keyring
// export LoadKey
func LoadKey(keyData string) (loadedPrivateKeys int, err error) {
	err, loadedPrivateKeys = pgpBackend.LoadKey(keyData)
	return
}

// UnlockKey unlocks a private key to be used
// export UnlockKey
func UnlockKey(fingerprint, password string) (err error) {
	return pgpBackend.UnlockKey(fingerprint, password)
}

// VerifySignature verifies a signature using a already loaded public key
// export VerifySignature
func VerifySignature(data []byte, signature string) (result bool, err error) {
	return pgpBackend.VerifySignature(data, signature)
}

// VerifyBase64DataSignature verifies a signature using a already loaded public key. The b64data is a raw binary data encoded in base64 string
// export VerifyBase64DataSignature
func VerifyBase64DataSignature(b64data, signature string) (result bool, err error) {
	var data []byte
	data, err = base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return
	}

	return VerifySignature(data, signature)
}

// SignData signs data using a already loaded and unlocked private key
// export SignData
func SignData(data []byte, fingerprint string) (result string, err error) {
	return pgpBackend.SignData(fingerprint, data, crypto.SHA512)
}

// SignBase64Data signs data using a already loaded and unlocked private key. The b64data is a raw binary data encoded in base64 string
// export SignBase64Data
func SignBase64Data(b64data, fingerprint string) (result string, err error) {
	var data []byte
	data, err = base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return
	}

	return SignData(data, fingerprint)
}
