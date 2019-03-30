package magicBuilder

import (
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/keymagic"
)

func MakeSM() etc.SMInterface {
	return keymagic.MakeSecretsManager()
}
