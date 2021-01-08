package models

type GPGDecryptedData struct {
	FingerPrint          string `example:"C1CF31FB8C2A8B59"`
	Base64Data           string `example:"SGVsbG8gd29ybGQK"`
	Filename             string `example:"hello world.txt"`
	IsIntegrityProtected bool   `example:"false"`
	IsIntegrityOK        bool   `example:"false"`
}
