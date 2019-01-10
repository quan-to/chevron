package remote_signer

import (
	"io/ioutil"
	"testing"
)

func TestStringIndexOf(t *testing.T) {
	v := []string{"a", "v", "o"}

	for i := range v {
		idx := stringIndexOf(v[i], v)
		if idx != i {
			t.Errorf("Expected %d got %d", i, idx)
		}
	}

	idx := stringIndexOf("huebrbrb", v)
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
	z, err := ioutil.ReadFile("testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	k, err := GetFingerPrintFromKey(string(z))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if k != testKeyFingerprint {
		t.Errorf("Expected %s got %s", testKeyFingerprint, k)
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

func TestGetFingerPrintsFromEncryptedMessageRaw(t *testing.T) {
	// TODO
}

func TestReadKeyToEntity(t *testing.T) {
	z, err := ioutil.ReadFile("testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e, err := ReadKeyToEntity(string(z))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if IssuerKeyIdToFP16(e.PrimaryKey.KeyId) != testKeyFingerprint {
		t.Errorf("Expected %s got %s", testKeyFingerprint, IssuerKeyIdToFP16(e.PrimaryKey.KeyId))
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
	o := crc24(z)
	if o != 8124930 {
		t.Errorf("Expected %d got %d", 8124930, o)
	}
}
