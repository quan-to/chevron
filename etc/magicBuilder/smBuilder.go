package magicBuilder

import (
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keymagic"
	"github.com/quan-to/slog"
)

func MakeSM(log slog.Instance) etc.SMInterface {
	return keymagic.MakeSecretsManager(log)
}
