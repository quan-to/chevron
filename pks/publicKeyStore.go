package pks

import (
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/database"
	"github.com/quan-to/remote-signer/models"
)

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
		_, _, err := models.AddGPGKey(conn, models.AsciiArmored2GPGKey(pubKey))

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
