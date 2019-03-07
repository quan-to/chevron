package models

type FieldDecipherInput struct {
	KeyFingerprint string
	EncryptedKey   string
	EncryptedJSON  map[string]interface{}
}
