package pg

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/quan-to/chevron/pkg/models"
)

func (h *PostgreSQLDBDriver) addUser(tx *sqlx.Tx, um models.User) (string, error) {
	existing, err := h.getUser(tx, um.Username)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return "", err
	}

	if existing != nil {
		return "", fmt.Errorf("already exists")
	}

	newUser := pgUserFromUser(um)
	err = newUser.save(tx)

	return newUser.ID, err
}

func (h *PostgreSQLDBDriver) getUser(tx *sqlx.Tx, username string) (um *models.User, err error) {
	user := &pgUser{}
	err = tx.Get(user, "SELECT * FROM chevron_user WHERE user_username = $1 LIMIT 1", username)
	if err != nil && !strings.EqualFold("sql: no rows in result set", err.Error()) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("not found")
	}

	um = user.toUser()

	return um, nil
}

func (h *PostgreSQLDBDriver) updateUser(tx *sqlx.Tx, um models.User) error {
	pguser := pgUserFromUser(um)

	if pguser.ID == "" {
		// Fetch user
		u, err := h.getUser(tx, pguser.Username)
		if err != nil {
			return err
		}
		pguser.ID = u.ID
	}

	return pguser.save(tx)
}
