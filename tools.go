package remote_signer

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/quan-to/remote-signer/models"
	"github.com/quan-to/remote-signer/openpgp"
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
	return strings.ToUpper(fmt.Sprintf("%016x", issuerKeyId))
}

func Quanto2GPG(signature string) string {
	sig := "-----BEGIN PGP SIGNATURE-----\nVersion: Quanto\n"

	s := strings.Split(signature, "$")
	if len(s) != 3 {
		s = strings.Split(signature, "_")
	}

	if len(s) != 3 {
		return ""
	}

	gpgSig := s[2]
	checkSum := gpgSig[len(gpgSig)-5:]
	for i := 0; i < len(gpgSig)-5; i++ {
		if i%64 == 0 {
			sig += "\n"
		}

		sig += string(gpgSig[i])
	}

	return sig + "\n" + checkSum + "\n-----END PGP SIGNATURE-----"
}

func GPG2Quanto(signature, fingerPrint, hash string) string {
	hashName := strings.ToUpper(hash)
	cutSig := ""

	s := strings.Split(strings.Trim(signature, " \r"), "\n")

	for i := 2; i < len(s)-1; i++ {
		cutSig += s[i]
	}

	return fmt.Sprintf("%s_%s_%s", fingerPrint, hashName, cutSig)
}

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

func GetFingerPrintFromKey(armored string) (string, error) {
	kr := strings.NewReader(armored)
	keys, err := openpgp.ReadArmoredKeyRing(kr)
	if err != nil {
		return "", err
	}

	for _, key := range keys {
		if key.PrimaryKey != nil {
			fp := ByteFingerPrint2FP16(key.PrimaryKey.Fingerprint[:])

			return fp, nil
		}
	}

	return "", fmt.Errorf("cannot read key")
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

func CreateEntityForSubKey(masterFingerPrint string, pubKey *packet.PublicKey, privKey *packet.PrivateKey) *openpgp.Entity {
	uid := packet.NewUserId(fmt.Sprintf("Subkey for %s", masterFingerPrint), "", "")

	e := openpgp.Entity{
		PrimaryKey: pubKey,
		PrivateKey: privKey,
		Identities: make(map[string]*openpgp.Identity),
	}

	e.Identities[uid.Id] = &openpgp.Identity{
		Name:   uid.Name,
		UserId: uid,
	}

	e.Subkeys = make([]openpgp.Subkey, 0)
	return &e
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

func IdentityMapToArray(m map[string]*openpgp.Identity) []*openpgp.Identity {
	arr := make([]*openpgp.Identity, 0)

	for _, v := range m {
		arr = append(arr, v)
	}

	return arr
}

func SimpleIdentitiesToString(ids []*openpgp.Identity) string {
	identifier := ""
	for _, k := range ids {
		identifier = k.Name
		break
	}

	return identifier
}

func ReadKeyToEntity(asciiArmored string) (*openpgp.Entity, error) {
	r := strings.NewReader(asciiArmored)
	e, err := openpgp.ReadArmoredKeyRing(r)

	if err != nil {
		return nil, err
	}

	if len(e) > 0 {
		return e[0], nil
	}

	return nil, errors.New("no keys found")
}

func CompareFingerPrint(fpA, fpB string) bool {
	if fpA == "" || fpB == "" {
		return false
	}

	if len(fpA) == len(fpB) {
		return fpA == fpB
	}

	if len(fpA) > len(fpB) {
		return fpA[len(fpA)-len(fpB):] == fpB
	}

	return fpB[len(fpB)-len(fpA):] == fpA
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
