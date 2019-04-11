package agent

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
)

func MakeTokenManager() etc.TokenManager {
	if remote_signer.RethinkTokenManager {
		return MakeRethinkTokenManager()
	}

	return MakeMemoryTokenManager()
}

func MakeAuthManager() etc.AuthManager {
	if remote_signer.RethinkAuthManager {
		return MakeRethinkAuthManager()
	}

	return MakeJSONAuthManager()
}
