package interfaces

// AuthManager is an interface to a Authentication Manager
// Used in Chevron Agent for Authentication StorageBackend
type AuthManager interface {
	// UserExists checks if a user with specified username exists in AuthManager
	UserExists(username string) bool
	// LoginAuth performs a login with the specified username and password
	LoginAuth(username, password string) (fingerPrint, fullname string, err error)
	// LoginAdd creates a new user in AuthManager
	LoginAdd(username, password, fullname, fingerprint string) error
	// ChangePassword changes the password of the specified user
	ChangePassword(username, password string) error
}
