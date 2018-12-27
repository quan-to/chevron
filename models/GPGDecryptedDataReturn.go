package models


type GPGDecryptedData struct {
	FingerPrint string
	Base64Data string
	Filename string
	IsIntegrityProtected bool
	IsIntegrityOK bool
}