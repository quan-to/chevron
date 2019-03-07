package etc

import (
	"crypto"
	"github.com/quan-to/remote-signer/models"
	"github.com/quan-to/remote-signer/openpgp"
	"github.com/quan-to/remote-signer/openpgp/packet"
)

type PGPInterface interface {
	LoadKeys()
	LoadKeyWithMetadata(armoredKey, metadata string) (error, int)
	LoadKey(armoredKey string) (error, int)
	FixFingerPrint(fp string) string
	IsKeyLocked(fp string) bool
	UnlockKey(fp, password string) error
	GetLoadedPrivateKeys() []models.KeyInfo
	SavePrivateKey(fingerPrint, armoredData string, password interface{}) error
	SignData(fingerPrint string, data []byte, hashAlgorithm crypto.Hash) (string, error)
	GetPublicKeyEntity(fingerPrint string) *openpgp.Entity
	GetPublicKey(fingerPrint string) *packet.PublicKey
	GetPublicKeyAscii(fingerPrint string) (string, error)
	VerifySignatureStringData(data string, signature string) (bool, error)
	VerifySignature(data []byte, signature string) (bool, error)
	GeneratePGPKey(identifier, password string, numBits int) (string, error)
	Encrypt(filename, fingerPrint string, data []byte, dataOnly bool) (string, error)
	Decrypt(data string, dataOnly bool) (*models.GPGDecryptedData, error)
	GetCachedKeys() []models.KeyInfo
	SetKeysBase64Encoded(bool)
	MinKeyBits() int
	GenerateTestKey() (string, error)
	GetPrivate(fingerPrint string) openpgp.EntityList
}
