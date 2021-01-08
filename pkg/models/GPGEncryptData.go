package models

type GPGEncryptData struct {
	FingerPrint string `example:"0551F452ABE463A4"`
	Base64Data  string `example:"SGVsbG8gd29ybGQK"`
	Filename    string `example:"hello world.txt"`
	DataOnly    bool   `example:"true"`
}
