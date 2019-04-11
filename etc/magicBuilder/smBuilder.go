package magicBuilder

import (
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keymagic"
)

func MakeSM() etc.SMInterface {
	return keymagic.MakeSecretsManager()
}
