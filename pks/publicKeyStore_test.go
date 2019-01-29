package remote_signer

import (
	"github.com/quan-to/remote-signer/models"
	"io/ioutil"
	"testing"
)

func TestPKSGetKey(t *testing.T) {
	pushVariables()
	defer popVariables()

	// Test Internal
	c := GetConnection()

	z, err := ioutil.ReadFile("testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	gpgKey := models.AsciiArmored2GPGKey(string(z))

	_, _, err = models.AddGPGKey(c, gpgKey)
	if err != nil {
		t.Errorf("Fail to add key to database: %s", err)
		t.FailNow()
	}

	key := PKSGetKey(gpgKey.FullFingerPrint)

	fp, err := GetFingerPrintFromKey(key)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !CompareFingerPrint(gpgKey.FullFingerPrint, fp) {
		t.Errorf("Expected %s got %s", gpgKey.FullFingerPrint, fp)
	}

	// Test External
	EnableRethinkSKS = false
	SKSServer = "https://keyserver.ubuntu.com/"

	key = PKSGetKey(externalKeyFingerprint)

	fp, err = GetFingerPrintFromKey(key)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !CompareFingerPrint(externalKeyFingerprint, fp) {
		t.Errorf("Expected %s got %s", externalKeyFingerprint, fp)
	}
}

func TestPKSSearchByName(t *testing.T) {
	pushVariables()
	defer popVariables()

	// Test Panics
	EnableRethinkSKS = false
	assertPanic(t, func() {
		_ = PKSSearchByName("", 0, 1)
	}, "SearchByName without RethinkSKS Should panic!")
}

func TestPKSSearchByFingerPrint(t *testing.T) {
	pushVariables()
	defer popVariables()

	// Test Panics
	EnableRethinkSKS = false
	assertPanic(t, func() {
		_ = PKSSearchByFingerPrint("", 0, 1)
	}, "SearchByFingerPrint without RethinkSKS Should panic!")
}

func TestPKSSearchByEmail(t *testing.T) {
	pushVariables()
	defer popVariables()

	// Test Panics
	EnableRethinkSKS = false
	assertPanic(t, func() {
		_ = PKSSearchByEmail("", 0, 1)
	}, "SearchByEmail without RethinkSKS Should panic!")
}

func TestPKSSearch(t *testing.T) {
	// TODO: Implement method and test
	// For now, should always panic

	assertPanic(t, func() {
		_ = PKSSearch("", 0, 1)
	}, "Search should always panic (NOT IMPLEMENTED)")
}

func TestPKSAdd(t *testing.T) {
	pushVariables()
	defer popVariables()
	// Test Internal
	z, err := ioutil.ReadFile("testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fp, err := GetFingerPrintFromKey(string(z))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	o := PKSAdd(string(z))

	if o != "OK" {
		t.Errorf("Expected %s got %s", "OK", o)
	}

	p := PKSGetKey(fp)

	if p == "" {
		t.Errorf("Key was not found")
		t.FailNow()
	}

	fp2, err := GetFingerPrintFromKey(string(p))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !CompareFingerPrint(fp, fp2) {
		t.Errorf("FingerPrint does not match. Expected %s got %s", fp, fp2)
	}

	// Test External
	EnableRethinkSKS = false
	// TODO: How to be a good test without stuffying SKS?
}
