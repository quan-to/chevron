package models

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var GPGKeyTableInit = TableInitStruct{
	TableName:    "gpgKey",
	TableIndexes: []string{"FullFingerPrint", "Names", "Emails"},
}

type GPGKey struct {
	Id                     string
	FullFingerPrint        string
	Names                  []string
	Emails                 []string
	KeyUids                []GPGKeyUid
	KeyBits                int
	AsciiArmoredPublicKey  string
	AsciiArmoredPrivateKey string
}

func AddGPGKey(conn *r.Session, data GPGKey) (string, error) {
	existing, err := r.
		Table(GPGKeyTableInit.TableName).
		GetAllByIndex("FullFingerPrint", data.FullFingerPrint).
		Run(conn)

	if err != nil {
		return "", err
	}

	var gpgKey GPGKey

	if existing.Next(gpgKey) {
		// Update
		_, err := r.Table(GPGKeyTableInit.TableName).
			Get(gpgKey.Id).
			Update(data).
			RunWrite(conn)

		if err != nil {
			return "", err
		}

		return gpgKey.Id, err
	} else {
		// Create
		wr, err := r.Table(GPGKeyTableInit.TableName).
			Insert(data).
			RunWrite(conn)

		if err != nil {
			return "", err
		}

		return wr.GeneratedKeys[0], err
	}
}
