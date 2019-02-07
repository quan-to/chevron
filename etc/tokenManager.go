package etc

type TokenManager interface {
	AddUser(user UserData) string
	Verify(token string) error
	GetUserData(token string) UserData
}
