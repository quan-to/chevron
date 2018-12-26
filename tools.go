package remote_signer

import (
	"encoding/hex"
	"strings"
)

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
