package main

import (
	"crypto"
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/keyBackend"
	"github.com/quan-to/remote-signer/keymagic"
	"github.com/quan-to/remote-signer/models"
	"os"
)

var (
	pgp etc.PGPInterface
	krm etc.KRMInterface
)

func Begin() {
	_ = os.Mkdir("store", os.ModePerm)
	kb := keyBackend.MakeSaveToDiskBackend("store", "key_")
	krm = keymagic.MakeKeyRingManager()
	pgp = keymagic.MakePGPManagerWithKRM(kb, krm)
	pgp.LoadKeys()
}

func AddPrivateKey(privateKeyData string) (string, error) {
	err, n := pgp.LoadKey(privateKeyData)
	if err != nil {
		return err.Error(), err
	}

	return fmt.Sprintf("Loaded %d keys", n), nil
}

func UnlockKey(fingerPrint, password string) (string, error) {
	err := pgp.UnlockKey(fingerPrint, password)
	if err != nil {
		log.Error("Error unlocking key %s: %s", fingerPrint, err)
		return err.Error(), err
	}

	return fmt.Sprintf("Key %s unlocked!", fingerPrint), nil
}

func Sign(fingerPrint, data string) (string, error) {
	signature, err := pgp.SignData(fingerPrint, []byte(data), crypto.SHA512)
	if err != nil {
		return err.Error(), err
	}
	quantoSig := remote_signer.GPG2Quanto(signature, fingerPrint, "SHA512")
	return quantoSig, nil
}

func ListPrivateKeys() ([]models.KeyInfo, error) {
	return pgp.GetLoadedPrivateKeys(), nil
}
