package rql

import (
	"fmt"

	"github.com/quan-to/chevron/pkg/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func (h *RethinkDBDriver) InitCursor() error {
	res, err := r.Table(gpgKeyTableInit.TableName).Run(h.conn)
	if err != nil {
		return err
	}

	h.gpgKeysMigrationCursor = res

	res, err = r.Table(userModelTableInit.TableName).Run(h.conn)
	if err != nil {
		_ = h.gpgKeysMigrationCursor.Close()
		h.gpgKeysMigrationCursor = nil
		return err
	}

	h.userMigrationCursor = res
	return nil
}

func (h *RethinkDBDriver) FinishCursor() error {
	var gpgKeysError error
	var userError error
	if h.gpgKeysMigrationCursor != nil {
		gpgKeysError = h.gpgKeysMigrationCursor.Close()
		h.gpgKeysMigrationCursor = nil
	}

	if h.userMigrationCursor != nil {
		userError = h.userMigrationCursor.Close()
		h.userMigrationCursor = nil
	}

	if gpgKeysError != nil {
		return gpgKeysError
	}

	return userError
}

func (h *RethinkDBDriver) NextGPGKey(key *models.GPGKey) bool {
	var gpgKey map[string]interface{}

	if h.gpgKeysMigrationCursor.Next(&gpgKey) {
		gpgKey = h.fixGPGKey(gpgKey)
		err := convertFromRethinkDB(gpgKey, key)
		if err != nil {
			h.log.Error("Error fetching next GPG Key: %s", err)
			return false
		}
		return true
	}

	return false
}

func (h *RethinkDBDriver) NextUser(user *models.User) bool {
	rdata := map[string]interface{}{}

	if h.userMigrationCursor.Next(&rdata) {
		err := convertFromRethinkDB(rdata, user)
		if err != nil {
			h.log.Error("Error fetching next User: %s", err)
			return false
		}
		return true
	}

	return false
}

func (h *RethinkDBDriver) NumGPGKeys() (int, error) {
	res, err := r.Table(gpgKeyTableInit.TableName).Count().Run(h.conn)
	if err != nil {
		return -1, err
	}

	count := -1

	if res.Next(&count) {
		return count, nil
	}

	return -1, fmt.Errorf("no table to count from")
}

// AddGPGKey adds a list GPG Key to the database or update an existing one by fingerprint
// Same as AddGPGKey but in a single transaction
func (h *RethinkDBDriver) AddGPGKeys(keys []models.GPGKey) (ids []string, addeds []bool, err error) {
	// TODO: Proper parallel add if we ever need
	for _, key := range keys {
		id, added, err := h.AddGPGKey(key)
		if err != nil {
			return ids, addeds, err
		}
		ids = append(ids, id)
		addeds = append(addeds, added)
	}

	return ids, addeds, err
}
