package pg

import (
	"github.com/jmoiron/sqlx"
	"github.com/quan-to/chevron/pkg/models"
)

// database calls

// UpdateGPGKey updates the specified GPG key by using it's ID
func (h *PostgreSQLDBDriver) UpdateGPGKey(key models.GPGKey) (err error) {
	h.log.Debug("UpdateGPGKey(%s)", key.FullFingerprint)
	tx, err := h.conn.Beginx()
	if err != nil {
		return err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	return h.updateGPGKey(tx, key)
}

// DeleteGPGKey deletes the specified GPG key by using it's ID
func (h *PostgreSQLDBDriver) DeleteGPGKey(key models.GPGKey) error {
	h.log.Debug("DeleteGPGKey(%s)", key.FullFingerprint)
	tx, err := h.conn.Beginx()
	if err != nil {
		return err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	return h.deleteGPGKey(tx, key)
}

// AddGPGKey adds a GPG Key to the database or update an existing one by fingerprint
// Returns generated id / hasBeenAdded / error
func (h *PostgreSQLDBDriver) AddGPGKey(key models.GPGKey) (string, bool, error) {
	h.log.Debug("AddGPGKey(%s)", key.FullFingerprint)
	tx, err := h.conn.Beginx()
	if err != nil {
		return "", false, err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	return h.addGPGKey(tx, key)
}

// FetchGPGKeysWithoutSubKeys fetch all keys that does not have a subkey
// This query is not implemented on PostgreSQL
func (h *PostgreSQLDBDriver) FetchGPGKeysWithoutSubKeys() (res []models.GPGKey, err error) {
	h.log.Debug("FetchGPGKeysWithoutSubKeys()")
	tx, err := h.conn.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	keys, err := h.fetchGPGKeysWithoutSubKeys(tx)
	if err != nil {
		return nil, err
	}

	return convertArray(keys, tx)
}

// FetchGPGKeyByFingerprint fetch a GPG Key by its fingerprint
func (h *PostgreSQLDBDriver) FetchGPGKeyByFingerprint(fingerprint string) (*models.GPGKey, error) {
	h.log.Debug("FetchGPGKeyByFingerprint(%s)", fingerprint)
	tx, err := h.conn.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	k, err := h.fetchGPGKeyByFingerprint(tx, fingerprint)
	if err != nil {
		return nil, err
	}

	return k.toGPGKey(tx)
}

// FindGPGKeyByEmail find all keys that has a underlying UID that contains that email
func (h *PostgreSQLDBDriver) FindGPGKeyByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByEmail(%s, %d, %d)", email, pageStart, pageEnd)
	tx, err := h.conn.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	keys, err := h.findGPGKeyByEmail(tx, email, pageStart, pageEnd)
	if err != nil {
		return nil, err
	}

	return convertArray(keys, tx)
}

// FindGPGKeyByFingerPrint find all keys that has a fingerprint that matches the specified fingerprint
func (h *PostgreSQLDBDriver) FindGPGKeyByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByFingerPrint(%s, %d, %d)", fingerPrint, pageStart, pageEnd)
	tx, err := h.conn.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	keys, err := h.findGPGKeyByFingerPrint(tx, fingerPrint, pageStart, pageEnd)
	if err != nil {
		return nil, err
	}

	return convertArray(keys, tx)
}

// FindGPGKeyByValue find all keys that has a underlying UID that contains that email, name or fingerprint specified by value
func (h *PostgreSQLDBDriver) FindGPGKeyByValue(value string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByValue(%s, %d, %d)", value, pageStart, pageEnd)
	tx, err := h.conn.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	keys, err := h.findGPGKeyByValue(tx, value, pageStart, pageEnd)
	if err != nil {
		return nil, err
	}

	return convertArray(keys, tx)
}

// FindGPGKeyByName find all keys that has a underlying UID that contains that name
func (h *PostgreSQLDBDriver) FindGPGKeyByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	h.log.Debug("FindGPGKeyByName(%s, %d, %d)", name, pageStart, pageEnd)
	tx, err := h.conn.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	keys, err := h.findGPGKeyByName(tx, name, pageStart, pageEnd)
	if err != nil {
		return nil, err
	}

	return convertArray(keys, tx)
}

func convertArray(keys []pgGPGKey, tx *sqlx.Tx) (res []models.GPGKey, err error) {
	for _, v := range keys {
		k, err := v.toGPGKey(tx)
		if err != nil {
			return nil, err
		}
		res = append(res, *k)
	}
	return res, nil
}
