package models

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/quan-to/chevron/pkg/openpgp"
	"strings"
)

const DefaultValue = -1
const DefaultPageStart = 0
const DefaultPageEnd = 100

type GPGKey struct {
	ID                     string
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
