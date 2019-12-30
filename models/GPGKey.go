package models

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/quan-to/chevron/openpgp"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"strings"
)

const DefaultValue = -1
const DefaultPageStart = 0
const DefaultPageEnd = 100

var GPGKeyTableInit = TableInitStruct{
	TableName:    "gpgKey",
	TableIndexes: []string{"FullFingerPrint", "Names", "Emails", "Subkeys"},
}

type GPGKey struct {
	Id                     string `rethinkdb:"id,omitempty"`
	FullFingerPrint        string
	Names                  []string
	Emails                 []string
	KeyUids                []GPGKeyUid
	KeyBits                int
	Subkeys                []string
	AsciiArmoredPublicKey  string
	AsciiArmoredPrivateKey string
}

func (key *GPGKey) GetShortFingerPrint() string {
	return key.FullFingerPrint[len(key.FullFingerPrint)-16:]
}

func (key *GPGKey) Save(conn *r.Session) error {
	return r.Table(GPGKeyTableInit.TableName).
		Get(key.Id).
		Update(key).
		Exec(conn)
}

func (key *GPGKey) Delete(conn *r.Session) error {
	return r.Table(GPGKeyTableInit.TableName).
		Get(key.Id).
		Delete().
		Exec(conn)
}

func AddGPGKey(conn *r.Session, data GPGKey) (string, bool, error) {
	existing, err := r.
		Table(GPGKeyTableInit.TableName).
		GetAllByIndex("FullFingerPrint", data.FullFingerPrint).
		Run(conn)

	if err != nil {
		return "", false, err
	}

	defer existing.Close()

	var gpgKey GPGKey

	if existing.Next(gpgKey) {
		// Update
		_, err := r.Table(GPGKeyTableInit.TableName).
			Get(gpgKey.Id).
			Update(data).
			RunWrite(conn)

		if err != nil {
			return "", false, err
		}

		return gpgKey.Id, false, err
	} else {
		// Create
		wr, err := r.Table(GPGKeyTableInit.TableName).
			Insert(data).
			RunWrite(conn)

		if err != nil {
			return "", false, err
		}
		return wr.GeneratedKeys[0], true, err
	}
}

func FetchKeysWithoutSubKeys(conn *r.Session) ([]GPGKey, error) {
	res, err := r.Table(GPGKeyTableInit.TableName).
		Filter(r.Row.HasFields("Subkeys").Not().Or(r.Row.Field("Subkeys").Count().Eq(0))).
		CoerceTo("array").
		Run(conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	results := make([]GPGKey, 0)
	var gpgKey GPGKey

	for res.Next(&gpgKey) {
		results = append(results, gpgKey)
	}

	return results, nil
}

func GetGPGKeyByFingerPrint(conn *r.Session, fingerPrint string) (*GPGKey, error) {
	res, err := r.Table(GPGKeyTableInit.TableName).
		Filter(r.Row.Field("FullFingerPrint").Match(fmt.Sprintf("%s$", fingerPrint)).
			Or(r.Row.HasFields("Subkeys").And(r.Row.Field("Subkeys").Filter(func(p r.Term) interface{} {
				return p.Match(fmt.Sprintf("%s$", fingerPrint))
			}).Count().Gt(0)))).
		Limit(1).
		CoerceTo("array").
		Run(conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	var gpgKey GPGKey

	if res.Next(&gpgKey) {
		return &gpgKey, nil
	}

	return nil, fmt.Errorf("not found")
}

func SearchGPGKeyByEmail(conn *r.Session, email string, pageStart, pageEnd int) ([]GPGKey, error) {
	if pageStart < 0 {
		pageStart = DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = DefaultPageEnd
	}

	var filterEmailList = func(r r.Term) interface{} {
		return r.Match(email)
	}
	res, err := r.Table(GPGKeyTableInit.TableName).
		Filter(func(r r.Term) interface{} {
			return r.Field("Emails").
				Filter(filterEmailList).
				Count().
				Gt(0)
		}).
		Slice(pageStart, pageEnd).
		CoerceTo("array").
		Run(conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()
	results := make([]GPGKey, 0)
	var gpgKey GPGKey

	for res.Next(&gpgKey) {
		results = append(results, gpgKey)
	}

	return results, nil
}

func SearchGPGKeyByFingerPrint(conn *r.Session, fingerPrint string, pageStart, pageEnd int) ([]GPGKey, error) {
	if pageStart < 0 {
		pageStart = DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = DefaultPageEnd
	}

	res, err := r.Table(GPGKeyTableInit.TableName).
		Filter(r.Row.Field("FullFingerPrint").Match(fmt.Sprintf("%s$", fingerPrint)).
			Or(r.Row.HasFields("Subkeys").And(r.Row.Field("Subkeys").Filter(func(p r.Term) interface{} {
				return p.Match(fmt.Sprintf("%s$", fingerPrint))
			}).Count().Gt(0)))).
		Slice(pageStart, pageEnd).
		CoerceTo("array").
		Run(conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()
	results := make([]GPGKey, 0)
	var gpgKey GPGKey

	for res.Next(&gpgKey) {
		results = append(results, gpgKey)
	}

	return results, nil
}

func SearchGPGKeyByValue(conn *r.Session, value string, pageStart, pageEnd int) ([]GPGKey, error) {
	if pageStart < 0 {
		pageStart = DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = DefaultPageEnd
	}
	var filterEmailList = func(r r.Term) interface{} {
		return r.Match(value)
	}

	var filterNames = func(r r.Term) interface{} {
		return r.Match(value)
	}

	var filterSub = func(r r.Term) interface{} {
		return r.Field("Emails").Filter(filterEmailList).Count().Gt(0).
			Or(r.Field("Names").Filter(filterNames).Count().Gt(0)).
			Or(r.Field("FullFingerPrint").Match(fmt.Sprintf("%s$", value)))
	}

	res, err := r.Table(GPGKeyTableInit.TableName).
		Filter(filterSub).
		Slice(pageStart, pageEnd).
		CoerceTo("array").
		Run(conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	results := make([]GPGKey, 0)
	var gpgKey GPGKey

	for res.Next(&gpgKey) {
		results = append(results, gpgKey)
	}

	return results, nil
}

func SearchGPGKeyByName(conn *r.Session, name string, pageStart, pageEnd int) ([]GPGKey, error) {
	if pageStart < 0 {
		pageStart = DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = DefaultPageEnd
	}

	var filterNames = func(r r.Term) interface{} {
		return r.Match(name)
	}
	res, err := r.Table(GPGKeyTableInit.TableName).
		Filter(func(r r.Term) interface{} {
			return r.Field("Names").
				Filter(filterNames).
				Count().
				Gt(0)
		}).
		Slice(pageStart, pageEnd).
		CoerceTo("array").
		Run(conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	results := make([]GPGKey, 0)
	var gpgKey GPGKey

	for res.Next(&gpgKey) {
		results = append(results, gpgKey)
	}

	return results, nil
}

func AsciiArmored2GPGKey(asciiArmored string) (GPGKey, error) {
	var key GPGKey
	reader := bytes.NewBuffer([]byte(asciiArmored))
	z, err := openpgp.ReadArmoredKeyRing(reader)

	if err != nil {
		return key, err
	}

	if len(z) > 0 {
		entity := z[0]
		pubKey := entity.PrimaryKey
		keyBits, _ := pubKey.BitLength()
		key = GPGKey{
			FullFingerPrint:       strings.ToUpper(hex.EncodeToString(pubKey.Fingerprint[:])),
			AsciiArmoredPublicKey: asciiArmored,
			Emails:                make([]string, 0),
			Names:                 make([]string, 0),
			KeyUids:               make([]GPGKeyUid, 0),
			KeyBits:               int(keyBits),
			Subkeys:               make([]string, 0),
		}

		fp := strings.ToUpper(hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
		key.Subkeys = append(key.Subkeys, fp[len(fp)-16:])

		for _, v := range entity.Subkeys {
			fp := strings.ToUpper(hex.EncodeToString(v.PublicKey.Fingerprint[:]))
			key.Subkeys = append(key.Subkeys, fp[len(fp)-16:])
		}

		for _, v := range entity.Identities {
			z := GPGKeyUid{
				Name:        v.UserId.Name,
				Email:       v.UserId.Email,
				Description: v.UserId.Comment,
			}
			if z.Name != "" || z.Email != "" {
				key.KeyUids = append(key.KeyUids, z)

				if z.Name != "" {
					key.Names = append(key.Names, z.Name)
				}

				if z.Email != "" {
					key.Emails = append(key.Emails, z.Email)
				}
			}
		}

		return key, nil
	}

	return key, fmt.Errorf("cannot parse GPG Key")
}
