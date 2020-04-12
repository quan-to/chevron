package interfaces

import (
	"context"
	"github.com/quan-to/chevron/internal/models"

	"github.com/quan-to/chevron/pkg/openpgp"
)

// KeyRingManager is an interface to a Key Ring Manager Instance
type KeyRingManager interface {
	// GetCachedKeys returns a list of the memory-cached keys
	GetCachedKeys(ctx context.Context) []models.KeyInfo
	// ContainsKey checks if a key with the specified fingerprint exists in Key Ring
	ContainsKey(ctx context.Context, fingerprint string) bool
	// GetKey returns a key with the specified fingerprint if exists. Returns nil if it does not
	GetKey(ctx context.Context, fingerprint string) *openpgp.Entity
	// AddKey adds a key to key ring manager. If nonErasable is true it will be persistent in cache
	AddKey(ctx context.Context, key *openpgp.Entity, nonErasable bool)
	// GetFingerprints returns a list of stored key fingerpints
	GetFingerPrints(ctx context.Context) []string
	// DeleteKey erases the specified key from the key ring
	DeleteKey(ctx context.Context, fingerprint string) error
}
