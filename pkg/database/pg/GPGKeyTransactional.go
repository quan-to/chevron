package pg

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/quan-to/chevron/pkg/models"
)

func (h *PostgreSQLDBDriver) updateGPGKey(tx *sqlx.Tx, key models.GPGKey) error {
	gpgKey, err := h.fetchGPGKeyByFingerprint(tx, key.FullFingerprint)
	if errorIsNotNilAndNotNotFound(err) {
		return err
	}

	if gpgKey != nil {
		// Preload UIDs
		_, err = gpgKey.getKeyUids(tx)
		if err != nil {
			return err
		}
		// Update UIDs
		err = gpgKey.updateUIDs(tx, key.KeyUids)
		if err != nil {
			return err
		}

		// Update Fields if needed
		if gpgKey.fieldsChanged(key) {
			gpgKey.KeyBits = key.KeyBits
			gpgKey.ASCIIArmoredPublicKey = key.AsciiArmoredPublicKey
			gpgKey.ASCIIArmoredPrivateKey = key.AsciiArmoredPrivateKey
			err := gpgKey.save(tx)
			if err != nil {
				return err
			}
		}
		return nil
	}

	gpgKey = pgGPGKeyFromGPGKey(key)
	return gpgKey.save(tx)
}

func (h *PostgreSQLDBDriver) deleteGPGKey(tx *sqlx.Tx, key models.GPGKey) error {
	pgKey := pgGPGKeyFromGPGKey(key)

	_, err := tx.NamedExec(`DELETE FROM chevron_gpg_key WHERE gpg_key_id = :gpg_key_id`, &pgKey)

	return err
}

// Transactional Methods
func (h *PostgreSQLDBDriver) addGPGKey(tx *sqlx.Tx, key models.GPGKey) (string, bool, error) {
	existingKey, err := h.fetchGPGKeyByFingerprint(tx, key.FullFingerprint)
	if errorIsNotNilAndNotNotFound(err) {
		return "", false, err
	}

	if existingKey != nil {
		key.ID = existingKey.ID
		err := h.updateGPGKey(tx, key)
		return existingKey.ID, false, err
	}

	key.ID = "" // Avoid bugs on inserting
	pgKey := pgGPGKeyFromGPGKey(key)
	err = pgKey.save(tx)
	if err != nil {
		return "", false, err
	}

	return pgKey.ID, true, nil
}

func (h *PostgreSQLDBDriver) fetchGPGKeysWithoutSubKeys(tx *sqlx.Tx) ([]pgGPGKey, error) {
	// Not supported. This would be a heavy query and it isn't used as much
	return nil, fmt.Errorf("not supported")
}

func (h *PostgreSQLDBDriver) fetchGPGKeyByFingerprint(tx *sqlx.Tx, fingerprint string) (*pgGPGKey, error) {
	if len(fingerprint) < 16 {
		return nil, fmt.Errorf("expected fingerprint length >= 16")
	}
	if len(fingerprint) > 16 {
		fingerprint = fingerprint[len(fingerprint)-16:]
	}

	key := &pgGPGKey{}

	err := tx.Get(key, "SELECT * FROM chevron_gpg_key WHERE gpg_key_fingerprint16 = $1 LIMIT 1", fingerprint)
	if err != nil && !strings.EqualFold("sql: no rows in result set", err.Error()) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("not found")
	}

	return key, nil
}

func (h *PostgreSQLDBDriver) findGPGKeyByEmail(tx *sqlx.Tx, email string, pageStart, pageEnd int) (res []pgGPGKey, err error) {
	return nil, fmt.Errorf("not supported") // Slow query
}

func (h *PostgreSQLDBDriver) findGPGKeyByFingerPrint(tx *sqlx.Tx, fingerPrint string, pageStart, pageEnd int) (res []pgGPGKey, err error) {
	if pageStart < 0 {
		pageStart = models.DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = models.DefaultPageEnd
	}

	numItems := pageEnd - pageStart

	if len(fingerPrint) > 16 {
		fingerPrint = fingerPrint[len(fingerPrint)-16:]
	}

	fingerPrint = strings.ToUpper("%" + fingerPrint)

	err = tx.Select(&res, "SELECT * FROM chevron_gpg_key WHERE gpg_key_fingerprint16 LIKE $1 LIMIT $2 OFFSET $3", fingerPrint, numItems, pageStart)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *PostgreSQLDBDriver) findGPGKeyByValue(tx *sqlx.Tx, value string, pageStart, pageEnd int) (res []pgGPGKey, err error) {
	return nil, fmt.Errorf("not supported") // Slow query
}

func (h *PostgreSQLDBDriver) findGPGKeyByName(tx *sqlx.Tx, name string, pageStart, pageEnd int) (res []pgGPGKey, err error) {
	return nil, fmt.Errorf("not supported") // Slow query
}
