package chevronlib

import (
	"fmt"
	"os"

	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/openpgp"
)

// GenerateKey generates a new key using specified bits and identifier and encrypts it using the specified password
func GenerateKey(password, identifier string, bits int) (result string, err error) {
	if password == "" {
		err = fmt.Errorf("no password supplied")
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, "Generating key. This might take a while...")

	result, err = pgpBackend.GeneratePGPKey(ctx, identifier, password, bits)

	return
}

// GetKeyFingerprints returns all fingerprints in a ASCII Armored PGP Keychain
func GetKeyFingerprints(keyData string) (fps []string, err error) {
	return tools.GetFingerPrintsFromKey(keyData)
}

// ChangeKeyPassword re-encrypts the input key using newPassword
func ChangeKeyPassword(keyData, currentPassword, newPassword string) (newKeyData string, err error) {
	n, e := pgpBackend.LoadKey(ctx, keyData)

	if e != nil {
		err = e
		return
	}

	if n == 0 {
		err = fmt.Errorf("no private key")
		return
	}

	fp, _ := tools.GetFingerPrintFromKey(keyData)

	err = pgpBackend.UnlockKey(ctx, fp, currentPassword)
	if err != nil {
		return
	}
	newKeyData, err = pgpBackend.GetPrivateKeyASCIIReencrypt(ctx, fp, currentPassword, newPassword)

	_ = pgpBackend.DeleteKey(ctx, fp) // Clean key after changing password
	return
}

// GetPublicKey returns the cached public key from the specified fingerprint
func GetPublicKey(fingerprint string) (keyData string, err error) {
	return pgpBackend.GetPublicKeyASCII(ctx, fingerprint)
}

// GetPublicKeyEntity returns the public key entity for the specified fingerprint
func GetPublicKeyEntity(fingerprint string) *openpgp.Entity {
	return pgpBackend.GetPublicKeyEntity(ctx, fingerprint)
}

// GetPrivateKeyEntity returns the private key entity list for the specified fingerprint
func GetPrivateKeyEntity(fingerprint string) openpgp.EntityList {
	return pgpBackend.GetPrivate(ctx, fingerprint)
}

// GetPGPManager returns the PGP Manager instance of Chevron Lib
// Use with care
func GetPGPManager() interfaces.PGPManager {
	return pgpBackend
}
