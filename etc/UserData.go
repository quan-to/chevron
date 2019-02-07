package etc

import "time"

type UserData interface {
	GetUsername() string
	GetUserdata() interface{}
	GetToken() string
	GetCreatedAt() time.Time
	GetFingerPrint() string
}
