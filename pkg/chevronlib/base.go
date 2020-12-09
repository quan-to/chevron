package chevronlib

import (
	"context"

	"github.com/quan-to/chevron/internal/etc/magicbuilder"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/database/memory"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
)

var ctx = context.Background()
var pgpBackend interfaces.PGPManager
var mem *memory.MemoryDBDriver

func init() {
	slog.SetDebug(false)
	mem = memory.MakeMemoryDBDriver(nil)
	ctx = context.WithValue(ctx, tools.CtxDatabaseHandler, mem)
	pgpBackend = magicbuilder.MakeVoidPGP(nil, mem)
}

func FolderExists(folder string) bool {
	return tools.FolderExists(folder)
}

func CopyFiles(src, dst string) error {
	return tools.CopyFiles(src, dst)
}
