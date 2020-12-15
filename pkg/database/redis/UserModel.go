package redis

import "github.com/quan-to/chevron/pkg/models"

// AddUser adds a user in the database if not exists
func (h *Driver) AddUser(um models.User) (string, error) {
	h.log.Debug("AddUser(%s)", um.Username)

	return h.proxy.AddUser(um)
}

// GetUser fetchs a user from the database by it's username
func (h *Driver) GetUser(username string) (um *models.User, err error) {
	h.log.Debug("GetUser(%s)", username)

	return h.proxy.GetUser(username)
}

// UpdateUser updates user fingerprint, password and / or fullname by it's ID
func (h *Driver) UpdateUser(um models.User) error {
	h.log.Debug("UpdateUser(%s)", um.Username)

	return h.proxy.UpdateUser(um)
}
