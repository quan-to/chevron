package models

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/openpgp"
)

const DefaultValue = -1
const DefaultPageStart = 0
const DefaultPageEnd = 100

type GPGKey struct {
	ID                     string `json:"id,omitempty"`
	FullFingerprint        string
	Names                  []string
	Emails                 []string
	KeyUids                []GPGKeyUid
	KeyBits                int
	Subkeys                []string
	AsciiArmoredPublicKey  string
	AsciiArmoredPrivateKey string
	ParentKey              *GPGKey
}

func (key *GPGKey) GetShortFingerPrint() string {
	return tools.FPto16(key.FullFingerprint)
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
			FullFingerprint:       strings.ToUpper(hex.EncodeToString(pubKey.Fingerprint[:])),
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
