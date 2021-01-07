package pg

import (
	"fmt"

	"github.com/quan-to/chevron/pkg/models"
)

func (h *PostgreSQLDBDriver) AddUserToken(ut models.UserToken) (string, error) {
	return "", fmt.Errorf("token is not supported on postgres. please use redis wrapper around it")
}

// RemoveUserToken removes a user token from the database
func (h *PostgreSQLDBDriver) RemoveUserToken(token string) (err error) {
	return fmt.Errorf("token is not supported on postgres. please use redis wrapper around it")
}

func (h *PostgreSQLDBDriver) GetUserToken(token string) (ut *models.UserToken, err error) {
	return nil, fmt.Errorf("token is not supported on postgres. please use redis wrapper around it")
}

func (h *PostgreSQLDBDriver) InvalidateUserTokens() (int, error) {
	return 0, fmt.Errorf("token is not supported on postgres. please use redis wrapper around it")
}
