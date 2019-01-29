package etc

import (
	"github.com/quan-to/remote-signer/models"
	"github.com/quan-to/remote-signer/openpgp"
)

type KRMInterface interface {
	GetCachedKeys() []models.KeyInfo
	ContainsKey(fp string) bool
	GetKey(fp string) *openpgp.Entity
	AddKey(key *openpgp.Entity, nonErasable bool)
	GetFingerPrints() []string
}
