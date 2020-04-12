package interfaces

// TokenManager is an interface to a Login Token Manager
type TokenManager interface {
	// AddUser adds a user to Token Manager and returns a login token
	AddUser(user UserData) string
	// AddUserWithExpiration adds an user to Token Manager that will expires in `expiration` seconds.
	AddUserWithExpiration(user UserData, expiration int) string
	// Verify verifies if the specified token is valid
	Verify(token string) error
	// GetUserData returns the user data for the specified token
	GetUserData(token string) UserData
	// InvalidateToken invalidates the specified token
	InvalidateToken(token string) error
}
