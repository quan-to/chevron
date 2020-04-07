package chevronlib

import (
	"context"
	"github.com/quan-to/chevron/pkg/interfaces"

	"github.com/quan-to/chevron/internal/etc/magicBuilder"
	"github.com/quan-to/slog"
)

var ctx = context.Background()
var pgpBackend interfaces.PGPInterface

func init() {
	slog.SetDebug(false)
	pgpBackend = magicBuilder.MakeVoidPGP(nil)
}
