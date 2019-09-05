package agent

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/slog"
)

func MakeTokenManager(logger slog.Instance) etc.TokenManager {
	if remote_signer.RethinkTokenManager {
		return MakeRethinkTokenManager(logger)
	}

	return MakeMemoryTokenManager(logger)
}

func MakeAuthManager(logger slog.Instance) etc.AuthManager {
	if remote_signer.RethinkAuthManager {
		return MakeRethinkAuthManager(logger)
	}

	return MakeJSONAuthManager(logger)
}
