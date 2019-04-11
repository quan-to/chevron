package bootstrap

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/database"
	"github.com/quan-to/slog"
)

var log = slog.Scope("Bootstrap")

func RunBootstraps() {
	if remote_signer.EnableRethinkSKS || remote_signer.RethinkTokenManager || remote_signer.RethinkAuthManager {
		conn := database.GetConnection()
		AddSubkeysToGPGKey(conn)
	}
}
