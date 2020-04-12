package fieldcipher

import (
	"encoding/base64"
	"strings"
)

func CipherPathCombine(args ...string) string {
	combined := "/"

	for _, v := range args {
		s := strings.Split(v, "/")
		if len(s) > 1 {
			// If there is a slash, then its already in cipher mode
			for _, o := range s {
				combined += o + "/"
			}
		} else {
			// If not, its a plain text
			combined += base64.StdEncoding.EncodeToString([]byte(v)) + "/"
		}
	}

	return combined
}

func CipherPathUnmangle(path string) string {
	blocks := strings.Split(path, "/")

	for i, v := range blocks {
		blkBytes, err := base64.StdEncoding.DecodeString(v)
		if err == nil {
			blocks[i] = string(blkBytes)
		}
	}

	return strings.Join(blocks, "/")
}
