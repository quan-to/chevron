package uuid

import (
	"github.com/google/uuid"
	"github.com/quan-to/slog"
)

const maxUUIDTries = 10

// EnsureUUID ensures a creation of a new UUID or panics if cannot
// It does retry 10 times before panic'ing instead of uuid.Must
// The UUID.NewRandom shouldn't fail in normal scenarios
func EnsureUUID(log slog.Instance) string {
	uniqueString := ""

	for tries := 0; tries < maxUUIDTries; tries++ {
		u, err := uuid.NewRandom()
		if err == nil {
			uniqueString = u.String()
			break
		}
		if log == nil { // Only generate log if we actually
			log = slog.Scope("UUID")
		}
		log.Warn("Error generating UUID: %q. Trying again", err)
	}

	if len(uniqueString) == 0 {
		panic("cannot generate uuid. max tries reached which probably means a server machine issue")
	}

	return uniqueString
}
