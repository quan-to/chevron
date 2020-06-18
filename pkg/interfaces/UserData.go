package interfaces

import "time"

// UserData is an interface for user data
type UserData interface {
	// GetUsername returns the username
	GetUsername() string
	// GetFullName returns the user full name
	GetFullName() string
	// GetUserdata returns the raw user data
	GetUserdata() interface{}
	// GetToken returns the user token
	GetToken() string
	// GetCreatedAt returns when the user was created
	GetCreatedAt() time.Time
	// GetFingerPrint returns the user key fingerprint
	GetFingerPrint() string
}
