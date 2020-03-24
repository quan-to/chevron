package chevronlib

import (
	"crypto"
	"encoding/base64"
	remote_signer "github.com/quan-to/chevron"
)

// LoadKey loads a private or public key into the memory keyring
// export LoadKey
func LoadKey(keyData string) (loadedPrivateKeys int, err error) {
	err, loadedPrivateKeys = pgpBackend.LoadKey(ctx, keyData)
	return
}

// UnlockKey unlocks a private key to be used
// export UnlockKey
func UnlockKey(fingerprint, password string) (err error) {
	return pgpBackend.UnlockKey(ctx, fingerprint, password)
}

// VerifySignature verifies a signature using a already loaded public key
// export VerifySignature
func VerifySignature(data []byte, signature string) (result bool, err error) {
	return pgpBackend.VerifySignature(ctx, data, signature)
}

// QuantoVerifySignature verifies a signature in Quanto Signature Format using a already loaded public key
// export VerifySignature
func QuantoVerifySignature(data []byte, signature string) (result bool, err error) {
	signature = remote_signer.Quanto2GPG(signature)

	return pgpBackend.VerifySignature(ctx, data, signature)
}

// QuantoVerifyBase64DataSignature verifies a signature using a already loaded public key.
// The b64data is a raw binary data encoded in base64 string
// export VerifyBase64DataSignature
func QuantoVerifyBase64DataSignature(b64data, signature string) (result bool, err error) {
	var data []byte
	data, err = base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return
	}

	return QuantoVerifySignature(data, signature)
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
	return pgpBackend.SignData(ctx, fingerprint, data, crypto.SHA512)
}

// QuantoSignData signs the data using a already loaded and unlocked private key and returning in Quanto PGP Signature format
func QuantoSignData(data []byte, fingerprint string) (result string, err error) {
	result, err = pgpBackend.SignData(ctx, fingerprint, data, crypto.SHA512)
	if err != nil {
		return "", err
	}
	result = remote_signer.GPG2Quanto(result, fingerprint, "SHA512")

	return result, nil
}

// SignBase64Data signs data using a already loaded and unlocked private key.
// The b64data is a raw binary data encoded in base64 string
// export SignBase64Data
func SignBase64Data(b64data, fingerprint string) (result string, err error) {
	var data []byte
	data, err = base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return
	}

	return SignData(data, fingerprint)
}

// QuantoSignBase64Data signs the data using a already loaded and unlocked private key and returning in Quanto Signature format.
//  The b64data is a raw binary data encoded in base64 string
func QuantoSignBase64Data(b64data, fingerprint string) (result string, err error) {
	var data []byte
	data, err = base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return
	}

	return QuantoSignData(data, fingerprint)
}
