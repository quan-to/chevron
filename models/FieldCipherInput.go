package models

type FieldCipherInput struct {
	JSON       map[string]interface{}
	Keys       []string
	SkipFields []string
}
