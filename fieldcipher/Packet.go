package fieldcipher

type CipherPacket struct {
	EncryptedKey  string
	EncryptedJSON map[string]interface{}
}

type UnmatchedField struct {
	Expected string
	Got      string
}

type DecipherPacket struct {
	DecryptedData   map[string]interface{}
	JSONChanged     bool
	UnmatchedFields []UnmatchedField
}
