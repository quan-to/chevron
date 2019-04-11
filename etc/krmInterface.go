package etc

import (
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/openpgp"
)

type KRMInterface interface {
	GetCachedKeys() []models.KeyInfo
	ContainsKey(fp string) bool
	GetKey(fp string) *openpgp.Entity
	AddKey(key *openpgp.Entity, nonErasable bool)
	GetFingerPrints() []string
}
