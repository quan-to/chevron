package fieldcipher

import (
	"github.com/quan-to/remote-signer"
	"io/ioutil"
	"testing"
)

func TestDecipher_DecipherPacket(t *testing.T) {

	keyData, err := ioutil.ReadFile("../tests/testkey_privateTestKey.gpg")

	if err != nil {
		t.Fatalf("Error reading private key: %s", err)
	}

	keyPass, err := ioutil.ReadFile("../tests/testprivatekeyPassword.txt")

	if err != nil {
		t.Fatalf("Error reading private key password: %s", err)
	}

	cipher := MakeCipherFromASCIIArmoredKeys([]string{remote_signer.TestPublicKey})

	dataToCipher := map[string]interface{}{
		"a": "b",
		"c": "d",
		"e": map[string]interface{}{
			"o": []string{"1", "2", "4"},
			"v": []interface{}{1, "2", true},
			"t": nil,
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

	// Invert two fields
	bb := packet.EncryptedJSON["bb"]
	oe := packet.EncryptedJSON["oe"]
	packet.EncryptedJSON["bb"] = oe
	packet.EncryptedJSON["oe"] = bb

	// Decrypt
	decipher, err := MakeDecipherWithASCIIPrivateKey(string(keyData))

	if err != nil {
		t.Fatalf("Error loading private key: %s", err)
	}

	if !decipher.Unlock(string(keyPass)) {
		t.Fatalf("Error decrypting private key")
	}

	decPacket, err := decipher.DecipherPacket(*packet)

	if err != nil {
		t.Fatalf("Error decrypting packet: %s", err)
	}

	if len(decPacket.UnmatchedFields) != 1 {
		t.Errorf("expected a single unmatched field")
	}

	if len(decPacket.UnmatchedFields) == 1 {
		um := decPacket.UnmatchedFields[0]
		if um.Expected != "/bb/" || um.Got != "/oe/" {
			t.Errorf("expected unmatched field to be (Expected: /bb/, Got: /oe/) but got (Expected: %s, Got: %s)", um.Expected, um.Got)
		}
	}
}
