package cache

import "github.com/quan-to/chevron/pkg/models"

// ProxiedMigrationHandler
type ProxiedMigrationHandler interface {
	InitCursor() error
	FinishCursor() error
	NextGPGKey(key *models.GPGKey) bool
	NextUser(user *models.User) bool
	NumGPGKeys() (int, error)
}

// ProxiedUserRepository a proxy to a User Repository
type ProxiedUserRepository interface {
	GetUser(username string) (um *models.User, err error)
	// AddUserToken adds a new user token to be valid and returns its token ID
	AddUserToken(ut models.UserToken) (string, error)
	// RemoveUserToken removes a user token from the database
	RemoveUserToken(token string) (err error)
	// GetUserToken fetch a UserToken object by the specified token
	GetUserToken(token string) (ut *models.UserToken, err error)
	// InvalidateUserTokens removes all user tokens that had been already expired
	InvalidateUserTokens() (int, error)
	AddUser(um models.User) (string, error)
	UpdateUser(um models.User) error
}

// ProxiedUserRepository a proxy to a GPG Repository
type ProxiedGPGRepository interface {
	// AddGPGKey adds a GPG Key to the database or update an existing one by fingerprint
	// Returns generated id / hasBeenAdded / error
	AddGPGKey(key models.GPGKey) (string, bool, error)
	// FindGPGKeyByEmail find all keys that has a underlying UID that contains that email
	FindGPGKeyByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error)
	// FindGPGKeyByFingerPrint find all keys that has a fingerprint that matches the specified fingerprint
	FindGPGKeyByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error)
	// FindGPGKeyByValue find all keys that has a underlying UID that contains that email, name or fingerprint specified by value
	FindGPGKeyByValue(value string, pageStart, pageEnd int) ([]models.GPGKey, error)
	// FindGPGKeyByName find all keys that has a underlying UID that contains that name
	FindGPGKeyByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error)
	// FetchGPGKeyByFingerprint fetch a GPG Key by its fingerprint
	FetchGPGKeyByFingerprint(fingerprint string) (*models.GPGKey, error)
	// FetchGPGKeysWithoutSubKeys fetch all keys that does not have a subkey
	// This query is not implemented on PostgreSQL
	FetchGPGKeysWithoutSubKeys() (res []models.GPGKey, err error)
	// DeleteGPGKey deletes the specified GPG key by using it's ID
	DeleteGPGKey(key models.GPGKey) error
	// UpdateGPGKey updates the specified GPG key by using it's ID
	UpdateGPGKey(key models.GPGKey) (err error)
}

// ProxiedUserRepository a proxy to a Health Checker
type ProxiedHealthChecker interface {
	// HealthCheck returns nil if everything is OK with the handler
	HealthCheck() error
}

// ProxiedHandler is a handler to be proxied by REDIS caching
type ProxiedHandler interface {
	ProxiedUserRepository
	ProxiedGPGRepository
	ProxiedHealthChecker
	ProxiedMigrationHandler
}
