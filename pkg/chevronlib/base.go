package chevronlib

import (
	"context"
	"github.com/quan-to/chevron/pkg/interfaces"

	"github.com/quan-to/chevron/internal/etc/magicbuilder"
	"github.com/quan-to/slog"
)

var ctx = context.Background()
var pgpBackend interfaces.PGPManager

func init() {
	slog.SetDebug(false)
	pgpBackend = magicbuilder.MakeVoidPGP(nil)
}
