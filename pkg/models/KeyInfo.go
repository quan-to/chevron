package models

type KeyInfo struct {
	FingerPrint           string `example:"0551F452ABE463A4"`
	Identifier            string `example:"Remote Signer Test <test@quan.to>"`
	Bits                  int    `example:"3072"`
	ContainsPrivateKey    bool   `example:"false"`
	PrivateKeyIsDecrypted bool   `example:"false"`
}
