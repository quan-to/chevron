package models


type KeyInfo struct {
	 FingerPrint string
	 Identifier string
	 Bits int
	 ContainsPrivateKey bool
	 PrivateKeyIsDecrypted bool
}