package fieldcipher

import (
	"github.com/quan-to/chevron/rstest"
	"testing"
)

func TestCipher_GenerateEncryptedPacket(t *testing.T) {
	cipher := MakeCipherFromASCIIArmoredKeys([]string{rstest.TestPublicKey})

	dataToCipher := map[string]interface{}{
		"a": "b",
		"c": "d",
		"e": map[string]interface{}{
			"o": []string{"1", "2", "4"},
			"v": []interface{}{1, "2", true},
			"k": nil,
		},
		"bb": true,
		"oe": 1234.5,
		"v":  nil,
	}

	skipFields := []string{CipherPathCombine("a"), CipherPathCombine("oe")}

	packet, err := cipher.GenerateEncryptedPacket(dataToCipher, skipFields)
	if err != nil {
		t.Errorf(err.Error())
	}

	cipheredData := packet.EncryptedJSON
	// Test for skipFields

	if dataToCipher["a"] != cipheredData["a"] {
		t.Errorf("expected /a to be %v got %v", dataToCipher["a"], cipheredData["a"])
	}

	if dataToCipher["oe"] != cipheredData["oe"] {
		t.Errorf("expected /oe to be %v got %v", dataToCipher["oe"], cipheredData["oe"])
	}
}
