package keymagic

import (
	"context"
	"fmt"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/database"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/test"
	"io/ioutil"
	"testing"
)

func TestPKSGetKey(t *testing.T) {
	config.PushVariables()
	defer config.PopVariables()

	ctx := context.Background()

	// Test Internal
	c := database.GetConnection()

	z, err := ioutil.ReadFile("../../test/data/testkey_privateTestKey.gpg")
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

	fp, err := tools.GetFingerPrintFromKey(key)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !tools.CompareFingerPrint(gpgKey.FullFingerPrint, fp) {
		t.Errorf("Expected %s got %s", gpgKey.FullFingerPrint, fp)
	}

	// Test External
	config.EnableRethinkSKS = false
	config.SKSServer = "https://keyserver.ubuntu.com/"

	key, _ = PKSGetKey(ctx, test.ExternalKeyFingerprint)

	fp, err = tools.GetFingerPrintFromKey(key)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !tools.CompareFingerPrint(test.ExternalKeyFingerprint, fp) {
		t.Errorf("Expected %s got %s", test.ExternalKeyFingerprint, fp)
	}
}

func TestPKSSearchByName(t *testing.T) {
	config.PushVariables()
	defer config.PopVariables()

	// Test Panics
	config.EnableRethinkSKS = false
	_, err := PKSSearchByName("", 0, 1)
	if err == nil {
		t.Fatalf("Search should fail as not implemented for rethinkdb disabled!")
	}
}

func TestPKSSearchByFingerPrint(t *testing.T) {
	config.PushVariables()
	defer config.PopVariables()

	// Test Panics
	config.EnableRethinkSKS = false
	_, err := PKSSearchByFingerPrint("", 0, 1)
	if err == nil {
		t.Fatalf("Search should fail as not implemented for rethinkdb disabled!")
	}
}

func TestPKSSearchByEmail(t *testing.T) {
	config.PushVariables()
	defer config.PopVariables()

	// Test Panics
	config.EnableRethinkSKS = false
	_, err := PKSSearchByEmail("", 0, 1)
	if err == nil {
		t.Fatalf("Search should fail as not implemented for rethinkdb disabled!")
	}
}

func TestPKSSearch(t *testing.T) {
	// TODO: Implement method and test
	config.EnableRethinkSKS = false
	_, err := PKSSearch("", 0, 1)
	if err == nil {
		t.Fatalf("Search should fail as not implemented for rethinkdb disabled!")
	}
}

func TestPKSAdd(t *testing.T) {
	config.PushVariables()
	defer config.PopVariables()
	config.EnableRethinkSKS = true
	ctx := context.Background()
	// Test Internal
	z, err := ioutil.ReadFile("../../test/data/testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fp, err := tools.GetFingerPrintFromKey(string(z))

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

	fp2, err := tools.GetFingerPrintFromKey(string(p))

	if err != nil {
		t.Error(err)
		t.Error(fmt.Errorf("key data: %s", string(p)))
		t.FailNow()
	}

	if !tools.CompareFingerPrint(fp, fp2) {
		t.Errorf("Fingerprint does not match. Expected %s got %s", fp, fp2)
	}

	// Test External
	config.EnableRethinkSKS = false
	// TODO: How to be a good test without stuffying SKS?
}
