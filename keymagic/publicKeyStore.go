package keymagic

import (
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/database"
	"github.com/quan-to/remote-signer/models"
)

var pksLog = SLog.Scope("PKS")

func PKSGetKey(fingerPrint string) string {
	if !remote_signer.EnableRethinkSKS {
		return GetSKSKey(fingerPrint)
	}

	conn := database.GetConnection()
	v := models.GetGPGKeyByFingerPrint(conn, fingerPrint)

	if v != nil {
		return v.AsciiArmoredPublicKey
	}

	return ""
}

func PKSSearchByName(name string, pageStart, pageEnd int) []models.GPGKey {
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByName(conn, name, pageStart, pageEnd)
	}
	panic("The server does not have RethinkDB enabled so it cannot serve search")
}

func PKSSearchByFingerPrint(fingerPrint string, pageStart, pageEnd int) []models.GPGKey {
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByFingerPrint(conn, fingerPrint, pageStart, pageEnd)
	}
	panic("The server does not have RethinkDB enabled so it cannot serve search")
}

func PKSSearchByEmail(email string, pageStart, pageEnd int) []models.GPGKey {
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByEmail(conn, email, pageStart, pageEnd)
	}
	panic("The server does not have RethinkDB enabled so it cannot serve search")
}

func PKSSearch(value string, pageStart, pageEnd int) []models.GPGKey {
	//if EnableRethinkSKS {
	//	conn := GetConnection()
	//	return models.SearchGPGKeyByEmail(conn, email, pageStart, pageEnd)
	//}
	//panic("The server does not have RethinkDB enabled so it cannot serve search")
	panic("Not implemented")
}

func PKSAdd(pubKey string) string {
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		key := models.AsciiArmored2GPGKey(pubKey)
		pksLog.Info("Adding public key %s to PKS", key.GetShortFingerPrint())
		_, _, err := models.AddGPGKey(conn, key)

		if err != nil {
			panic(err)
		}

		return "OK"
	}

	if PutSKSKey(pubKey) {
		return "OK"
	}

	return "NOK"
}
