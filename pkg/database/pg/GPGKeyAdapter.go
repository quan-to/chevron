package pg

import (
	"strings"
	"time"

	"github.com/quan-to/chevron/internal/tools"

	"github.com/quan-to/chevron/pkg/uuid"

	"github.com/jmoiron/sqlx"
	"github.com/quan-to/chevron/pkg/models"
)

type pgGPGKey struct {
	ID                     string    `db:"gpg_key_id"`
	FullFingerprint        string    `db:"gpg_key_full_fingerprint"`
	Fingerprint16          string    `db:"gpg_key_fingerprint16"`
	KeyBits                int       `db:"gpg_key_keybits"`
	ASCIIArmoredPublicKey  string    `db:"gpg_key_public_key"`
	ASCIIArmoredPrivateKey string    `db:"gpg_key_private_key"`
	CreatedAt              time.Time `db:"gpg_key_created_at"`
	UpdatedAt              time.Time `db:"gpg_key_updated_at"`
	DeletedAt              time.Time `db:"gpg_key_deleted_at"`
	ParentKeyID            *string   `db:"gpg_key_parent"`

	// Relations
	keyUids       []*pgGPGKeyUID
	subkeys       []*pgGPGKey
	keyUidsLoaded bool
	subkeysLoaded bool
}

func (k *pgGPGKey) getNames(tx *sqlx.Tx) (names []string, err error) {
	uids, err := k.getKeyUids(tx)
	if err != nil {
		return names, err
	}

	for _, v := range uids {
		names = append(names, v.Name)
	}

	return names, err
}

func (k *pgGPGKey) getEmails(tx *sqlx.Tx) (emails []string, err error) {
	uids, err := k.getKeyUids(tx)
	if err != nil {
		return emails, err
	}

	for _, v := range uids {
		emails = append(emails, v.Email)
	}

	return emails, err
}

func (k *pgGPGKey) getKeyUids(tx *sqlx.Tx) ([]*pgGPGKeyUID, error) {
	if k.keyUidsLoaded {
		return k.keyUids, nil
	}

	err := tx.Select(&k.keyUids, "SELECT * FROM chevron_gpg_key_uid WHERE gpg_key_uid_parent = $1", k.ID)
	if err != nil {
		return nil, err
	}
	k.keyUidsLoaded = true

	return k.keyUids, nil
}

func (k *pgGPGKey) loadSubkeys(tx *sqlx.Tx) error {
	err := tx.Select(&k.subkeys, "SELECT * FROM chevron_gpg_key WHERE gpg_key_parent = $1", k.ID)
	if err != nil {
		return err
	}
	k.subkeysLoaded = true

	return nil
}

func pgGPGKeyFromGPGKey(key models.GPGKey) *pgGPGKey {
	k := &pgGPGKey{
		ID:                     key.ID,
		FullFingerprint:        key.FullFingerprint,
		Fingerprint16:          tools.FPto16(key.FullFingerprint),
		KeyBits:                key.KeyBits,
		ASCIIArmoredPublicKey:  key.AsciiArmoredPublicKey,
		ASCIIArmoredPrivateKey: key.AsciiArmoredPrivateKey,

		// Relations
		keyUidsLoaded: true,
		subkeysLoaded: true,
	}

	if key.ParentKey != nil {
		k.ParentKeyID = &key.ParentKey.ID
	}

	for _, uid := range key.KeyUids {
		k.keyUids = append(k.keyUids, &pgGPGKeyUID{
			Parent:      key.ID,
			Name:        uid.Name,
			Description: uid.Description,
			Email:       uid.Email,
		})
	}

	for _, fp := range key.Subkeys {
		k.subkeys = append(k.subkeys, &pgGPGKey{
			FullFingerprint: fp,
			Fingerprint16:   tools.FPto16(fp),
			keyUids:         k.keyUids,
		})
	}

	return k
}

func (k *pgGPGKey) getSubkeysFingerprints(tx *sqlx.Tx) (fps []string, err error) {
	if !k.subkeysLoaded {
		err := k.loadSubkeys(tx)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range k.subkeys {
		fps = append(fps, v.FullFingerprint)
	}

	return fps, err
}

func (k *pgGPGKey) toGPGKey(tx *sqlx.Tx) (*models.GPGKey, error) {
	_, err := k.getKeyUids(tx)
	if err != nil {
		return nil, err
	}

	subkeys, err := k.getSubkeysFingerprints(tx)
	if err != nil {
		return nil, err
	}

	var keyUids []models.GPGKeyUid

	names, _ := k.getNames(tx)
	emails, _ := k.getEmails(tx)

	for _, v := range k.keyUids {
		keyUids = append(keyUids, models.GPGKeyUid{
			Name:        v.Name,
			Email:       v.Email,
			Description: v.Description,
		})
	}

	key := models.GPGKey{
		ID:                     k.ID,
		FullFingerprint:        k.FullFingerprint,
		Names:                  names,
		Emails:                 emails,
		KeyUids:                keyUids,
		KeyBits:                k.KeyBits,
		Subkeys:                subkeys,
		AsciiArmoredPublicKey:  k.ASCIIArmoredPublicKey,
		AsciiArmoredPrivateKey: k.ASCIIArmoredPrivateKey,
	}

	return &key, nil
}

func (k *pgGPGKey) fieldsChanged(key models.GPGKey) bool {
	if strings.EqualFold(k.ASCIIArmoredPrivateKey, key.AsciiArmoredPrivateKey) ||
		strings.EqualFold(k.ASCIIArmoredPublicKey, key.AsciiArmoredPublicKey) ||
		k.KeyBits != key.KeyBits {
		return true
	}

	if k.subkeysLoaded {
		for _, v := range k.subkeys {
			found := false
			for _, subfp := range key.Subkeys {
				if tools.CompareFingerPrint(v.FullFingerprint, subfp) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return false
}

func (k *pgGPGKey) updateUIDs(tx *sqlx.Tx, newUIDs []models.GPGKeyUid) error {
	oldUids, err := k.getKeyUids(tx)
	if err != nil {
		return err
	}

	for _, newUID := range newUIDs {
		found := false
		for _, oldUID := range oldUids {
			if oldUID.compareWithKeyUID(newUID) {
				// Found
				found = true
				if oldUID.fieldsChanged(newUID) {
					oldUID.Name = newUID.Name
					oldUID.Email = newUID.Email
					oldUID.Description = newUID.Description
					err = oldUID.save(tx)
					if err != nil {
						return err
					}
				}
				break
			}
		}
		if !found {
			newUid := &pgGPGKeyUID{
				Parent:      k.ID,
				Name:        newUID.Name,
				Description: newUID.Description,
				Email:       newUID.Email,
			}
			err = newUid.save(tx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (k *pgGPGKey) save(tx *sqlx.Tx) error {
	if k.ID == "" { // Insert
		k.ID = uuid.EnsureUUID(nil)
		_, err := tx.NamedExec(`INSERT INTO 
    		chevron_gpg_key(gpg_key_id, gpg_key_full_fingerprint, gpg_key_fingerprint16, gpg_key_keybits, gpg_key_parent, gpg_key_public_key, gpg_key_private_key) 
    		VALUES (:gpg_key_id, :gpg_key_full_fingerprint, :gpg_key_fingerprint16, :gpg_key_keybits, :gpg_key_parent, :gpg_key_public_key, :gpg_key_private_key)`, k)
		if err != nil {
			return err
		}

		for _, v := range k.keyUids {
			v.Parent = k.ID
			err = v.save(tx)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Update
	_, err := tx.NamedExec(`UPDATE chevron_gpg_key SET 
                           gpg_key_private_key = :gpg_key_private_key,
                           gpg_key_public_key = :gpg_key_public_key 
                           WHERE gpg_key_id = :gpg_key_id`, k)
	return err
}

type pgGPGKeyUID struct {
	ID          string    `db:"gpg_key_uid_id"`
	Parent      string    `db:"gpg_key_uid_parent"`
	Name        string    `db:"gpg_key_uid_name"`
	Email       string    `db:"gpg_key_uid_email"`
	Description string    `db:"gpg_key_uid_description"`
	CreatedAt   time.Time `db:"gpg_key_uid_created_at"`
	UpdatedAt   time.Time `db:"gpg_key_uid_updated_at"`
	DeletedAt   time.Time `db:"gpg_key_uid_deleted_at"`
}

func (k *pgGPGKeyUID) fieldsChanged(m models.GPGKeyUid) bool {
	return k.Name != m.Name || k.Email != m.Email || k.Description != m.Description
}

func (k *pgGPGKeyUID) compareWithKeyUID(m models.GPGKeyUid) bool {
	return strings.EqualFold(k.Name, m.Name) || strings.EqualFold(k.Email, m.Email)
}

func (k *pgGPGKeyUID) save(tx *sqlx.Tx) error {
	if k.ID == "" { // Insert
		k.ID = uuid.EnsureUUID(nil)
		_, err := tx.NamedExec(`INSERT INTO
                               chevron_gpg_key_uid(gpg_key_uid_id, gpg_key_uid_name, gpg_key_uid_email, gpg_key_uid_description, gpg_key_uid_parent)
                               VALUES (:gpg_key_uid_id, :gpg_key_uid_name, :gpg_key_uid_email, :gpg_key_uid_description, :gpg_key_uid_parent)`, k)
		return err
	}
	// Update
	_, err := tx.NamedExec(`UPDATE chevron_gpg_key_uid SET
                           gpg_key_uid_name = :gpg_key_uid_name,
                           gpg_key_uid_email = :gpg_key_uid_email,
                           gpg_key_uid_description = :gpg_key_uid_description,
                           gpg_key_uid_updated_at = now()
                           WHERE gpg_key_uid_id = :gpg_key_uid_id`, k)
	return err
}
