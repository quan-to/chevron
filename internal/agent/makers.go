package agent

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
)

// MakeTokenManager creates an instance of token manager. If Rethink is enabled returns an RethinkTokenManager, if not a MemoryTokenManager
func MakeTokenManager(logger slog.Instance) interfaces.TokenManager {
	if config.RethinkTokenManager {
		return MakeRethinkTokenManager(logger)
	}

	return MakeMemoryTokenManager(logger)
}

// MakeAuthManager creates an instance of auth manager. If Rethink is enabled returns an RethinkAuthManager, if not a JSONAuthManager
func MakeAuthManager(logger slog.Instance) interfaces.AuthManager {
	if config.RethinkAuthManager {
		return MakeRethinkAuthManager(logger)
	}

	return MakeJSONAuthManager(logger)
}
