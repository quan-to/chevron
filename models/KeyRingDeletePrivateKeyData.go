package models

type KeyRingDeletePrivateKeyData struct {
	EncryptedPrivateKey string
	Password            interface{} // Actually string, but we want to nil check it
}
