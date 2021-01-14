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
	ID                     string   `json:"id,omitempty"`
	FullFingerprint        string   `example:"0551F452ABE463A4"`
	Names                  []string `example:"Remote Signer Test"`
	Emails                 []string `example:"test@quan.to"`
	KeyUids                []GPGKeyUid
	KeyBits                int      `example:"3072"`
	Subkeys                []string `example:""`
	AsciiArmoredPublicKey  string   `example:"-----BEGIN PGP PUBLIC KEY BLOCK-----\nVersion: GnuPG v2\nComment: Generated by Chevron\n\n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\nFXkbQHgIocl4wcuq7ZegIr7Z7hYVd0EfOA==\n=Ufu+\n-----END PGP PUBLIC KEY BLOCK-----"`
	AsciiArmoredPrivateKey string   `example:""`
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
