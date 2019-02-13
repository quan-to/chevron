package etc

type SMInterface interface {
	PutKeyPassword(fingerPrint, password string)
	PutEncryptedPassword(fingerPrint, encryptedPassword string)
	GetPasswords() map[string]string
	UnlockLocalKeys(gpg PGPInterface)
	GetMasterKeyFingerPrint() string
}
