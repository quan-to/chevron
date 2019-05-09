package chevronlib

import (
	"fmt"
	"github.com/quan-to/chevron"
	"os"
)

// GenerateKey generates a new key using specified bits and identifier and encrypts it using the specified password
func GenerateKey(password, identifier string, bits int) (result string, err error) {
	if password == "" {
		err = fmt.Errorf("no password supplied")
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, "Generating key. This might take a while...")

	result, err = pgpBackend.GeneratePGPKey(identifier, password, bits)

	return
}

// GetKeyFingerprints returns all fingerprints in a ASCII Armored PGP Keychain
func GetKeyFingerprints(keyData string) (fps []string, err error) {
	return remote_signer.GetFingerPrintsFromKey(keyData)
}

// ChangeKeyPassword re-encrypts the input key using newPassword
func ChangeKeyPassword(keyData, currentPassword, newPassword string) (newKeyData string, err error) {
	e, n := pgpBackend.LoadKey(keyData)

	if e != nil {
		err = e
		return
	}

	if n == 0 {
		err = fmt.Errorf("no private key")
		return
	}

	fp, _ := remote_signer.GetFingerPrintFromKey(keyData)

	err = pgpBackend.UnlockKey(fp, currentPassword)
	if err != nil {
		return
	}

	newKeyData, err = pgpBackend.GetPrivateKeyAscii(fp, newPassword)
	return
}

// GetPublicKey returns the cached public key from the specified fingerprint
func GetPublicKey(fingerprint string) (keyData string, err error) {
	return pgpBackend.GetPublicKeyAscii(fingerprint)
}
