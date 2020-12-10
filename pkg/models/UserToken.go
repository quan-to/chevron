package models

import (
	"time"
)

type UserToken struct {
	ID          string `json:"id,omitempty"`
	Fingerprint string
	Username    string
	Fullname    string
	Token       string
	CreatedAt   time.Time
	Expiration  time.Time
}

// GetId returns the id
func (ut *UserToken) GetId() string {
	return ut.ID
}

// GetFingerPrint returns the user key fingerprint
func (ut *UserToken) GetFingerPrint() string {
	return ut.Fingerprint
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

func (ut *UserToken) GetFingerprint() string {
	return ut.Fingerprint
}
