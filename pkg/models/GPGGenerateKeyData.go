package models

type GPGGenerateKeyData struct {
	Identifier string `example:"John HUEBR <john@huebr.com>"`
	Password   string `example:"I think you will never guess"`
	Bits       int    `example:"4096"`
}
