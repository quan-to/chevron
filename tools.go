package remote_signer

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
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
