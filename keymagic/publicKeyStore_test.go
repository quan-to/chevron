package keymagic

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	remote_signer "github.com/quan-to/chevron"
	"github.com/quan-to/chevron/database"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/rstest"
)

func TestPKSGetKey(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	ctx := context.Background()

	// Test Internal
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

	key, _ := PKSGetKey(ctx, gpgKey.FullFingerPrint)

	fp, err := remote_signer.GetFingerPrintFromKey(key)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !remote_signer.CompareFingerPrint(gpgKey.FullFingerPrint, fp) {
		t.Errorf("Expected %s got %s", gpgKey.FullFingerPrint, fp)
	}

	// Test External
	remote_signer.EnableRethinkSKS = false
	remote_signer.SKSServer = "https://keyserver.ubuntu.com/"

	key, _ = PKSGetKey(ctx, rstest.ExternalKeyFingerprint)

	fp, err = remote_signer.GetFingerPrintFromKey(key)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !remote_signer.CompareFingerPrint(rstest.ExternalKeyFingerprint, fp) {
		t.Errorf("Expected %s got %s", rstest.ExternalKeyFingerprint, fp)
	}
}

func TestPKSSearchByName(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Panics
	remote_signer.EnableRethinkSKS = false
	_, err := PKSSearchByName("", 0, 1)
	if err == nil {
		t.Fatalf("Search should fail as not implemented for rethinkdb disabled!")
	}
}

func TestPKSSearchByFingerPrint(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Panics
	remote_signer.EnableRethinkSKS = false
	_, err := PKSSearchByFingerPrint("", 0, 1)
	if err == nil {
		t.Fatalf("Search should fail as not implemented for rethinkdb disabled!")
	}
}

func TestPKSSearchByEmail(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Panics
	remote_signer.EnableRethinkSKS = false
	_, err := PKSSearchByEmail("", 0, 1)
	if err == nil {
		t.Fatalf("Search should fail as not implemented for rethinkdb disabled!")
	}
}

func TestPKSSearch(t *testing.T) {
	// TODO: Implement method and test
	remote_signer.EnableRethinkSKS = false
	_, err := PKSSearch("", 0, 1)
	if err == nil {
		t.Fatalf("Search should fail as not implemented for rethinkdb disabled!")
	}
}

func TestPKSAdd(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()
	remote_signer.EnableRethinkSKS = true
	ctx := context.Background()
	// Test Internal
	z, err := ioutil.ReadFile("../tests/testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fp, err := remote_signer.GetFingerPrintFromKey(string(z))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	o := PKSAdd(ctx, string(z))

	if o != "OK" {
		t.Errorf("Expected %s got %s", "OK", o)
	}

	p, _ := PKSGetKey(ctx, fp)

	if p == "" {
		t.Errorf("Key was not found")
		t.FailNow()
	}

	fp2, err := remote_signer.GetFingerPrintFromKey(string(p))

	if err != nil {
		t.Error(err)
		t.Error(fmt.Errorf("key data: %s", string(p)))
		t.FailNow()
	}

	if !remote_signer.CompareFingerPrint(fp, fp2) {
		t.Errorf("FingerPrint does not match. Expected %s got %s", fp, fp2)
	}

	// Test External
	remote_signer.EnableRethinkSKS = false
	// TODO: How to be a good test without stuffying SKS?
}
