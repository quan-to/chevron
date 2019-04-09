package bootstrap

import (
	"github.com/quan-to/remote-signer/database"
	"github.com/quan-to/slog"
)

var log = slog.Scope("Bootstrap")

func RunBootstraps() {
	conn := database.GetConnection()
	AddSubkeysToGPGKey(conn)
}
