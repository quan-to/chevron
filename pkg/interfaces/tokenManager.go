package interfaces

type TokenManager interface {
	AddUser(user UserData) string
	AddUserWithExpiration(user UserData, expiration int) string
	Verify(token string) error
	GetUserData(token string) UserData
	InvalidateToken(token string) error
}
