package etc

import (
	"context"

	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/openpgp"
)

type KRMInterface interface {
	GetCachedKeys(ctx context.Context) []models.KeyInfo
	ContainsKey(ctx context.Context, fp string) bool
	GetKey(ctx context.Context, fp string) *openpgp.Entity
	AddKey(ctx context.Context, key *openpgp.Entity, nonErasable bool)
	GetFingerPrints(ctx context.Context) []string
	DeleteKey(ctx context.Context, fp string) error
}
