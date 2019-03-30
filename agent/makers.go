package agent

import (
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/etc"
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
