package etc

type AuthManager interface {
	UserExists(username string) bool
	LoginAuth(username, password string) (fingerPrint string, err error)
	LoginAdd(username, password, fullname, fingerprint string) error
	ChangePassword(username, password string) error
}
