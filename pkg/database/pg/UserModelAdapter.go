package pg

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quan-to/chevron/pkg/uuid"

	"github.com/quan-to/chevron/pkg/models"
)

type pgUser struct {
	ID          string     `db:"user_id"`
	Fingerprint string     `db:"user_fingerprint"`
	Username    string     `db:"user_username"`
	Password    []byte     `db:"user_password"`
	FullName    string     `db:"user_full_name"`
	CreatedAt   time.Time  `db:"user_created_at"`
	UpdatedAt   time.Time  `db:"user_updated_at"`
	DeletedAt   *time.Time `db:"user_deleted_at"`
}

func (u *pgUser) toUser() *models.User {
	return &models.User{
		ID:          u.ID,
		Fingerprint: u.Fingerprint,
		Username:    u.Username,
		Password:    string(u.Password),
		FullName:    u.FullName,
		CreatedAt:   u.CreatedAt,
	}
}

func pgUserFromUser(um models.User) *pgUser {
	return &pgUser{
		ID:          um.ID,
		Fingerprint: um.Fingerprint,
		Username:    um.Username,
		Password:    []byte(um.Password),
		FullName:    um.FullName,
		CreatedAt:   um.CreatedAt,
	}
}

func (u *pgUser) save(tx *sqlx.Tx) error {
	if u.ID == "" { // Insert
		u.ID = uuid.EnsureUUID(nil)
		_, err := tx.NamedExec(`INSERT INTO 
            chevron_user(user_id, user_fingerprint, user_username, user_password, user_full_name, user_created_at) 
            VALUES (:user_id, :user_fingerprint, :user_username, :user_password, :user_full_name, now())`, u)
		if err != nil {
			return err
		}
		return nil
	}

	// Update
	_, err := tx.NamedExec(`UPDATE chevron_user SET
                           user_fingerprint = :user_fingerprint,
                           user_password = :user_password,
                           user_full_name = :user_full_name,
                           user_updated_at = now()
                           WHERE user_id = :user_id`, u)
	return err
}
