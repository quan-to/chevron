package redis

import (
	"github.com/quan-to/chevron/pkg/models"
)

// UpdateGPGKey updates the specified GPG key by using it's ID
func (h *Driver) UpdateGPGKey(key models.GPGKey) (err error) {
	h.log.Debug("UpdateGPGKey(%s)", key.FullFingerprint)
	return h.proxy.UpdateGPGKey(key)
}

// DeleteGPGKey deletes the specified GPG key by using it's ID
func (h *Driver) DeleteGPGKey(key models.GPGKey) error {
	h.log.Debug("DeleteGPGKey(%s)", key.FullFingerprint)
	return h.proxy.DeleteGPGKey(key)
}

// AddGPGKey adds a GPG Key to the database or update an existing one by fingerprint
// Returns generated id / hasBeenAdded / error
func (h *Driver) AddGPGKey(key models.GPGKey) (string, bool, error) {
	h.log.Debug("AddGPGKey(%s)", key.FullFingerprint)
	return h.proxy.AddGPGKey(key)
}

// FetchGPGKeysWithoutSubKeys fetch all keys that does not have a subkey
// This query is not implemented on PostgreSQL
func (h *Driver) FetchGPGKeysWithoutSubKeys() (res []models.GPGKey, err error) {
	h.log.Debug("FetchGPGKeysWithoutSubKeys()")
	return h.proxy.FetchGPGKeysWithoutSubKeys()
}

// FetchGPGKeyByFingerprint fetch a GPG Key by its fingerprint
func (h *Driver) FetchGPGKeyByFingerprint(fingerprint string) (*models.GPGKey, error) {
	h.log.Debug("FetchGPGKeyByFingerprint(%s)", fingerprint)
	return h.proxy.FetchGPGKeyByFingerprint(fingerprint)
}

// FindGPGKeyByEmail find all keys that has a underlying UID that contains that email
func (h *Driver) FindGPGKeyByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByEmail(%s, %d, %d)", email, pageStart, pageEnd)
	return h.proxy.FindGPGKeyByEmail(email, pageStart, pageEnd)
}

// FindGPGKeyByFingerPrint find all keys that has a fingerprint that matches the specified fingerprint
func (h *Driver) FindGPGKeyByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByFingerPrint(%s, %d, %d)", fingerPrint, pageStart, pageEnd)
	return h.proxy.FindGPGKeyByFingerPrint(fingerPrint, pageStart, pageEnd)
}

// FindGPGKeyByValue find all keys that has a underlying UID that contains that email, name or fingerprint specified by value
func (h *Driver) FindGPGKeyByValue(value string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByValue(%s, %d, %d)", value, pageStart, pageEnd)
	return h.proxy.FindGPGKeyByValue(value, pageStart, pageEnd)
}

// FindGPGKeyByName find all keys that has a underlying UID that contains that name
func (h *Driver) FindGPGKeyByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByName(%s, %d, %d)", name, pageStart, pageEnd)
	return h.proxy.FindGPGKeyByName(name, pageStart, pageEnd)
}
