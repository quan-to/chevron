package interfaces

import (
	"context"
	"crypto"
	"github.com/quan-to/chevron/internal/models"

	"github.com/quan-to/chevron/pkg/openpgp"
	"github.com/quan-to/chevron/pkg/openpgp/packet"
)

// PGPManager is a interface for handling PGP Operations
type PGPManager interface {
	// LoadKeys loads the keys stored on the PGP Manager key backend
	LoadKeys(ctx context.Context)
	// LoadKeyWithMetadata loads a armored ascii key with the specified json metadata
	LoadKeyWithMetadata(ctx context.Context, armoredKey, metadata string) (int, error)
	// LoadKey loads a armored ascii key
	LoadKey(ctx context.Context, armoredKey string) (int, error)
	// FixFingerPrint fixes and trims the fingerprint to 16 Char Hex
	FixFingerPrint(fingerprint string) string
	// IsKeyLocked returns if the specified key is currently locked inside the PGP Manager
	IsKeyLocked(fingerprint string) bool
	// UnlockKey unlocks the specified key with the specified password
	UnlockKey(ctx context.Context, fingerprint, password string) error
	// GetLoadedPrivateKeys returns the information of each loaded private key
	GetLoadedPrivateKeys(ctx context.Context) []models.KeyInfo
	// GetLoadedKeys returns the information for all keys in PGP Manager
	GetLoadedKeys() []models.KeyInfo
	// SaveKey saves the specified key in PGP Manager Key Backend
	SaveKey(fingerprint, armoredData string, password interface{}) error
	// DeleteKey removes the specified key from the memory and key backend
	DeleteKey(ctx context.Context, fingerprint string) error
	// SignData signs the specified data with a unlocked private key
	SignData(ctx context.Context, fingerprint string, data []byte, hashAlgorithm crypto.Hash) (string, error)
	// GetPublicKeyEntity returns the public key entity
	GetPublicKeyEntity(ctx context.Context, fingerprint string) *openpgp.Entity
	// GetPublicKey returns the public key
	GetPublicKey(ctx context.Context, fingerprint string) *packet.PublicKey
	// GetPublicKeyASCII returns the public key in ASCII Armored format
	GetPublicKeyASCII(ctx context.Context, fingerprint string) (string, error)
	// GetPublicKeyASCII returns the encrypted private key in ASCII Armored format
	GetPrivateKeyASCII(ctx context.Context, fingerprint, password string) (string, error)
	// GetPublicKeyASCII returns the encrypted private key in ASCII Armored format changing it's password
	GetPrivateKeyASCIIReencrypt(ctx context.Context, fingerprint, currentPassword, newPassword string) (string, error)
	// VerifySignatureStringData verifies signature of specified data in string format
	VerifySignatureStringData(ctx context.Context, data string, signature string) (bool, error)
	// VerifySignatureStringData verifies signature of specified data
	VerifySignature(ctx context.Context, data []byte, signature string) (bool, error)
	// GeneratePGPKey generates a new PGP Key with the specified information
	GeneratePGPKey(ctx context.Context, identifier, password string, numBits int) (string, error)
	// Encrypt encrypts data using the specified public key.
	// Filename is a metadata from GPG
	// dataOnly field specifies that it will encrypt as binary content instead ASCII Armored
	Encrypt(ctx context.Context, filename, fingerprint string, data []byte, dataOnly bool) (string, error)
	// Decrypt decrypts data using any available unlocked private key
	Decrypt(ctx context.Context, data string, dataOnly bool) (*models.GPGDecryptedData, error)
	// GetCachedKeys returns all cached public keys in memory
	GetCachedKeys(ctx context.Context) []models.KeyInfo
	// SetKeysBase64Encoded sets if keys should be stored in Base64 Encoded format
	SetKeysBase64Encoded(bool)
	// MinKeyBits returns the minimum key bits allowed for generating PGP Keys
	MinKeyBits() int
	// GenerateTestKey generates a private key for testing
	// Bits: MinKeyBits
	// Password: 1234
	// Identity: *empty string*
	GenerateTestKey() (string, error)
	// GetPrivate returns the private key entity list for a specified private key
	GetPrivate(ctx context.Context, fingerprint string) openpgp.EntityList
	// GetPrivateKeyInfo returns the information of the specified private key
	GetPrivateKeyInfo(ctx context.Context, fingerprint string) *models.KeyInfo
}
