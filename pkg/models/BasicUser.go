package models

import "time"

type BasicUser struct {
	ID          string
	FingerPrint string
	Username    string
	FullName    string
	CreatedAt   time.Time
}

func (bu *BasicUser) GetId() string {
	return bu.ID
}

func (bu *BasicUser) GetFullName() string {
	return bu.FullName
}

func (bu *BasicUser) GetUsername() string {
	return bu.Username
}

func (bu *BasicUser) GetUserdata() interface{} {
	return nil
}

func (bu *BasicUser) GetCreatedAt() time.Time {
	return bu.CreatedAt
}

func (bu *BasicUser) GetFingerPrint() string {
	return bu.FingerPrint
}

func (bu *BasicUser) GetToken() string {
	return ""
}
