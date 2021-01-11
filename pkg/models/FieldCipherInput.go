package models

type FieldCipherInput struct {
	JSON       map[string]interface{}
	Keys       []string `example:"0551F452ABE463A4"`
	SkipFields []string `example:""`
}
