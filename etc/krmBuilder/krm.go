package krmBuilder

import (
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/pks"
)

func MakeKRM() etc.KRMInterface {
	return pks.MakeKeyRingManager()
}
