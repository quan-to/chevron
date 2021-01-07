package pg

import (
	"github.com/quan-to/chevron/pkg/models"
)

// AddUser adds a user in the database if not exists
func (h *PostgreSQLDBDriver) AddUser(um models.User) (string, error) {
	h.log.Debug("AddUser(%s)", um.Username)
	tx, err := h.conn.Beginx()
	if err != nil {
		return "", err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	return h.addUser(tx, um)
}

// GetUser fetchs a user from the database by it's username
func (h *PostgreSQLDBDriver) GetUser(username string) (um *models.User, err error) {
	h.log.Debug("GetUser(%s)", username)
	tx, err := h.conn.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	return h.getUser(tx, username)
}

// UpdateUser updates user fingerprint, password and / or fullname by it's ID
func (h *PostgreSQLDBDriver) UpdateUser(um models.User) error {
	h.log.Debug("UpdateUser(%s)", um.Username)
	tx, err := h.conn.Beginx()
	if err != nil {
		return err
	}
	defer func() { h.rollbackIfErrorCommitIfNot(err, tx) }()

	return h.updateUser(tx, um)
}
