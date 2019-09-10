package bootstrap

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/database"
	"github.com/quan-to/slog"
)

var log = slog.Scope("Bootstrap").Tag(remote_signer.DefaultTag)

func RunBootstraps() {
	if remote_signer.EnableRethinkSKS || remote_signer.RethinkTokenManager || remote_signer.RethinkAuthManager {
		log.Note("Running database bootstrap because RethinkDB is enabled.")
		conn := database.GetConnection()
		AddSubkeysToGPGKey(conn)
	} else {
		log.WarnNote("RethinkDB is disabled. Skipping database bootstrap.")
	}
}
