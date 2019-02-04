package keymagic

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/database"
	"github.com/quan-to/remote-signer/models"
)

var pksLog = SLog.Scope("PKS")

func PKSGetKey(fingerPrint string) (string, error) {
	if !remote_signer.EnableRethinkSKS {
		return GetSKSKey(fingerPrint)
	}

	conn := database.GetConnection()
	v, err := models.GetGPGKeyByFingerPrint(conn, fingerPrint)

	if v != nil {
		return v.AsciiArmoredPublicKey, nil
	}

	return "", err
}

func PKSSearchByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByName(conn, name, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have RethinkDB enabled so it cannot serve search")
}

func PKSSearchByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByFingerPrint(conn, fingerPrint, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have RethinkDB enabled so it cannot serve search")
}

func PKSSearchByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByEmail(conn, email, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have RethinkDB enabled so it cannot serve search")
}

func PKSSearch(value string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByValue(conn, value, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have RethinkDB enabled so it cannot serve search")
}

func PKSAdd(pubKey string) string {
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		key, err := models.AsciiArmored2GPGKey(pubKey)
		if err != nil {
			return "NOK"
		}

		keys, err := models.SearchGPGKeyByFingerPrint(conn, key.FullFingerPrint, 0, 1)

		if err != nil {
			return "NOK"
		}

		if len(keys) > 0 {
			pksLog.Info("Tried to add key %s to PKS but already exists.", key.GetShortFingerPrint())
			return "OK"
		}

		pksLog.Info("Adding public key %s to PKS", key.GetShortFingerPrint())
		_, _, err = models.AddGPGKey(conn, key)

		if err != nil {
			return "NOK"
		}

		return "OK"
	}

	res, _ := PutSKSKey(pubKey)

	if res {
		return "OK"
	}

	return "NOK"
}
