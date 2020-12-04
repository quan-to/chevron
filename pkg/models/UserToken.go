package models

import (
	"time"
)

type UserToken struct {
	ID          string `json:"id,omitempty"`
	FingerPrint string
	Username    string
	Fullname    string
	Token       string
	CreatedAt   time.Time
	Expiration  time.Time
}

func (ut *UserToken) GetUsername() string {
	return ut.Username
}

func (ut *UserToken) GetFullName() string {
	return ut.Fullname
}

func (ut *UserToken) GetUserdata() interface{} {
	return nil
}

func (ut *UserToken) GetToken() string {
	return ut.Token
}

func (ut *UserToken) GetCreatedAt() time.Time {
	return ut.CreatedAt
}

func (ut *UserToken) GetFingerPrint() string {
	return ut.FingerPrint
}
