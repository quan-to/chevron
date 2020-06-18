package magicbuilder

import (
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
)

// MakeSM creates a new Instance of SecretsManager
func MakeSM(log slog.Instance) interfaces.SecretsManager {
	return keymagic.MakeSecretsManager(log)
}
