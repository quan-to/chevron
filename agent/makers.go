package agent

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/slog"
)

// MakeTokenManager creates an instance of token manager. If Rethink is enabled returns an RethinkTokenManager, if not a MemoryTokenManager
func MakeTokenManager(logger slog.Instance) etc.TokenManager {
	if remote_signer.RethinkTokenManager {
		return MakeRethinkTokenManager(logger)
	}

	return MakeMemoryTokenManager(logger)
}

// MakeAuthManager creates an instance of auth manager. If Rethink is enabled returns an RethinkAuthManager, if not a JSONAuthManager
func MakeAuthManager(logger slog.Instance) etc.AuthManager {
	if remote_signer.RethinkAuthManager {
		return MakeRethinkAuthManager(logger)
	}

	return MakeJSONAuthManager(logger)
}
