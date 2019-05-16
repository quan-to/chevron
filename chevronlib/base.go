package chevronlib

import (
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/etc/magicBuilder"
)

var pgpBackend etc.PGPInterface

func init() {
	pgpBackend = magicBuilder.MakeVoidPGP()
}
