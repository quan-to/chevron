package memory

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/quan-to/chevron/pkg/models"
	"strings"
)

func (h *MemoryDBDriver) UpdateGPGKey(key models.GPGKey) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.log.Debug("UpdateGPGKey(%q)", key.FullFingerprint)

	if key.FullFingerprint == "" {
		return fmt.Errorf("not found")
	}

	for i, v := range h.keys {
		if strings.EqualFold(v.FullFingerprint, key.FullFingerprint) {
			h.keys[i] = key
			return nil
		}
	}

	return fmt.Errorf("not found")
}

func (h *MemoryDBDriver) DeleteGPGKey(key models.GPGKey) error {
	h.log.Debug("DeleteGPGKey(%q, %q)", key.ID, key.FullFingerprint)
	h.lock.Lock()
	defer h.lock.Unlock()

	for i, v := range h.keys {
		if v.ID == key.ID || v.FullFingerprint == key.FullFingerprint {
			h.keys = append(h.keys[:i], h.keys[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("not found")
}

// AddGPGKey adds a GPG Key to the database or update an existing one by fingerprint
// Returns generated id / hasBeenAdded / error
func (h *MemoryDBDriver) AddGPGKey(key models.GPGKey) (string, bool, error) {
	h.log.Debug("AddGPGKey(%q)", key.FullFingerprint)
	h.lock.Lock()
	defer h.lock.Unlock()

	if key.FullFingerprint == "" {
		return "", false, fmt.Errorf("invalid key fingerprint")
	}

	for i, v := range h.keys {
		if strings.EqualFold(v.FullFingerprint, key.FullFingerprint) {
			key.ID = v.ID
			h.keys[i] = key
			return key.ID, false, nil
		}
	}

	key.ID = uuid.New().String()

	h.keys = append(h.keys, key)

	return key.ID, true, nil
}

func (h *MemoryDBDriver) FetchGPGKeysWithoutSubKeys() ([]models.GPGKey, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	var keys []models.GPGKey

	for _, v := range h.keys {
		if len(v.Subkeys) == 0 {
			keys = append(keys, v)
		}
	}

	return keys, nil
}

func (h *MemoryDBDriver) FetchGPGKeyByFingerprint(fingerprint string) (*models.GPGKey, error) {
	h.log.Debug("FetchGPGKeyByFingerprint(%q)", fingerprint)
	h.lock.RLock()
	defer h.lock.RUnlock()

	if fingerprint == "" {
		return nil, fmt.Errorf("not found")
	}

	for _, v := range h.keys {
		if strings.EqualFold(v.FullFingerprint[len(v.FullFingerprint)-len(fingerprint):], fingerprint) {
			h := v
			return &h, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (h *MemoryDBDriver) FindGPGKeyByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByEmail(%q, %d, %d)", email, pageStart, pageEnd)
	h.lock.RLock()
	defer h.lock.RUnlock()

	if pageStart < 0 {
		pageStart = models.DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = models.DefaultPageEnd
	}

	var items []models.GPGKey
	delta := pageEnd - pageStart

	for _, v := range h.keys {
		for _, e := range v.Emails {
			if strings.EqualFold(e, email) {
				items = append(items, v)
				break
			}
		}
		if len(items) > delta {
			break
		}
	}

	if len(items) > pageStart {
		items = items[pageStart:]
	}

	if len(items) > delta {
		items = items[:delta]
	}

	return items, nil
}

func (h *MemoryDBDriver) FindGPGKeyByFingerPrint(fingerprint string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByFingerPrint(%q, %d, %d)", fingerprint, pageStart, pageEnd)
	h.lock.RLock()
	defer h.lock.RUnlock()

	if pageStart < 0 {
		pageStart = models.DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = models.DefaultPageEnd
	}

	var items []models.GPGKey
	delta := pageEnd - pageStart

	for _, v := range h.keys {
		if strings.EqualFold(v.FullFingerprint[len(v.FullFingerprint)-len(fingerprint):], fingerprint) {
			items = append(items, v)
		}
		if len(items) > delta {
			break
		}
	}

	if len(items) > pageStart {
		items = items[pageStart:]
	}

	if len(items) > delta {
		items = items[:delta]
	}

	return items, nil
}

func (h *MemoryDBDriver) FindGPGKeyByValue(value string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByValue(%q, %d, %d)", value, pageStart, pageEnd)
	h.lock.RLock()
	defer h.lock.RUnlock()

	if pageStart < 0 {
		pageStart = models.DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = models.DefaultPageEnd
	}

	var items []models.GPGKey
	delta := pageEnd - pageStart

	for _, v := range h.keys {
		foundEmail := false
		foundName := false
		for _, e := range v.Emails {
			if strings.Contains(e, value) {
				foundEmail = true
				break
			}
		}
		for _, e := range v.Names {
			if strings.Contains(e, value) {
				foundName = true
				break
			}
		}
		if foundEmail || foundName || strings.Contains(v.FullFingerprint, value) {
			items = append(items, v)
		}

		if len(items) > delta {
			break
		}
	}

	if len(items) > pageStart {
		items = items[pageStart:]
	}

	if len(items) > delta {
		items = items[:delta]
	}

	return items, nil
}

func (h *MemoryDBDriver) FindGPGKeyByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByName(%q, %d, %d)", name, pageStart, pageEnd)
	h.lock.RLock()
	defer h.lock.RUnlock()

	if pageStart < 0 {
		pageStart = models.DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = models.DefaultPageEnd
	}

	var items []models.GPGKey
	delta := pageEnd - pageStart

	for _, v := range h.keys {
		for _, e := range v.Names {
			if strings.EqualFold(e, name) {
				items = append(items, v)
				break
			}
		}

		if len(items) > delta {
			break
		}
	}

	if len(items) > pageStart {
		items = items[pageStart:]
	}

	if len(items) > delta {
		items = items[:delta]
	}

	return items, nil
}
