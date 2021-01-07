package memory

import (
	"fmt"

	"github.com/quan-to/chevron/pkg/models"
)

func (h *DbDriver) InitCursor() error {
	return fmt.Errorf("not implemented")
}

func (h *DbDriver) FinishCursor() error {
	return fmt.Errorf("not implemented")
}

func (h *DbDriver) NextGPGKey(key *models.GPGKey) bool {
	return false
}

func (h *DbDriver) NextUser(user *models.User) bool {
	return false
}

func (h *DbDriver) NumGPGKeys() (int, error) {
	return 0, fmt.Errorf("not implemented")
}

// AddGPGKey adds a list GPG Key to the database or update an existing one by fingerprint
// Same as AddGPGKey but in a single transaction
func (h *DbDriver) AddGPGKeys(keys []models.GPGKey) (ids []string, addeds []bool, err error) {
	h.log.Debug("AddGPGKeys(...%d)", len(keys))
	h.lock.Lock()
	defer h.lock.Unlock()

	for _, key := range keys {
		id, added, err := h.addGpgKey(key)
		if err != nil {
			return ids, addeds, err
		}
		ids = append(ids, id)
		addeds = append(addeds, added)
	}

	return ids, addeds, err
}
