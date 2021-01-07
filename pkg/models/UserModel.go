package models

import "time"

type User struct {
	ID          string `json:"id,omitempty"`
	Fingerprint string
	Username    string
	Password    string
	FullName    string
	CreatedAt   time.Time
}

// GetID returns the id
func (u User) GetID() string {
	return u.ID
}

// GetUsername returns the username
func (u User) GetUsername() string {
	return u.Username
}

// GetFullName returns the user full name
func (u User) GetFullName() string {
	return u.FullName
}

// GetUserdata returns the raw user data
func (u User) GetUserdata() interface{} {
	return &u
}

// GetCreatedAt returns when the user was created
func (u User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

// GetFingerprint returns the user key fingerprint
func (u User) GetFingerprint() string {
	return u.Fingerprint
}
