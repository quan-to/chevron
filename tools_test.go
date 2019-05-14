package remote_signer

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"github.com/quan-to/chevron/openpgp"
	"github.com/quan-to/chevron/openpgp/armor"
	"github.com/quan-to/chevron/openpgp/packet"
	"github.com/quan-to/chevron/rstest"
	"io/ioutil"
	"path"
	"strings"
	"testing"
	"time"
)

func TestStringIndexOf(t *testing.T) {
	v := []string{"a", "v", "o"}

	for i := range v {
		idx := StringIndexOf(v[i], v)
		if idx != i {
			t.Errorf("Expected %d got %d", i, idx)
		}
	}

	idx := StringIndexOf("huebrbrb", v)
	if idx != -1 {
		t.Errorf("Expected %d got %d", -1, idx)
	}
}

func TestIssuerKeyIdToFP16(t *testing.T) {
	v := uint64(0xFFFF0000FFFF0000)
	o := IssuerKeyIdToFP16(v)

	if o != "FFFF0000FFFF0000" {
		t.Errorf("Expected FFFF0000FFFF0000 got %s", o)
	}

	v = uint64(0xFFFF0000)
	o = IssuerKeyIdToFP16(v)

	if o != "00000000FFFF0000" {
		t.Errorf("Expected 00000000FFFF0000 got %s", o)
	}
}

const sigConvertQuanto = "0ADF79401F28C569_SHA512_iQIzBAEBCgAdFiEEab8JRxWM7/xGsOGsRxsIMMDGp/EFAlw3bgwACgkQRxsIMMDGp/Gq+hAAooiGdBZl0z1+uZs6voUEPloIl0qYxSuDdgI2QAdTiALcbasuzhYge04exIgpXf6Exik3TH4Qop5RqpvbDRK5J5AYvWdst377NSIL/m00X44hU3Mq3oJ52LyTCj3qShMDkviXtm7GynoXNFaloPwxs3hXze3E+ddWVn17Nw9tIAJbdeWOMRbWSdpijAsOZP6qGvrjejNCA3eQSTb2G15zB69yS///mgeRVLNGC7YHzbgX3VROXix6pcdc8LOgZolloey7VkrOkvBg7t9n2VpqMti1qUQ3qGVLx27YyKjjI+mykUnoO2i5KzsMfZVCB9iQC3FgVmaGElLUxVJGGToByw4QNuTsLNeVchd+nA20dhQmmZ2dmaMpUIOl0TbL3wxPxa7eJ72fx3+6EQIqQw0t6ScauPfEQ7Ad0ORIEhGvRXhNYykNUVgdoH09FoF1eEZv2yvJK5UDQNDUifTnhJ+7A1r7jgykE3vqcrcegbJahC0Qjn66316+D1O/6I5E/ZZtx3zuzJQT9kTawDTslnmgg5XhQ9LmsrjBYpSKNspAvlhonue07XVyekO1u6UaKTOmGG060dInWby5Xf+YAK7W8a7Iucoq3zPM0Y6eMVDMNcGcLWhcyCnnFRhOrGJSIfo/sifdCmZyXLG0VQHljkLcKhYsWgAn9br9YTWrpEQPIRs==55cZ"
const sigConvertGPG = `-----BEGIN PGP SIGNATURE-----
Version: Quanto

iQIzBAEBCgAdFiEEab8JRxWM7/xGsOGsRxsIMMDGp/EFAlw3bgwACgkQRxsIMMDG
p/Gq+hAAooiGdBZl0z1+uZs6voUEPloIl0qYxSuDdgI2QAdTiALcbasuzhYge04e
xIgpXf6Exik3TH4Qop5RqpvbDRK5J5AYvWdst377NSIL/m00X44hU3Mq3oJ52LyT
Cj3qShMDkviXtm7GynoXNFaloPwxs3hXze3E+ddWVn17Nw9tIAJbdeWOMRbWSdpi
jAsOZP6qGvrjejNCA3eQSTb2G15zB69yS///mgeRVLNGC7YHzbgX3VROXix6pcdc
8LOgZolloey7VkrOkvBg7t9n2VpqMti1qUQ3qGVLx27YyKjjI+mykUnoO2i5KzsM
fZVCB9iQC3FgVmaGElLUxVJGGToByw4QNuTsLNeVchd+nA20dhQmmZ2dmaMpUIOl
0TbL3wxPxa7eJ72fx3+6EQIqQw0t6ScauPfEQ7Ad0ORIEhGvRXhNYykNUVgdoH09
FoF1eEZv2yvJK5UDQNDUifTnhJ+7A1r7jgykE3vqcrcegbJahC0Qjn66316+D1O/
6I5E/ZZtx3zuzJQT9kTawDTslnmgg5XhQ9LmsrjBYpSKNspAvlhonue07XVyekO1
u6UaKTOmGG060dInWby5Xf+YAK7W8a7Iucoq3zPM0Y6eMVDMNcGcLWhcyCnnFRhO
rGJSIfo/sifdCmZyXLG0VQHljkLcKhYsWgAn9br9YTWrpEQPIRs=
=55cZ
-----END PGP SIGNATURE-----`

func TestQuanto2GPG(t *testing.T) {
	z := Quanto2GPG(sigConvertQuanto)
	if z != sigConvertGPG {
		t.Errorf("Expected %s got %s", sigConvertGPG, z)
	}

	z = Quanto2GPG("asdausigheioygase")
	if z != "" {
		t.Errorf("Expected empty got %s", z)
	}
}

func TestGPG2Quanto(t *testing.T) {
	z := GPG2Quanto(sigConvertGPG, "0ADF79401F28C569", "SHA512")
	if z != sigConvertQuanto {
		t.Errorf("Expected %s got %s", sigConvertQuanto, z)
	}
}

func TestGetFingerPrintFromKey(t *testing.T) {
	z, err := ioutil.ReadFile("./tests/testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	k, err := GetFingerPrintFromKey(string(z))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if k != rstest.TestKeyFingerprint {
		t.Errorf("Expected %s got %s", rstest.TestKeyFingerprint, k)
	}

	// Test Error Scenarios
	_, err = GetFingerPrintFromKey("huebr")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	_, err = GetFingerPrintFromKey(sigConvertGPG) // Test Non Key GPG
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetFingerPrintsFromEncryptedMessage(t *testing.T) {
	fps, err := GetFingerPrintsFromEncryptedMessage(rstest.TestDecryptDataAscii)

	if err != nil {
		t.Fatalf("Got error in test: %s", err)
		t.FailNow()
	}

	if len(fps) != 1 {
		t.Fatalf("Expected 1 fingerprint. Got %d", len(fps))
		t.FailNow()
	}

	if fps[0] != "AB8917A1CA8BCF0E" {
		t.Fatalf("Expected AB8917A1CA8BCF0E got %s", fps[0])
	}

	// Try invalid data

	fps, err = GetFingerPrintsFromEncryptedMessage("huebrinvalidpayload")

	if err == nil {
		t.Fatalf("Expected error")
		t.FailNow()
	}

	if fps != nil {
		t.Fatalf("Expected fingerprints to be null")
	}

	// Test Non PGP Data
	fps, err = GetFingerPrintsFromEncryptedMessage(strings.Replace(rstest.TestDecryptDataAscii, "PGP MESSAGE", "HUE MESSAGE", -1))

	if err == nil {
		t.Fatalf("Expected error")
		t.FailNow()
	}

	if fps != nil {
		t.Fatalf("Expected fingerprints to be null")
	}

}

func TestGetFingerPrintsFromEncryptedMessageRaw(t *testing.T) {
	fps, err := GetFingerPrintsFromEncryptedMessageRaw(rstest.TestDecryptDataRawB64)

	if err != nil {
		t.Fatalf("Got error in test: %s", err)
		t.FailNow()
	}

	if len(fps) != 1 {
		t.Fatalf("Expected 1 fingerprint. Got %d", len(fps))
		t.FailNow()
	}

	if fps[0] != "AB8917A1CA8BCF0E" {
		t.Fatalf("Expected AB8917A1CA8BCF0E got %s", fps[0])
	}

	// Try invalid data

	fps, err = GetFingerPrintsFromEncryptedMessageRaw("huebrinvalidpayload")

	if err == nil {
		t.Fatalf("Expected error")
		t.FailNow()
	}

	if fps != nil {
		t.Fatalf("Expected fingerprints to be null")
	}

	fps, err = GetFingerPrintsFromEncryptedMessageRaw(base64.StdEncoding.EncodeToString([]byte("huebrinvalidpayload")))

	if err == nil {
		t.Fatalf("Expected error")
		t.FailNow()
	}

	if fps != nil {
		t.Fatalf("Expected fingerprints to be null")
	}
}

func TestReadKeyToEntity(t *testing.T) {
	z, err := ioutil.ReadFile("./tests/testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e, err := ReadKeyToEntity(string(z))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if IssuerKeyIdToFP16(e.PrimaryKey.KeyId) != rstest.TestKeyFingerprint {
		t.Errorf("Expected %s got %s", rstest.TestKeyFingerprint, IssuerKeyIdToFP16(e.PrimaryKey.KeyId))
	}

	_, err = ReadKeyToEntity("hueheueheuehue")
	if err == nil {
		t.Errorf("Expected error got nil")
	}
}

func TestCompareFingerPrint(t *testing.T) {

	// fpA == ""
	if CompareFingerPrint("", "auisiehuase") {
		t.Error("Expected false got true")
	}

	// fpB == ""
	if CompareFingerPrint("asuieha", "") {
		t.Error("Expected false got true")
	}

	// fpA == "" && fpB == ""
	if CompareFingerPrint("", "") {
		t.Error("Expected false got true")
	}

	if !CompareFingerPrint("ABCDEFHG", "ABCDEFHG") {
		t.Error("Expected true got false")
	}

	// fpA > fpB && true
	if !CompareFingerPrint("1234567890", "4567890") {
		t.Error("Expected true got false")
	}
	// fpA > fpB && false
	if CompareFingerPrint("1234567890", "4569990") {
		t.Error("Expected false got true")
	}

	// fpB > fpA && true
	if !CompareFingerPrint("4567890", "1234567890") {
		t.Error("Expected true got false")
	}
	// fpB > fpA && false
	if CompareFingerPrint("4569990", "1234567890") {
		t.Error("Expected false got true")
	}
}

func TestCrc24(t *testing.T) {
	z := []byte{1, 2, 3, 3, 41, 23, 12, 31, 23, 12, 31, 23, 12, 41, 24, 15, 12, 43, 12, 31, 23, 12, 31, 23, 123, 12, 4, 12, 31, 23, 12, 31, 23, 120}
	o := CRC24(z)
	if o != 8124930 {
		t.Errorf("Expected %d got %d", 8124930, o)
	}
}

func TestCreateEntityForSubKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	var cTimestamp = time.Now()

	pgpPubKey := packet.NewRSAPublicKey(cTimestamp, &privateKey.PublicKey)
	pgpPrivKey := packet.NewRSAPrivateKey(cTimestamp, privateKey)

	e := CreateEntityForSubKey(rstest.TestKeyFingerprint, pgpPubKey, pgpPrivKey)

	if e.PrimaryKey != pgpPubKey {
		t.Errorf("Expected Primary Key to be the Public key")
	}

	if e.PrivateKey != pgpPrivKey {
		t.Errorf("Expected Private Key to be the Private Key")
	}
}

func TestCreateEntityFromKeys(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	var cTimestamp = time.Now()

	pgpPubKey := packet.NewRSAPublicKey(cTimestamp, &privateKey.PublicKey)
	pgpPrivKey := packet.NewRSAPrivateKey(cTimestamp, privateKey)

	e := CreateEntityFromKeys("huebr", "comment", "a@a.com", 0, pgpPubKey, pgpPrivKey)

	if e.PrimaryKey != pgpPubKey {
		t.Errorf("Expected Primary Key to be the Public key")
	}

	if e.PrivateKey != pgpPrivKey {
		t.Errorf("Expected Private Key to be the Private Key")
	}

	if len(e.Identities) == 0 {
		t.Errorf("Expected one identity")
	}
	fullName := "huebr (comment) <a@a.com>"
	if e.Identities[fullName] != nil {
		id := e.Identities[fullName]
		if id.Name != "huebr" {
			t.Errorf("Expected identity name to be huebr")
		}
		uid := id.UserId

		if uid != nil {
			if uid.Name != "huebr" {
				t.Errorf("Expected UID.name to be huebr")
			}
			if uid.Email != "a@a.com" {
				t.Errorf("Expected UID.Email to be a@a.com")
			}
			if uid.Comment != "comment" {
				t.Errorf("Expected UID.Comment to be comment")
			}
		} else {
			t.Errorf("Expected Identity to have UID")
		}
	} else {
		t.Errorf("Expected identity called huebr")
	}
}

func TestSignatureFix(t *testing.T) {
	s := SignatureFix(rstest.TestSignatureSignature)

	original := GPG2Quanto(rstest.TestSignatureSignature, "", "")
	fixed := GPG2Quanto(s, "", "")

	if original != fixed {
		t.Errorf("Expected: %s\nGot %s", original, fixed)
	}

	s = SignatureFix(rstest.TestSignatureSignatureNoCRC)
	fixed = GPG2Quanto(s, "", "")

	if original != fixed {
		t.Errorf("Expected: %s\nGot %s", original, fixed)
	}

	s = SignatureFix(rstest.TestSignatureSignatureNoCRCSingleLine)
	fixed = GPG2Quanto(s, "", "")

	if original != fixed {
		t.Errorf("Expected: %s\nGot %s", original, fixed)
	}

	// Test invalid base64
	assertPanic(t, func() {
		SignatureFix(strings.Replace(rstest.TestSignatureSignatureNoCRC, "wsFcBAA", "iQ-----", -1))
	}, "Expected panic on invalid base64")

	s = SignatureFix(rstest.BrokenMacOSXSignature)
	fixed = GPG2Quanto(s, "", "")

	if original != fixed {
		t.Errorf("Expected: %s\nGot %s", original, fixed)
	}

	// Test Embedded CRC Case
	s = SignatureFix(rstest.TestEmbeddedCRCSignature)

	//          Try get the fingerprint of signature
	b := bytes.NewReader([]byte(s))
	block, err := armor.Decode(b)
	if err != nil {
		t.Errorf("Failed to read Armored format from TestEmbeddedCRCSignature")
		t.FailNow()
	}

	if block.Type != openpgp.SignatureType {
		t.Errorf("TestEmbeddedCRCSignature does not have signature type")
		t.FailNow()
	}
	fingerPrint := ""
	reader := packet.NewReader(block.Body)
	for {
		pkt, err := reader.Next()

		if err != nil {
			break
		}

		switch sig := pkt.(type) {
		case *packet.Signature:
			fingerPrint = IssuerKeyIdToFP16(*sig.IssuerKeyId)
		case *packet.SignatureV3:
			fingerPrint = IssuerKeyIdToFP16(sig.IssuerKeyId)
		}

		if fingerPrint != "" {
			break
		}
	}

	if fingerPrint == "" {
		t.Errorf("Failed to read Armored format from TestEmbeddedCRCSignature")
		t.FailNow()
	}
}

func TestSimpleIdentitiesToString(t *testing.T) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	var cTimestamp = time.Now()

	pgpPubKey := packet.NewRSAPublicKey(cTimestamp, &privateKey.PublicKey)
	pgpPrivKey := packet.NewRSAPrivateKey(cTimestamp, privateKey)

	e := CreateEntityFromKeys("huebr", "comment", "a@a.com", 0, pgpPubKey, pgpPrivKey)

	ids := IdentityMapToArray(e.Identities)
	if len(ids) != 1 {
		t.Fatalf("Expected one ID got %d", len(ids))
		t.FailNow()
	}

	idsString := SimpleIdentitiesToString(ids)

	if idsString != "huebr" {
		t.Fatalf("Expected idsString to be huebr got %s", idsString)
	}
}

func TestCopyFiles(t *testing.T) {
	folderA, _ := ioutil.TempDir("/tmp", "")
	folderB, _ := ioutil.TempDir("/tmp", "")
	folderAFiles := make([]string, 0)

	for i := 0; i < 4; i++ {
		f, _ := ioutil.TempFile(folderA, "")
		folderAFiles = append(folderAFiles, path.Base(f.Name()))
		_, _ = f.WriteString("Test")
		_ = f.Close()
	}

	err := CopyFiles(folderA, folderB)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		t.FailNow()
	}

	files, err := ioutil.ReadDir(folderB)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		t.FailNow()
	}

	folderBFiles := make([]string, 0)

	for _, f := range files {
		if !f.IsDir() {
			folderBFiles = append(folderBFiles, f.Name())
			found := false
			for _, f2 := range folderAFiles {
				if f2 == f.Name() {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Cannot find %s in folderA", f.Name())
			}
		}
	}

	for _, v := range folderAFiles {
		found := false
		for _, v2 := range folderBFiles {
			if v2 == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Cannot find %s in folderB", v)
		}
	}
}

func TestBrokenMacOSXKey(t *testing.T) {
	s := strings.Split(rstest.BrokenMacOSXSignature, "\n")
	s = brokenMacOSXArrayFix(s, true)

	fixed := strings.Join(s, "\n")
	if fixed != rstest.BrokenMacOSXSignatureFixed {
		t.Errorf("macosx signature not fixed. Expected:\n%s\nGot:\n%s", rstest.BrokenMacOSXSignatureFixed, fixed)
	}
}

func TestOneLineSignature(t *testing.T) {
	sig := SignatureFix(rstest.OneLineSignature)
	if sig != rstest.BrokenMacOSXSignatureFixed {
		t.Errorf("expected one line signature to be fixed. Expected:\n%s\nGot:\n%s", rstest.BrokenMacOSXSignatureFixed, sig)
	}
}

func TestGetFingerPrintsFromKey(t *testing.T) {
	fps, err := GetFingerPrintsFromKey(rstest.TestPublicKeyManySubkeys)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range fps {
		if StringIndexOf(v, rstest.SubKeysFromTestPublicKeyManySubkeys) == -1 {
			t.Errorf("expected %s to be in fingerprints list", v)
		}
	}
}

func TestFolderExists(t *testing.T) {
	v := FolderExists("tests")
	if !v {
		t.Errorf("expected FolderExists(\"tests\") == true")
	}

	v = FolderExists("__heuerbabsueaius31i2u3n13ubae___")
	if v {
		t.Errorf("expected FolderExists(\"__heuerbabsueaius31i2u3n13ubae___\") == false")
	}

	v = FolderExists("tools_test.go")
	if v {
		t.Errorf("expected FolderExists(\"tools_test.go\") == false")
	}
}

func TestGeneratePassword(t *testing.T) {
	b := GeneratePassword()
	if len(b) != defaultPasswordLength {
		t.Errorf("expected password to be %d bytes long", defaultPasswordLength)
	}

	for _, v := range b {
		if strings.Index(passwordBytes, string(v)) == -1 {
			t.Errorf("char %s is not in passwordBytes list.", string(v))
		}
	}
}

func TestIsASCIIArmored(t *testing.T) {
	b, err := ioutil.ReadFile("tests/nonasciiencrypted.gpg")
	if err != nil {
		t.Errorf("Error loading file: %s", err)
	}

	nonascii := string(b)

	if IsASCIIArmored(nonascii) != false {
		t.Errorf("Expected NONASCII from tests/nonasciiencrypted.gpg")
	}

	if IsASCIIArmored(sigConvertGPG) != true {
		t.Errorf("Expected ASCII from sigConvertGPG")
	}
}

func TestNonASCIIFingerprints(t *testing.T) {
	b, err := ioutil.ReadFile("tests/nonasciiencrypted.gpg")
	if err != nil {
		t.Errorf("Error loading file: %s", err)
	}

	nonascii := string(b)

	fps, err := GetFingerPrintsFromEncryptedMessage(nonascii)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(fps) != 1 {
		t.Fatalf("Expected one fingerprint on encrypted got %d", len(fps))
	}

	if fps[0] != "344C911D5CA6B681" {
		t.Fatalf("Expected fingerprint to be 344C911D5CA6B681")
	}
}
