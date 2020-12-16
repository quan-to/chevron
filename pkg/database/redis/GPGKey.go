package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/models"
)

const (
	gpgKeyByIDPrefix             = "gpgKeyByID-"
	gpgKeyByFingerprintPrefix    = "gpgKeyByFingerprint-"
	gpgKeyEntryList              = "gpgKeyEntryList-"
	gpgKeysByEmailCriteria       = "gpgKeysByEmail-"
	gpgKeysByFingerprintCriteria = "gpgKeysByFingerprint-"
	gpgKeysByValueCriteria       = "gpgKeysByValue-"
	gpgKeysByNameCriteria        = "gpgKeysByName-"

	gpgKeyExpiration        = time.Hour * 24 * 7 // One week expiration for keys
	gpgKeyEntriesExpiration = time.Minute * 15   // 15 minutes for key entry list
)

func (h *Driver) cacheKeyList(keys []models.GPGKey, keyString string) error {
	return h.cache.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   keyString,
		Value: &keys,
		TTL:   gpgKeyEntriesExpiration,
	})
}

type keyListFallbackFunc = func(string, int, int) ([]models.GPGKey, error)

func (h *Driver) getKeyListCache(value, criteria string, pageStart, pageEnd int, fallback keyListFallbackFunc) (keys []models.GPGKey, err error) {
	keyString := fmt.Sprintf("%s%s%s%d%d", gpgKeyEntryList, criteria, value, pageStart, pageEnd)
	err = h.cache.Get(context.TODO(), keyString, &keys)
	if err == nil { // Cache hit
		return keys, nil
	}

	// Get fallback
	keys, err = fallback(value, pageStart, pageEnd)
	if err != nil {
		return nil, err
	}

	if cacheErr := h.cacheKeyList(keys, keyString); cacheErr != nil {
		// Cache errors are just logged
		h.log.Error("error caching key entry list for %s: %s", criteria, cacheErr)
	}
	return keys, err
}

func (h *Driver) cacheKey(key models.GPGKey) error {
	shortFingerprint := tools.FPto16(key.FullFingerprint)
	// Check if key has an ID
	// if not fetch by using it's fingerprint
	if key.ID == "" {
		h.log.Debug("key came without id. this might be a bug. fetching key")
		// To be safe, we call the proxy directly since we're in the cache call
		newGpgKey, err := h.proxy.FetchGPGKeyByFingerprint(key.FullFingerprint)
		// We need the ID to make the cache, so if we can't have return a error
		if err != nil {
			h.log.Error("error fetching key %s for caching: %s", shortFingerprint, err)
			return err
		}
		key.ID = newGpgKey.ID
	}

	// Cache by ID
	h.log.Debug("Caching key by id %s", key.ID)
	if err := h.cache.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   gpgKeyByIDPrefix + key.ID,
		Value: &key,
		TTL:   gpgKeyExpiration,
	}); err != nil {
		h.log.Error("error caching gpg key by ID(%s): %s", key.ID, err)
		return err
	}

	// Cache by Fingerprint
	h.log.Debug("Caching key by fingerprint %s", tools.FPto16(key.FullFingerprint))
	if err := h.cache.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   gpgKeyByFingerprintPrefix + tools.FPto16(key.FullFingerprint),
		Value: &key,
		TTL:   gpgKeyExpiration,
	}); err != nil {
		h.log.Error("error caching gpg key by Fingerprint(%s): %s", shortFingerprint, err)
		return err
	}

	return nil
}

func (h *Driver) getCachedKeyById(keyId string) (key *models.GPGKey, err error) {
	err = h.cache.Get(context.TODO(), gpgKeyByIDPrefix+keyId, &key)
	return key, err
}

func (h *Driver) getCachedKeyByFingerprint(fingerprint string) (key *models.GPGKey, err error) {
	fingerprint = tools.FPto16(fingerprint)
	err = h.cache.Get(context.TODO(), gpgKeyByFingerprintPrefix+fingerprint, &key)
	return key, err
}

func (h *Driver) invalidateCachedKey(key models.GPGKey) error {
	err := h.cache.Delete(context.TODO(), gpgKeyByIDPrefix+key.ID)
	if err != nil {
		return err
	}
	return h.cache.Delete(context.TODO(), gpgKeyByFingerprintPrefix+key.FullFingerprint)
}

// UpdateGPGKey updates the specified GPG key by using it's ID
func (h *Driver) UpdateGPGKey(key models.GPGKey) (err error) {
	h.log.Debug("UpdateGPGKey(%s)", key.FullFingerprint)
	err = h.proxy.UpdateGPGKey(key)
	if err == nil {
		// The cacheKey will log the error
		// and we don't want to break the flow
		_ = h.cacheKey(key)
	}
	return err
}

// DeleteGPGKey deletes the specified GPG key by using it's ID
func (h *Driver) DeleteGPGKey(key models.GPGKey) error {
	h.log.Debug("DeleteGPGKey(%s)", key.FullFingerprint)

	if key.ID == "" {
		h.log.Debug("no key id provided for deleting. Fetching it from database using fingerprint %s", key.GetShortFingerPrint())
		existingKey, err := h.FetchGPGKeyByFingerprint(key.FullFingerprint)
		if err != nil {
			return err
		}
		key.ID = existingKey.ID
	}

	err := h.invalidateCachedKey(key)
	if err != nil {
		// Invalidating cache is critical here, so we will return the error if we can't invalidate it.
		h.log.Error("error invalidating cache for key %s(%s): %s", key.ID, key.GetShortFingerPrint(), err)
		return err
	}

	return h.proxy.DeleteGPGKey(key)
}

// AddGPGKey adds a GPG Key to the database or update an existing one by fingerprint
// Returns generated id / hasBeenAdded / error
func (h *Driver) AddGPGKey(key models.GPGKey) (string, bool, error) {
	h.log.Debug("AddGPGKey(%s)", key.FullFingerprint)
	id, added, err := h.proxy.AddGPGKey(key)
	// Set the returning ID to the input key so we cache correctly
	key.ID = id
	// The cacheKey will log the error
	// and we don't want to break the flow
	_ = h.cacheKey(key)

	return id, added, err
}

// FetchGPGKeysWithoutSubKeys fetch all keys that does not have a subkey
// This query is not implemented on PostgreSQL
func (h *Driver) FetchGPGKeysWithoutSubKeys() (res []models.GPGKey, err error) {
	h.log.Debug("FetchGPGKeysWithoutSubKeys()")
	// No caching. This is for migrating data
	return h.proxy.FetchGPGKeysWithoutSubKeys()
}

// FetchGPGKeyByFingerprint fetch a GPG Key by its fingerprint
func (h *Driver) FetchGPGKeyByFingerprint(fingerprint string) (*models.GPGKey, error) {
	h.log.Debug("FetchGPGKeyByFingerprint(%s)", fingerprint)
	key, err := h.getCachedKeyByFingerprint(fingerprint)
	if err != nil { // Cache miss
		h.log.Debug("load cache %s error: %s", fingerprint, err)
		if key, err = h.proxy.FetchGPGKeyByFingerprint(fingerprint); err == nil {
			// The cacheKey will log the error
			// and we don't want to break the flow
			_ = h.cacheKey(*key)
		}
	}

	return key, err
}

// FindGPGKeyByEmail find all keys that has a underlying UID that contains that email
func (h *Driver) FindGPGKeyByEmail(email string, pageStart, pageEnd int) (res []models.GPGKey, err error) {
	h.log.Debug("FindGPGKeyByEmail(%s, %d, %d)", email, pageStart, pageEnd)
	return h.getKeyListCache(email, gpgKeysByEmailCriteria, pageStart, pageEnd, h.proxy.FindGPGKeyByEmail)
}

// FindGPGKeyByFingerPrint find all keys that has a fingerprint that matches the specified fingerprint
func (h *Driver) FindGPGKeyByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByFingerPrint(%s, %d, %d)", fingerPrint, pageStart, pageEnd)
	return h.getKeyListCache(fingerPrint, gpgKeysByFingerprintCriteria, pageStart, pageEnd, h.proxy.FindGPGKeyByFingerPrint)
}

// FindGPGKeyByValue find all keys that has a underlying UID that contains that email, name or fingerprint specified by value
func (h *Driver) FindGPGKeyByValue(value string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByValue(%s, %d, %d)", value, pageStart, pageEnd)
	return h.getKeyListCache(value, gpgKeysByValueCriteria, pageStart, pageEnd, h.proxy.FindGPGKeyByValue)
}

// FindGPGKeyByName find all keys that has a underlying UID that contains that name
func (h *Driver) FindGPGKeyByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByName(%s, %d, %d)", name, pageStart, pageEnd)
	return h.getKeyListCache(name, gpgKeysByNameCriteria, pageStart, pageEnd, h.proxy.FindGPGKeyByName)
}
