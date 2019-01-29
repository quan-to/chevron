package smBuilder

import (
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/secretsManager"
)

func MakeSM() etc.SMInterface {
	return secretsManager.MakeSecretsManager()
}
