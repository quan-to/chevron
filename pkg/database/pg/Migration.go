package pg

import (
	"github.com/quan-to/chevron/pkg/models"
)

func (h *PostgreSQLDBDriver) InitCursor() error {
	gpgRows, err := h.conn.Queryx("SELECT * FROM chevron_gpg_key")
	if err != nil {
		return err
	}

	h.gpgKeysRows = gpgRows

	userRows, err := h.conn.Queryx("SELECT * FROM chevron_user")
	if err != nil {
		_ = h.gpgKeysRows.Close()
		h.gpgKeysRows = nil
		return err
	}

	h.usersRows = userRows

	return nil
}

func (h *PostgreSQLDBDriver) FinishCursor() error {
	var gpgKeysError error
	var userError error
	if h.gpgKeysRows != nil {
		gpgKeysError = h.gpgKeysRows.Close()
		h.gpgKeysRows = nil
	}

	if h.usersRows != nil {
		userError = h.usersRows.Close()
		h.usersRows = nil
	}

	if gpgKeysError != nil {
		return gpgKeysError
	}

	return userError
}

func (h *PostgreSQLDBDriver) NextGPGKey(key *models.GPGKey) bool {
	pgKey := &pgGPGKey{}

	if h.gpgKeysRows.Next() {
		err := h.gpgKeysRows.StructScan(pgKey)
		if err != nil {
			h.log.Error("Error fetching next GPG Key: %s", err)
			return false
		}
		tx, err := h.conn.Beginx()
		if err != nil {
			h.log.Error("Error starting transaction: %s", err)
			return false
		}
		defer func() {
			_ = tx.Rollback()
		}()

		newKey, err := pgKey.toGPGKey(tx)
		if err != nil {
			h.log.Error("Error fetching data: %s", err)
			return false
		}
		key.ID = newKey.ID
		key.FullFingerprint = newKey.FullFingerprint
		key.Names = newKey.Names
		key.Emails = newKey.Emails
		key.KeyUids = newKey.KeyUids
		key.KeyBits = newKey.KeyBits
		key.Subkeys = newKey.Subkeys
		key.AsciiArmoredPublicKey = newKey.AsciiArmoredPublicKey
		key.AsciiArmoredPrivateKey = newKey.AsciiArmoredPrivateKey
		key.ParentKey = newKey.ParentKey
		return true
	}

	return false
}

func (h *PostgreSQLDBDriver) NextUser(user *models.User) bool {
	pgUser := &pgUser{}

	if h.usersRows.Next() {
		err := h.usersRows.StructScan(pgUser)
		if err != nil {
			h.log.Error("Error fetching next User: %s", err)
			return false
		}

		newUser := pgUser.toUser()

		user.ID = newUser.ID
		user.Fingerprint = newUser.Fingerprint
		user.Username = newUser.Username
		user.FullName = newUser.FullName
		user.CreatedAt = newUser.CreatedAt
		user.Password = newUser.Password

		return true
	}

	return false
}

func (h *PostgreSQLDBDriver) NumGPGKeys() (int, error) {
	count := -1
	err := h.conn.Get(&count, "SELECT COUNT(*) FROM chevron_gpg_key")
	return count, err
}

// AddGPGKey adds a list GPG Key to the database or update an existing one by fingerprint
// Same as AddGPGKey but in a single transaction
func (h *PostgreSQLDBDriver) AddGPGKeys(keys []models.GPGKey) (ids []string, addeds []bool, err error) {
	h.log.Debug("AddGPGKeys(...%d)", len(keys))
	tx, err := h.conn.Beginx()
	if err != nil {
		return nil, nil, err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	for _, key := range keys {
		id, added, err := h.addGPGKey(tx, key)
		if err != nil {
			return ids, addeds, err
		}
		ids = append(ids, id)
		addeds = append(addeds, added)
	}

	return ids, addeds, err
}
