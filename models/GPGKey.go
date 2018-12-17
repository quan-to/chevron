package models

var GPGKeyTableInit = TableInitStruct{
	TableName:    "gpgKey",
	TableIndexes: []string{"FullFingerPrint", "Names", "Emails"},
}

type GPGKey struct {
	Id                     string
	FullFingerPrint        string
	Names                  []string
	Emails                 []string
	KeyUids                []GPGKeyUid
	KeyBits                int
	AsciiArmoredPublicKey  string
	AsciiArmoredPrivateKey string
}
