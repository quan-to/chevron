package models

type KeyRingAddPrivateKeyData struct {
	EncryptedPrivateKey string
	SaveToDisk          bool
	Password            interface{} // Actually string, but we want to nil check it
}
