package tools

import (
	"bytes"
	"context"
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/quan-to/chevron/internal/models"
	"github.com/quan-to/chevron/pkg/openpgp"
	"github.com/quan-to/chevron/pkg/openpgp/armor"
	"github.com/quan-to/chevron/pkg/openpgp/packet"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/mewkiz/pkg/osutil"
	"github.com/pkg/errors"
	"github.com/quan-to/slog"
)

var toolsLog = slog.Scope("Tools")
var pgpsig = regexp.MustCompile("(?s)-----BEGIN PGP SIGNATURE-----\n(.*)-----END PGP SIGNATURE-----")

type ContextField string

const (
	CtxRequestID ContextField = "requestID"
)

func StringIndexOf(v string, a []string) int {
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
	checkSum := ""

	// check if Checksum is 5 bytes

	if !mightBeChecksum(gpgSig[:len(gpgSig)-5]) {
		// The GPG Has embedded checksum
		checkSum = ""
	} else {
		checkSum = gpgSig[len(gpgSig)-5:]
		gpgSig = gpgSig[:len(gpgSig)-5]
	}

	for i := 0; i < len(gpgSig); i++ {
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

	s := brokenMacOSXArrayFix(strings.Split(strings.Trim(signature, " \r"), "\n"), true)

	save := false

	for i := 1; i < len(s)-1; i++ {
		if !save {
			// Wait for first empty line
			if len(s[i]) == 0 {
				save = true
			}
		} else {
			cutSig += s[i]
		}
	}

	return fmt.Sprintf("%s_%s_%s", fingerPrint, hashName, cutSig)
}

func brokenMacOSXArrayFix(s []string, includeHead bool) []string {
	brokenMacOSX := true

	if includeHead {
		for i := 1; i < len(s)-1; i++ { // For Broken MacOSX Signatures
			// Search for empty lines, if there is none, its a Broken MacOSX Signature
			if len(s[i]) == 0 {
				brokenMacOSX = false
				break
			}
		}

		if brokenMacOSX {
			// Add a empty line as second line, to mandate empty header
			n := append([]string{s[0]}, "")
			s = append(n, s[1:]...)
		}
	} else {
		for i := 0; i < len(s)-1; i++ { // For Broken MacOSX Signatures, don't check last line, its not needed.
			// Search for empty lines, if there is none, its a Broken MacOSX Signature
			if len(s[i]) == 0 {
				brokenMacOSX = false
				break
			}
		}

		if brokenMacOSX {
			s = append([]string{""}, s...)
		}
	}

	return s
}

// CanDecodeBase64 returns true if s is a valid base64
func CanDecodeBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func mightBeChecksum(s string) bool {
	return (s[0] == '=' && CanDecodeBase64(s[1:])) || CanDecodeBase64(s)
}

// IsASCIIArmored returns true if data is ascii armored
func IsASCIIArmored(data string) bool {
	return data[:5] == "-----"
}

// SignatureFix recalculates the CRC
func SignatureFix(sig string) string {
	if pgpsig.MatchString(sig) {
		g := pgpsig.FindStringSubmatch(sig)
		if len(g) > 1 {
			sig = ""
			data := brokenMacOSXArrayFix(strings.Split(strings.Trim(g[1], " "), "\n"), false)
			save := false
			embeddedCrc := false

			if len(data) == 3 {
				sig = data[1] // Single Line
				if mightBeChecksum(sig[len(sig)-5:]) {
					sig = sig[:len(sig)-5]
				} else {
					embeddedCrc = true
				}
			} else {
				// PGP Has metadata header, wait for a single empty line before getting base64
				for _, v := range data {
					if !save {
						if len(v) == 0 {
							save = true // Empty line
						}
					} else if len(v) > 0 && string(v[0]) != "=" && len(v) != 5 {
						sig += v
						if len(v) == 4 {
							embeddedCrc = true
						}
					}
				}
			}

			d, err := base64.StdEncoding.DecodeString(sig)
			if err != nil {
				panic(fmt.Errorf("error decoding base64: %s", err))
			}

			crcU := make([]byte, 3)

			if !embeddedCrc {
				crc := CRC24(d)
				crcU[0] = byte((crc >> 16) & 0xFF)
				crcU[1] = byte((crc >> 8) & 0xFF)
				crcU[2] = byte(crc & 0xFF)
			}

			b64data := sig
			sig = "-----BEGIN PGP SIGNATURE-----\n"

			for i := 0; i < len(b64data); i++ {
				if i%64 == 0 {
					sig += "\n"
				}
				sig += string(b64data[i])
			}

			sig += "\n"
			if !embeddedCrc {
				sig += "=" + base64.StdEncoding.EncodeToString(crcU)
			}
			sig += "\n-----END PGP SIGNATURE-----"
		}
	}

	return sig
}

func ReadKey(armored string) (openpgp.EntityList, error) {
	kr := strings.NewReader(armored)
	keys, err := openpgp.ReadArmoredKeyRing(kr)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func GetFingerPrintFromKey(armored string) (string, error) {
	keys, err := ReadKey(armored)

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

func GetFingerPrintsFromKey(armored string) ([]string, error) {
	keys, err := ReadKey(armored)

	if err != nil {
		return nil, err
	}

	fps := make([]string, 0)

	for _, key := range keys {
		if key.PrimaryKey != nil {
			fp := ByteFingerPrint2FP16(key.PrimaryKey.Fingerprint[:])
			fps = append(fps, fp)
		}
		for _, v := range key.Subkeys {
			fp := ByteFingerPrint2FP16(v.PublicKey.Fingerprint[:])
			fps = append(fps, fp)
		}
	}

	return fps, nil
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

	if len(fps) == 0 {
		return nil, fmt.Errorf("no fingerprint found")
	}

	return fps, nil
}

func GetFingerPrintsFromEncryptedMessage(armored string) ([]string, error) {
	var fps = make([]string, 0)
	var r io.Reader

	if IsASCIIArmored(armored) {
		aem := strings.NewReader(armored)
		block, err := armor.Decode(aem)

		if err != nil {
			return nil, err
		}

		if block.Type != "PGP MESSAGE" {
			return nil, fmt.Errorf("expected pgp message but got: %s", block.Type)
		}

		r = block.Body
	} else {
		r = strings.NewReader(armored)
	}

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

	if len(fps) == 0 {
		return nil, fmt.Errorf("no fingerprints found")
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
			CreationTime:              currentTime,
			SigType:                   packet.SigTypePositiveCert,
			PubKeyAlgo:                packet.PubKeyAlgoRSA,
			Hash:                      config.Hash(),
			IsPrimaryId:               &isPrimaryId,
			FlagCertify:               true,
			FlagSign:                  true,
			FlagsValid:                true,
			FlagEncryptStorage:        true,
			FlagEncryptCommunications: true,
			IssuerKeyId:               &e.PrimaryKey.KeyId,
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
			FlagCertify:               true,
			FlagSign:                  true,
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

// CRC24 calculates the OpenPGP checksum as specified in RFC 4880, section 6.1
func CRC24(d []byte) uint32 {
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

// CopyFile copies file src to dst
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func CopyFiles(src, dst string) error {
	if !osutil.Exists(src) {
		return fmt.Errorf("folder %s does not exists", src)
	}
	if !osutil.Exists(dst) {
		return fmt.Errorf("folder %s does not exists", src)
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.Name() != "." && f.Name() != ".." && !f.IsDir() {
			srcPath := path.Join(src, f.Name())
			dstPath := path.Join(dst, f.Name())
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				toolsLog.Warn("Cannot copy %s to %s: %s", srcPath, dstPath, err)
			}
		}
	}

	return nil
}

func FolderExists(folder string) bool {
	f, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return false
	}
	if f.IsDir() {
		return true
	}

	return false
}

var identifierRegex = regexp.MustCompile("(.*) <(.*)>")

func ExtractIdentifierFields(identifier string) (name, email, comment string) {
	// TODO: Find Comment

	if identifierRegex.MatchString(identifier) {
		res := identifierRegex.FindStringSubmatch(identifier)
		if len(res) < 3 {
			return identifier, "", ""
		}
		name = res[1]
		email = res[2]

		return name, email, ""
	}

	return identifier, "", ""
}

func ContextWithValues(parent context.Context, values map[string]interface{}) context.Context {
	for k, v := range values {
		parent = context.WithValue(parent, k, v)
	}

	return parent
}

const passwordBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const defaultPasswordLength = 14

func GeneratePassword() string {
	b := make([]byte, defaultPasswordLength)
	for i := range b {
		b[i] = passwordBytes[rand.Int63()%int64(len(passwordBytes))]
	}
	return string(b)
}

// DefaultTag represents a default tag to be used in logs
const DefaultTag = "00000000+0000+0000+0000+000000000000"
const maxFieldLength = 40

// TruncateFieldForDisplay truncates the field to display in log files
func TruncateFieldForDisplay(fieldData string) string {
	fieldData = strings.Replace(fieldData, "\n", "\\n", -1)
	fieldData = strings.Replace(fieldData, "\r", "\\r", -1)

	if len(fieldData) <= maxFieldLength {
		return fieldData
	}

	return fieldData[:maxFieldLength/2] + "..." + fieldData[len(fieldData)-maxFieldLength/2:]
}

// GenerateTag generates a new UUID tag with - replaced by +
func GenerateTag() string {
	u := uuid.New()
	return strings.Replace(u.String(), "-", "+", -1)
}

func GetRequestIDFromContext(ctx context.Context) string {
	var requestID string
	ctxRequestID := ctx.Value(CtxRequestID)
	if ctxRequestID != nil {
		requestID = ctxRequestID.(string)
	} else {
		requestID = DefaultTag
	}

	return requestID
}
