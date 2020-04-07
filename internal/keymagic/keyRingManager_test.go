package keymagic

import (
	"context"
	"io/ioutil"
	"testing"

	remote_signer "github.com/quan-to/chevron"
	"github.com/quan-to/chevron/database"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/rstest"
	"github.com/quan-to/chevron/vaultManager"
)

func TestAddKey(t *testing.T) {
	ctx := context.Background()
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()
	remote_signer.MaxKeyRingCache = 10
	krm := MakeKeyRingManager(nil)
	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(nil, remote_signer.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(nil, remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
	}

	gpg := MakePGPManagerWithKRM(nil, kb, krm)

	str, err := gpg.GeneratePGPKey(ctx, "", "", gpg.MinKeyBits())

	if err != nil {
		t.Errorf("Cannot generate test key: %s", err)
		t.FailNow()
	}

	e, err := remote_signer.ReadKeyToEntity(str)

	if err != nil {
		t.Errorf("Error loading test key: %s", err)
		t.FailNow()
	}

	fp := remote_signer.IssuerKeyIdToFP16(e.PrimaryKey.KeyId)

	krm.AddKey(ctx, e, true)

	if !krm.ContainsKey(ctx, fp) {
		t.Error("Cannot find added key")
	}

	if krm.GetKey(ctx, fp) == nil {
		t.Error("Cannot find added key")
	}

	if remote_signer.StringIndexOf(fp, krm.GetFingerPrints(ctx)) != -1 {
		t.Error("Non Erasable Key should be on the fingerPrint list")
	}

	if len(krm.GetCachedKeys(ctx)) != 1 {
		t.Error("The generated key is not cached")
	}

	// Test Ring Cache
	str, err = gpg.GeneratePGPKey(ctx, "", "", gpg.MinKeyBits())

	if err != nil {
		t.Errorf("Cannot generate test key: %s", err)
		t.FailNow()
	}

	erasableKeyTest, err := remote_signer.ReadKeyToEntity(str)

	if err != nil {
		t.Errorf("Error loading test key: %s", err)
		t.FailNow()
	}

	krm.AddKey(ctx, erasableKeyTest, false) // Add to pool
	fpErasable := remote_signer.IssuerKeyIdToFP16(erasableKeyTest.PrimaryKey.KeyId)

	// Rotate until MaxKeyRingCache -1, so erasableKeyTest should be still there
	for i := 0; i < remote_signer.MaxKeyRingCache-1; i++ {
		str, err := gpg.GenerateTestKey()

		if err != nil {
			t.Errorf("Cannot generate test key: %s", err)
			t.FailNow()
		}

		e, err := remote_signer.ReadKeyToEntity(str)

		if err != nil {
			t.Errorf("Error loading test key: %s", err)
			t.FailNow()
		}
		krm.AddKey(ctx, e, false)
	}

	// fpErasable should be still there
	if !krm.ContainsKey(ctx, fpErasable) {
		t.Errorf("For MaxRingCache - 1, fpErasable should be still stored")
	}

	// Generate one more, should be erased
	str, err = gpg.GeneratePGPKey(ctx, "", "", gpg.MinKeyBits())

	if err != nil {
		t.Errorf("Cannot generate test key: %s", err)
		t.FailNow()
	}

	e, err = remote_signer.ReadKeyToEntity(str)

	if err != nil {
		t.Errorf("Error loading test key: %s", err)
		t.FailNow()
	}
	krm.AddKey(ctx, e, false)

	// fpErasable should not be there
	if krm.ContainsKey(ctx, fpErasable) {
		t.Errorf("For MaxRingCache, fpErasable should not be still stored")
	}
	krm.AddKey(ctx, e, false)
	krm.AddKey(ctx, e, false)
}

func TestGetKeyExternal(t *testing.T) {
	ctx := context.Background()
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()
	// Test External SKS Fetch
	remote_signer.SKSServer = "https://keyserver.ubuntu.com/"
	krm := MakeKeyRingManager(nil)
	remote_signer.EnableRethinkSKS = false
	e := krm.GetKey(ctx, rstest.ExternalKeyFingerprint)

	if e == nil {
		t.Error("Expected External key to be fetch")
		t.FailNow()
	}

	fp := remote_signer.IssuerKeyIdToFP16(e.PrimaryKey.KeyId)

	if fp != rstest.ExternalKeyFingerprint {
		t.Errorf("Expected key %s got %s", rstest.ExternalKeyFingerprint, fp)
	}

	// Test SKS Internal
	remote_signer.EnableRethinkSKS = true
	c := database.GetConnection()

	z, err := ioutil.ReadFile("../tests/testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	gpgKey, _ := models.AsciiArmored2GPGKey(string(z))

	_, _, err = models.AddGPGKey(c, gpgKey)
	if err != nil {
		t.Errorf("Fail to add key to database: %s", err)
		t.FailNow()
	}

	e = krm.GetKey(ctx, gpgKey.FullFingerPrint)

	if e == nil {
		t.Error("Expected Internal key to be fetch")
		t.FailNow()
	}

	fp = remote_signer.IssuerKeyIdToFP16(e.PrimaryKey.KeyId)

	if !remote_signer.CompareFingerPrint(fp, gpgKey.FullFingerPrint) {
		t.Errorf("Expected %s == %s", fp, gpgKey.FullFingerPrint)
	}

	//// Test Bad Scenario
	//gpgKey.FullFingerPrint = "ABCDABCDABCDABCDABCDABCDABCDABCD"
	//gpgKey.AsciiArmoredPublicKey = "HUEBR"
	//_, _, err = models.AddGPGKey(c, gpgKey)
	//if err != nil {
	//	t.Errorf("Fail to add key to database: %s", err)
	//	t.FailNow()
	//}
	//
	//e = krm.GetKey(gpgKey.FullFingerPrint)
	//
	//if e != nil {
	//	t.Errorf("Expected key to be null got %v", e)
	//}
}

//func TestBigPublicKeyExponent(t *testing.T) {
//    // Breaks Golang OpenPGP
//    bigPubKeyExponent := `-----BEGIN PGP PUBLIC KEY BLOCK-----
//Version: SKS 1.1.6
//Comment: Hostname: keyserver.ubuntu.com
//
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
//WDhTJw==
//=Nx0F
//-----END PGP PUBLIC KEY BLOCK-----`
//    _, err := ReadKeyToEntity(bigPubKeyExponent)
//
//    if err != nil {
//        t.Error(err)
//    }
//}
