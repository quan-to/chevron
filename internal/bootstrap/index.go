package bootstrap

import (
	config "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/database"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/slog"
)

var log = slog.Scope("Bootstrap").Tag(tools.DefaultTag)

func RunBootstraps() {
	if config.EnableRethinkSKS || config.RethinkTokenManager || config.RethinkAuthManager {
		log.Note("Running database bootstrap because RethinkDB is enabled.")
		conn := database.GetConnection()
		AddSubkeysToGPGKey(conn)
	} else {
		log.WarnNote("RethinkDB is disabled. Skipping database bootstrap.")
	}
}
