package main

import "C"
import (
	"github.com/quan-to/chevron/chevronlib"
)

// LoadKey loads a private or public key into the memory keyring
//export LoadKey
func LoadKey(keyData string) (loadedPrivateKeys int, err string) {
	l, e := chevronlib.LoadKey(keyData)
	loadedPrivateKeys = l
	if e != nil {
		err = e.Error()
	}
	return
}

// UnlockKey unlocks a private key to be used
//export UnlockKey
func UnlockKey(fingerprint, password string) (err string) {
	e := chevronlib.UnlockKey(fingerprint, password)
	if e != nil {
		err = e.Error()
	}
	return
}

// VerifySignature verifies a signature using a already loaded public key
//export VerifySignature
func VerifySignature(data []byte, signature string) (result bool, err string) {
	r, e := chevronlib.VerifySignature(data, signature)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

// VerifyBase64DataSignature verifies a signature using a already loaded public key. The b64data is a raw binary data encoded in base64 string
//export VerifyBase64DataSignature
func VerifyBase64DataSignature(b64data, signature string) (result bool, err string) {
	r, e := chevronlib.VerifyBase64DataSignature(b64data, signature)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

// SignData signs data using a already loaded and unlocked private key
//export SignData
func SignData(data []byte, fingerprint string) (result string, err string) {
	r, e := chevronlib.SignData(data, fingerprint)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

// SignBase64Data signs data using a already loaded and unlocked private key. The b64data is a raw binary data encoded in base64 string
//export SignBase64Data
func SignBase64Data(b64data, fingerprint string) (result string, err string) {
	r, e := chevronlib.SignBase64Data(b64data, fingerprint)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

// GetKeyFingerprints returns all fingerprints in a ASCII Armored PGP Keychain
//export GetKeyFingerprints
func GetKeyFingerprints(keyData string) (result []string, err string) {
	r, e := chevronlib.GetKeyFingerprints(keyData)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

// ChangeKeyPassword re-encrypts the input key using newPassword
//export ChangeKeyPassword
func ChangeKeyPassword(keyData, currentPassword, newPassword string) (result string, err string) {
	r, e := chevronlib.ChangeKeyPassword(keyData, currentPassword, newPassword)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

// GetPublicKey returns the cached public key from the specified fingerprint
//export GetPublicKey
func GetPublicKey(fingerprint string) (result string, err string) {
	r, e := chevronlib.GetPublicKey(fingerprint)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

// GenerateKey generates a new key using specified bits and identifier and encrypts it using the specified password
//export GenerateKey
func GenerateKey(password, identifier string, bits int) (result string, err string) {
	r, e := chevronlib.GenerateKey(password, identifier, bits)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

func main() {}
