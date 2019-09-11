package magicBuilder

import (
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keymagic"
	"github.com/quan-to/slog"
)

// MakeSM creates a new Instance of SecretsManager
func MakeSM(log slog.Instance) etc.SMInterface {
	return keymagic.MakeSecretsManager(log)
}
