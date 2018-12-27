package remote_signer

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/quan-to/remote-signer/models"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"regexp"
	"strings"
)

var pgpsig = regexp.MustCompile("-----BEGIN PGP SIGNATURE-----(.*)-----END PGP SIGNATURE-----")

func stringIndexOf(v string, a []string) int {
	for i, vo := range a {
		if vo == v {
			return i
		}
	}

	return -1
}

func ByteFingerPrint2FP16(raw []byte) string {
	fp := hex.EncodeToString(raw)
	return strings.ToUpper(fp[len(fp)-16:])
}

func IssuerKeyIdToFP16(issuerKeyId uint64) string {
	fp := strings.ToUpper(fmt.Sprintf("%016x", issuerKeyId))
	if len(fp) > 16 {
		return fp[len(fp)-16:]
	} else {
		return fp
	}
}

// region CRC24 from https://github.com/golang/crypto/blob/master/openpgp/armor/armor.go
const crc24Init = 0xb704ce
const crc24Poly = 0x1864cfb

// crc24 calculates the OpenPGP checksum as specified in RFC 4880, section 6.1
func crc24(d []byte) uint32 {
	crc := uint32(crc24Init)
	for _, b := range d {
		crc ^= uint32(b) << 16
		for i := 0; i < 8; i++ {
			crc <<= 1
			if crc&0x1000000 != 0 {
				crc ^= crc24Poly
			}
		}
	}
	return crc
}

// endregion

func signatureFix(sig string) string {
	if pgpsig.MatchString(sig) {
		g := pgpsig.FindStringSubmatch(sig)
		if len(g) > 1 {
			sig = ""
			data := strings.Split(strings.Trim(g[0], " "), "\n")
			save := false
			if len(data) == 1 {
				sig = data[0]
			} else {
				for _, v := range data {
					if !save {
						save = save || len(v) > 0
						if len(v) > 2 && v[:2] == "iQ" { // Workarround for GPG Bug in Production
							save = true
							sig += v
						} else {
							sig += v
						}
					}
				}
			}

			d, err := base64.StdEncoding.DecodeString(sig)
			if err != nil {
				panic(err)
			}

			crc := crc24(d)
			crcU := make([]byte, 3)
			crcU[0] = byte((crc >> 16) & 0xFF)
			crcU[1] = byte((crc >> 8) & 0xFF)
			crcU[2] = byte(crc & 0xFF)

			sig = "-----BEGIN PGP SIGNATURE-----\n\n" + sig + "\n=" + base64.StdEncoding.EncodeToString(crcU) + "\n-----END PGP SIGNATURE-----"
		}
	}

	return sig
}

func GetFingerPrintFromKey(armored string) string {
	kr := strings.NewReader(armored)
	keys, err := openpgp.ReadArmoredKeyRing(kr)
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		if key.PrivateKey != nil {
			fp := ByteFingerPrint2FP16(key.PrimaryKey.Fingerprint[:])

			return fp
		}
	}

	return ""
}

func GetFingerPrintsFromEncryptedMessageRaw(rawB64Data string) ([]string, error) {
	var fps = make([]string, 0)
	data, err := base64.StdEncoding.DecodeString(rawB64Data)

	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(data)

	reader := packet.NewReader(r)

	for {
		p, err := reader.Next()

		if err != nil {
			break
		}

		switch v := p.(type) {
		case *packet.EncryptedKey:
			fps = append(fps, IssuerKeyIdToFP16(v.KeyId))
		}
	}

	return fps, nil
}

func GetFingerPrintsFromEncryptedMessage(armored string) ([]string, error) {
	var fps = make([]string, 0)
	aem := strings.NewReader(armored)
	block, err := armor.Decode(aem)

	if err != nil {
		return nil, err
	}

	if block.Type != "PGP MESSAGE" {
		return nil, fmt.Errorf("expected pgp message but got: %s", block.Type)
	}

	reader := packet.NewReader(block.Body)

	for {
		p, err := reader.Next()

		if err != nil {
			break
		}

		switch v := p.(type) {
		case *packet.EncryptedKey:
			fps = append(fps, IssuerKeyIdToFP16(v.KeyId))
		}
	}

	return fps, nil
}

func CreateEntityFromKeys(name, comment, email string, lifeTimeInSecs uint32, pubKey *packet.PublicKey, privKey *packet.PrivateKey) *openpgp.Entity {
	bitLen, _ := privKey.BitLength()
	config := packet.Config{
		DefaultHash:            crypto.SHA512,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		CompressionConfig: &packet.CompressionConfig{
			Level: 9,
		},
		RSABits: int(bitLen),
	}
	currentTime := config.Now()
	uid := packet.NewUserId(name, comment, email)

	e := openpgp.Entity{
		PrimaryKey: pubKey,
		PrivateKey: privKey,
		Identities: make(map[string]*openpgp.Identity),
	}
	isPrimaryId := false

	e.Identities[uid.Id] = &openpgp.Identity{
		Name:   uid.Name,
		UserId: uid,
		SelfSignature: &packet.Signature{
			CreationTime: currentTime,
			SigType:      packet.SigTypePositiveCert,
			PubKeyAlgo:   packet.PubKeyAlgoRSA,
			Hash:         config.Hash(),
			IsPrimaryId:  &isPrimaryId,
			FlagsValid:   true,
			FlagSign:     true,
			FlagCertify:  true,
			IssuerKeyId:  &e.PrimaryKey.KeyId,
		},
	}

	e.Subkeys = make([]openpgp.Subkey, 1)
	e.Subkeys[0] = openpgp.Subkey{
		PublicKey:  pubKey,
		PrivateKey: privKey,
		Sig: &packet.Signature{
			CreationTime:              currentTime,
			SigType:                   packet.SigTypeSubkeyBinding,
			PubKeyAlgo:                packet.PubKeyAlgoRSA,
			Hash:                      config.Hash(),
			PreferredHash:             []uint8{models.GPG_SHA512},
			FlagsValid:                true,
			FlagEncryptStorage:        true,
			FlagEncryptCommunications: true,
			IssuerKeyId:               &e.PrimaryKey.KeyId,
			KeyLifetimeSecs:           &lifeTimeInSecs,
		},
	}
	return &e
}
