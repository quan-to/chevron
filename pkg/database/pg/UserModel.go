package pg

import (
	"fmt"

	"github.com/quan-to/chevron/pkg/models"
)

func (h *PostgreSQLDBDriver) AddUser(um models.User) (string, error) {

	return "", fmt.Errorf("not supported") // TODO Implement
}

func (h *PostgreSQLDBDriver) GetUser(username string) (um *models.User, err error) {

	return um, fmt.Errorf("not supported") // TODO Implement
}

func (h *PostgreSQLDBDriver) UpdateUser(um models.User) error {

	return fmt.Errorf("not supported") // TODO Implement
}
