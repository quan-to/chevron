package etc

import (
	"context"
	"crypto"

	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/openpgp"
	"github.com/quan-to/chevron/openpgp/packet"
)

type PGPInterface interface {
	LoadKeys(ctx context.Context)
	LoadKeyWithMetadata(ctx context.Context, armoredKey, metadata string) (error, int)
	LoadKey(ctx context.Context, armoredKey string) (error, int)
	FixFingerPrint(fp string) string
	IsKeyLocked(fp string) bool
	UnlockKey(ctx context.Context, fp, password string) error
	GetLoadedPrivateKeys(ctx context.Context) []models.KeyInfo
	GetLoadedKeys() []models.KeyInfo
	SaveKey(fingerPrint, armoredData string, password interface{}) error
	DeleteKey(ctx context.Context, fingerPrint string) error
	SignData(ctx context.Context, fingerPrint string, data []byte, hashAlgorithm crypto.Hash) (string, error)
	GetPublicKeyEntity(ctx context.Context, fingerPrint string) *openpgp.Entity
	GetPublicKey(ctx context.Context, fingerPrint string) *packet.PublicKey
	GetPublicKeyAscii(ctx context.Context, fingerPrint string) (string, error)
	GetPrivateKeyAscii(ctx context.Context, fingerPrint, password string) (string, error)
	GetPrivateKeyAsciiReencrypt(ctx context.Context, fingerPrint, currentPassword, newPassword string) (string, error)
	VerifySignatureStringData(ctx context.Context, data string, signature string) (bool, error)
	VerifySignature(ctx context.Context, data []byte, signature string) (bool, error)
	GeneratePGPKey(ctx context.Context, identifier, password string, numBits int) (string, error)
	Encrypt(ctx context.Context, filename, fingerPrint string, data []byte, dataOnly bool) (string, error)
	Decrypt(ctx context.Context, data string, dataOnly bool) (*models.GPGDecryptedData, error)
	GetCachedKeys(ctx context.Context) []models.KeyInfo
	SetKeysBase64Encoded(bool)
	MinKeyBits() int
	GenerateTestKey() (string, error)
	GetPrivate(ctx context.Context, fingerPrint string) openpgp.EntityList
	GetPrivateKeyInfo(ctx context.Context, fingerPrint string) *models.KeyInfo
}
