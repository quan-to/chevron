package chevronlib

import (
	"encoding/base64"
	"testing"

	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/internal/tools"
)

func TestGenerateKey(t *testing.T) {
	// Test Generate Key < MinBits
	_, err := GenerateKey(testKeyPassword, "", keymagic.MinKeyBits-1)
	if err == nil {
		t.Errorf("Expected GenerateKey with bits < %d to fail", keymagic.MinKeyBits)
	}

	// Test Generate Key = MinKeyBits
	result, err := GenerateKey(testKeyPassword, "", keymagic.MinKeyBits)
	if err != nil {
		t.Errorf("Expected GenerateKey with bits = %d to generate a key. Got %q", keymagic.MinKeyBits, err)
	}

	if len(result) == 0 {
		t.Error("Expected a generated key, got empty")
	}
}

func TestLoadKey(t *testing.T) {
	l, err := LoadKey(testKey)
	if err != nil {
		t.Errorf("Expected key to be loaded. But got %q", err)
	}

	if l == 0 {
		t.Error("Expected one private key to be loaded but got 0")
	}

	_, err = LoadKey("huebr")
	if err == nil {
		t.Error("Expected \"huebr\" to trigger a key load error but got none")
	}
}

func TestVerifySignature(t *testing.T) {
	_, _ = LoadKey(testKey)

	res, err := VerifySignature([]byte(payloadToSign), testSignature)

	if err != nil {
		t.Errorf("Expected signature to be valid but got %q", err)
	}

	if !res {
		t.Error("Expected signature go be valid but got false")
	}

	res, err = VerifySignature([]byte(payloadToSign+"BLA"), testSignature)

	if err == nil {
		t.Errorf("Expected signature to have an error but got nil")
	}

	if res {
		t.Error("Expected signature go be invalid but got true")
	}
}

func TestVerifyBase64DataSignature(t *testing.T) {
	_, _ = LoadKey(testKey)

	res, err := VerifyBase64DataSignature(base64.StdEncoding.EncodeToString([]byte(payloadToSign)), testSignature)

	if err != nil {
		t.Errorf("Expected signature to be valid but got %q", err)
	}

	if !res {
		t.Error("Expected signature go be valid but got false")
	}

	res, err = VerifyBase64DataSignature(base64.StdEncoding.EncodeToString([]byte(payloadToSign+"BLA")), testSignature)

	if err == nil {
		t.Errorf("Expected signature to have an error but got nil")
	}

	if res {
		t.Error("Expected signature go be invalid but got true")
	}

}

func TestQuantoVerifySignature(t *testing.T) {
	_, _ = LoadKey(testKey)

	sig := tools.GPG2Quanto(testSignature, testKeyFingerprint, "SHA512")

	res, err := QuantoVerifySignature([]byte(payloadToSign), sig)

	if err != nil {
		t.Errorf("Expected signature to be valid but got %q", err)
	}

	if !res {
		t.Error("Expected signature go be valid but got false")
	}

	res, err = QuantoVerifySignature([]byte(payloadToSign+"BLA"), sig)

	if err == nil {
		t.Errorf("Expected signature to have an error but got nil")
	}

	if res {
		t.Error("Expected signature go be invalid but got true")
	}
}

func TestQuantoVerifyBase64DataSignature(t *testing.T) {
	_, _ = LoadKey(testKey)

	sig := tools.GPG2Quanto(testSignature, testKeyFingerprint, "SHA512")

	res, err := QuantoVerifyBase64DataSignature(base64.StdEncoding.EncodeToString([]byte(payloadToSign)), sig)

	if err != nil {
		t.Errorf("Expected signature to be valid but got %q", err)
	}

	if !res {
		t.Error("Expected signature go be valid but got false")
	}

	res, err = QuantoVerifyBase64DataSignature(base64.StdEncoding.EncodeToString([]byte(payloadToSign+"BLA")), sig)

	if err == nil {
		t.Errorf("Expected signature to have an error but got nil")
	}

	if res {
		t.Error("Expected signature go be invalid but got true")
	}
}

func TestUnlockKey(t *testing.T) {
	_, _ = LoadKey(testKey)

	err := UnlockKey(testKeyFingerprint, testKeyPassword)

	if err != nil {
		t.Errorf("Expected key to be unlocked but got %q instead", err)
	}

	err = UnlockKey(testKeyFingerprint, "ABCD1029371n2cy39812y381jx")

	if err == nil {
		t.Errorf("Expected key unlock to return error but got nil")
	}
}

func TestSignData(t *testing.T) {
	_, _ = LoadKey(testKey)
	_ = UnlockKey(testKeyFingerprint, testKeyPassword)

	result, err := SignData([]byte(payloadToSign), testKeyFingerprint)

	if err != nil {
		t.Errorf("Expected signature to work but got %q", err)
	}

	valid, err := VerifySignature([]byte(payloadToSign), result)

	if err != nil {
		t.Errorf("Error validating signature: %q", err)
	}

	if !valid {
		t.Error("Expected signature to be valid, but got false")
	}
}

func TestSignBase64Data(t *testing.T) {
	_, _ = LoadKey(testKey)
	_ = UnlockKey(testKeyFingerprint, testKeyPassword)

	result, err := SignBase64Data(base64.StdEncoding.EncodeToString([]byte(payloadToSign)), testKeyFingerprint)

	if err != nil {
		t.Errorf("Expected signature to work but got %q", err)
	}

	valid, err := VerifySignature([]byte(payloadToSign), result)

	if err != nil {
		t.Errorf("Error validating signature: %q", err)
	}

	if !valid {
		t.Error("Expected signature to be valid, but got false")
	}
}

func TestQuantoSignData(t *testing.T) {
	_, _ = LoadKey(testKey)
	_ = UnlockKey(testKeyFingerprint, testKeyPassword)

	result, err := QuantoSignData([]byte(payloadToSign), testKeyFingerprint)

	if err != nil {
		t.Errorf("Expected signature to work but got %q", err)
	}

	valid, err := QuantoVerifySignature([]byte(payloadToSign), result)

	if err != nil {
		t.Errorf("Error validating signature: %q", err)
	}

	if !valid {
		t.Error("Expected signature to be valid, but got false")
	}
}

func TestQuantoSignBase64Data(t *testing.T) {
	_, _ = LoadKey(testKey)
	_ = UnlockKey(testKeyFingerprint, testKeyPassword)

	result, err := QuantoSignBase64Data(base64.StdEncoding.EncodeToString([]byte(payloadToSign)), testKeyFingerprint)

	if err != nil {
		t.Errorf("Expected signature to work but got %q", err)
	}

	valid, err := QuantoVerifySignature([]byte(payloadToSign), result)

	if err != nil {
		t.Errorf("Error validating signature: %q", err)
	}

	if !valid {
		t.Error("Expected signature to be valid, but got false")
	}
}

func TestGetPublicKey(t *testing.T) {
	_, _ = LoadKey(testKey)
	pubKey, err := GetPublicKey(testKeyFingerprint)

	if err != nil {
		t.Errorf("Expected public key but got error %q", err)
	}

	fps, err := GetKeyFingerprints(pubKey)

	if err != nil {
		t.Errorf("Expected public key to be valid got error %q", err)
	}

	if len(fps) == 0 {
		t.Errorf("Got no fingerprints, expected one")
	} else if fps[0] != testKeyFingerprint {
		t.Errorf("Expected fingerprint to be %s but got %s", testKeyFingerprint, fps[0])
	}
}

func TestChangeKeyPassword(t *testing.T) {
	const tmpPass = "anei1he9m298em1xh"
	key, _ := GenerateKey(tmpPass, "ACD123", 2048)

	newKey, err := ChangeKeyPassword(key, tmpPass, testKeyPassword)

	if err != nil {
		t.Fatalf("Unexpected error changing key password %q", err)
	}

	fps, _ := GetKeyFingerprints(newKey)

	_, _ = LoadKey(newKey)

	err = UnlockKey(fps[0], testKeyPassword)

	if err != nil {
		t.Fatalf("Unexpected error when unlocking key with new password: %q", err)
	}
}

const payloadToSign = "HUEBR"

const testSignature = `-----BEGIN PGP SIGNATURE-----

wsBcBAABCgAQBQJedt7dCRDORQOylH4iAgAArvYIAMY9G3oQA0ZL7CwZmhjRfu6d
BfNst48oHJ3mSLe1oTATkkiRvJBgXaNg9aFP0rNa5Cx8Hlkhjqp0QYRmPNhKjWhD
L+LpxLasESQeWLVZnDkaz75YHw6UK6TRIsAH7Py/zgiGlNehWx3+0rjkBxEEjJgK
aM5wBDoPi4+1yAf1ZSanypb2fq3Mzy77HKtew+9BCKhC8DE9uXN6PGC4Q/LhgDON
aXS7WqLoIaB9RjvTDHbR4V8rP5pEx5s/GhJHgYAYzcJvhxog3n03GAIW3DiQdbih
h72Ll0XnpYkiiSvwz7UtfeveDSNFUSfeoOKdGsM/8rrZIexHWSJVdUc0I53ITQY=
=K6nc
-----END PGP SIGNATURE-----`

const testKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: GnuPG v2
Comment: Generated by Chevron

xcL0BF52btoBCADtWsLWwS0AugQ98LnVnndA8WPafCkfUWtikBr1CyjXeDP4rqOz
2ovhancM1XDo025AMesfQCB38GFSWEaDJcFG5R1fnRdVqsl0G/+yaSWuEmqP6C7a
Jux//xUpe7jMc8O341SQ1P+FEJv7J8nHi8hXbweuS/mxn7Yiu3rmKzXTKBPJitkF
k/kU/TQcxuoN5LGrU87DEFYFNr1FQIR2q75/8yXhY6HsP6HQkoirVm8AEwo10+9O
6qcW4vs3UnR2/hilxefTOVa3nnTIp/4fGs1xl994olHNYLmIXPm2B3HvJCIHeJUs
yMIZCwMSpB4Q30RRO0LDXgpuQwrgd4WIThkLABEBAAH/BwMCaELZXGd6MaZg2ZUa
NOMX1pSewCSu9/cQy4U0O5kCDc29SFeysi/zvEn/QfuprNd5GgqEVqKXF20HQg5M
a17N4ibJAvw7QFuVAA8fElc2JchurAlfmE6qD0APt/Qo9Xf73P0nfdfvM77d6Xy2
+Pov/+i9O71SOQ72PfQzDy+TgVGAopVKFb3RkpeWivI0UmLnHiDh5ewpbDfbCB3t
F1O4JvdWHsJVe/lIJeHi2GMvLxLFkdBq9fnrgLjSB7aZ09XDqcjJqQmgXwC7bQwR
y4HKozBKAp/g4tgPM6S68ua+1K4s+aiKjjlgK3W+Yi8A0IBzPxYmlsovTn+CmlFo
iB9f+yufuMkIe64D1cHyxIPa+XrdFrUWv94S61FYPJMkEOlEJQWfe8zAtv0by0wX
26ouf4trG1qX7duYSGxRgzxCPIpmNU2JQaZZlR0NkXPPQfg9cc02SbePY2Lg0VG8
OAYO14maYrQ4Wr7kxXSbWPm5rbIvjbvTqRzjEOSz47P5vypNpRHIMQifpzmMWz+C
al4HaRKOh0qwUXvHh1cIhsBpKqUOr9mDV0Pa0MA+4OGGMXlhu5c63H+Pu+lt97Fc
yHspLSexxm6M8PSjiDerUaXroqtoJWL9ynM49mIXjq0xerjIII814fDarRQqcyp5
l0/0IgdASBuM2y9yh447WaB3t1WzsfgiqatyAjVw8IYyqxB/mWU6gedsD3AgtZqh
2vn5xQkooi5p3tAPMbJrSWuIv07HAeNAfITyXkOHxf/Ym9eUn+ZTDxWOAHsX9ZAJ
ic/s4QdL4Yd2iahtR61RwTIrQP/l5pprNftjaX3HWKZMjWxRfz7yUA8hGYhn6ivS
4T0+9zMnVcCAdfgfj+8pjXJIfuVJ5bApIu7R3oRZjsE+8WN1gjGXzQtMdWNhcyBU
ZXNrZcLAXwQTAQoAEwUCXnZu2gkQzkUDspR+IgICGw8AAN4cCAAq8QCLTz2xyg29
VpxrpL8K+3p4RcqsG3PBOigXsOgi/UlCkiNceBwWpyANN2FicYbKpsU2yHKn5Hgh
0BlcjhE9V/a8qzFftTnPTZeMDvBKBM9DAIBjdVbNdMDeZQ4XnckE+4wPNKiOA7pE
F1mpuarrDIaxVWESn9BZxsTUBgx7LoLK/MKQX+oATAD96gmmou0K6M1lt2m0Vqva
nTnwxSeghH3/w7eNXmfmpGAKec5TYVLrLFMPc34ELclONT3xuRx1BBs1HUpXs7so
IeWPJX37xC5oFnuh9MnF0lf/M1hwbNrqC0sT7nUXLFX40zO350Sv2kJxasWoZGov
dmwxVNpdxcL0BF52btoBCADtWsLWwS0AugQ98LnVnndA8WPafCkfUWtikBr1CyjX
eDP4rqOz2ovhancM1XDo025AMesfQCB38GFSWEaDJcFG5R1fnRdVqsl0G/+yaSWu
EmqP6C7aJux//xUpe7jMc8O341SQ1P+FEJv7J8nHi8hXbweuS/mxn7Yiu3rmKzXT
KBPJitkFk/kU/TQcxuoN5LGrU87DEFYFNr1FQIR2q75/8yXhY6HsP6HQkoirVm8A
Ewo10+9O6qcW4vs3UnR2/hilxefTOVa3nnTIp/4fGs1xl994olHNYLmIXPm2B3Hv
JCIHeJUsyMIZCwMSpB4Q30RRO0LDXgpuQwrgd4WIThkLABEBAAH/BwMCaELZXGd6
MaZg2ZUaNOMX1pSewCSu9/cQy4U0O5kCDc29SFeysi/zvEn/QfuprNd5GgqEVqKX
F20HQg5Ma17N4ibJAvw7QFuVAA8fElc2JchurAlfmE6qD0APt/Qo9Xf73P0nfdfv
M77d6Xy2+Pov/+i9O71SOQ72PfQzDy+TgVGAopVKFb3RkpeWivI0UmLnHiDh5ewp
bDfbCB3tF1O4JvdWHsJVe/lIJeHi2GMvLxLFkdBq9fnrgLjSB7aZ09XDqcjJqQmg
XwC7bQwRy4HKozBKAp/g4tgPM6S68ua+1K4s+aiKjjlgK3W+Yi8A0IBzPxYmlsov
Tn+CmlFoiB9f+yufuMkIe64D1cHyxIPa+XrdFrUWv94S61FYPJMkEOlEJQWfe8zA
tv0by0wX26ouf4trG1qX7duYSGxRgzxCPIpmNU2JQaZZlR0NkXPPQfg9cc02SbeP
Y2Lg0VG8OAYO14maYrQ4Wr7kxXSbWPm5rbIvjbvTqRzjEOSz47P5vypNpRHIMQif
pzmMWz+Cal4HaRKOh0qwUXvHh1cIhsBpKqUOr9mDV0Pa0MA+4OGGMXlhu5c63H+P
u+lt97FcyHspLSexxm6M8PSjiDerUaXroqtoJWL9ynM49mIXjq0xerjIII814fDa
rRQqcyp5l0/0IgdASBuM2y9yh447WaB3t1WzsfgiqatyAjVw8IYyqxB/mWU6geds
D3AgtZqh2vn5xQkooi5p3tAPMbJrSWuIv07HAeNAfITyXkOHxf/Ym9eUn+ZTDxWO
AHsX9ZAJic/s4QdL4Yd2iahtR61RwTIrQP/l5pprNftjaX3HWKZMjWxRfz7yUA8h
GYhn6ivS4T0+9zMnVcCAdfgfj+8pjXJIfuVJ5bApIu7R3oRZjsE+8WN1gjGXwsBi
BBgBCgAWBQJedm7aCRDORQOylH4iAgIbDwIVCgAACWMIAClFP9hMxTwPMl2L6QSr
cuFO2ujheqgeMwna8wNp/LwUybtVEQ5aZ67K0knAb1s64q1IYkFbWMJIGEd5hA+B
AHnqGT+7etcCYA24M3iJoKaBZPgdUkMGGndp8hgh7M/8aJmYal3jTNus+sLRWRFU
sslFLkdwa0XtSFh2Wh+9dXqi55N9vRH05aBWEMtlr2qt1X0Y0G3VAySrXUfrG9Yv
QD1eq2AWBXjYosjwDaUR2+v7kAN/euG1Bvs2Vd1/+mfSn2WwBmHhQHDhdWW0oGTp
pq29HT8vVfshyd7NLjDGf4Acifgpv2uIOOKw8QTyi5fMsPr5EGFWgvUnj7WUZzMe
bBM=
=36ik
-----END PGP PRIVATE KEY BLOCK-----`

const testKeyPassword = `1234567890`
const testKeyFingerprint = `CE4503B2947E2202`
