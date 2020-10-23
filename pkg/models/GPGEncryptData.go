package models

type GPGEncryptData struct {
	FingerPrint string
	Base64Data  string
	Filename    string
	DataOnly    bool
}
