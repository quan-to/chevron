package testmodels

import (
	"time"

	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/pkg/uuid"
)

// GpgKey is a models.GPGKey instance for use in tests
var GpgKey = models.GPGKey{
	ID:              "abcd",
	FullFingerprint: "DEADBEEFDEADBEEFDEADBEEFDEADBEEF",
	Names:           []string{"AbCE", "B ASD"},
	Emails:          []string{"a@a.com", "b@a.com"},
	KeyUids: []models.GPGKeyUid{
		{
			Name:        "AbCE",
			Email:       "a@a.com",
			Description: "desc",
		},
		{
			Name:        "B ASD",
			Email:       "b@a.com",
			Description: "desc",
		},
	},
	KeyBits:                1234,
	Subkeys:                []string{"BABABEBE"},
	AsciiArmoredPublicKey:  "PUBKEY",
	AsciiArmoredPrivateKey: "PRIVKEY",
}

// User is a models.User instance for using in tests
var User = models.User{
	ID:          "abcd",
	Username:    "johnhuebr",
	FullName:    "John HUEBR",
	Fingerprint: "DEADBEEFDEADBEEF",
	Password:    "I think you will never guess",
	CreatedAt:   time.Now().Truncate(time.Second),
}

// Time is a time constant used for tests
var Time = time.Now().Truncate(time.Second)

// Token is a models.UserToken instance for using in tests
var Token = models.UserToken{
	Fingerprint: "DEADBEEF",
	Username:    "johnhuebr",
	Fullname:    "John HUEBR",
	Token:       uuid.EnsureUUID(nil),
	CreatedAt:   Time.Truncate(time.Second),
	Expiration:  Time.Add(time.Hour).Truncate(time.Second),
}
