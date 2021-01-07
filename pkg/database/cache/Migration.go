package cache

import (
	"github.com/quan-to/chevron/pkg/models"
)

func (h *Driver) InitCursor() error {
	return h.proxy.InitCursor()
}

func (h *Driver) FinishCursor() error {
	return h.proxy.FinishCursor()
}

func (h *Driver) NextGPGKey(key *models.GPGKey) bool {
	return h.proxy.NextGPGKey(key)
}

func (h *Driver) NextUser(user *models.User) bool {
	return h.proxy.NextUser(user)
}

func (h *Driver) NumGPGKeys() (int, error) {
	return h.proxy.NumGPGKeys()
}

// AddGPGKey adds a list GPG Key to the database or update an existing one by fingerprint
// Same as AddGPGKey but in a single transaction
func (h *Driver) AddGPGKeys(keys []models.GPGKey) ([]string, []bool, error) {
	h.log.Debug("AddGPGKeys(...%d)", len(keys))
	id, added, err := h.proxy.AddGPGKeys(keys)
	// Set the returning ID to the input key so we cache correctly
	for k, v := range keys {
		v.ID = id[k]
		// The cacheKey will log the error
		// and we don't want to break the flow
		_ = h.cacheKey(v)
	}

	return id, added, err
}
