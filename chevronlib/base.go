package chevronlib

import (
	"context"

	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/etc/magicBuilder"
	"github.com/quan-to/slog"
)

var ctx = context.Background()
var pgpBackend etc.PGPInterface

func init() {
	slog.SetDebug(false)
	pgpBackend = magicBuilder.MakeVoidPGP(nil)
}
