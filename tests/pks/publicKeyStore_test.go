package pks

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/database"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/etc/pgpBuilder"
	"github.com/quan-to/remote-signer/etc/smBuilder"
	"github.com/quan-to/remote-signer/models"
	"github.com/quan-to/remote-signer/pks"
	"gopkg.in/rethinkdb/rethinkdb-go.v5"
	"io/ioutil"
	"os"
	"testing"
)

func ResetDatabase() {
	c := etc.GetConnection()
	dbs := etc.GetDatabases()
	if remote_signer.StringIndexOf(remote_signer.DatabaseName, dbs) > -1 {
		_, err := rethinkdb.DBDrop(remote_signer.DatabaseName).Run(c)
		if err != nil {
			panic(err)
		}
	}
}

func TestMain(m *testing.M) {
	QuantoError.EnableStackTrace()
	SLog.SetTestMode()

	remote_signer.PrivateKeyFolder = ".."
	remote_signer.KeyPrefix = "testkey_"
	remote_signer.KeysBase64Encoded = false

	remote_signer.MasterGPGKeyBase64Encoded = false
	remote_signer.MasterGPGKeyPath = "../testkey_privateTestKey.gpg"
	remote_signer.MasterGPGKeyPasswordPath = "../testprivatekeyPassword.txt"

	remote_signer.DatabaseName = "qrs_test"
	remote_signer.HttpPort = 40000
	remote_signer.SKSServer = fmt.Sprintf("http://localhost:%d/sks/", remote_signer.HttpPort)
	remote_signer.EnableRethinkSKS = true

	ResetDatabase()
	etc.DbSetup()
	etc.InitTables()

	sm := smBuilder.MakeSM()
	gpg := pgpBuilder.MakePGP()
	gpg.LoadKeys()

	err := gpg.UnlockKey(remote_signer.TestKeyFingerprint, remote_signer.TestKeyPassword)

	if err != nil {
		SLog.SetError(true)
		os.Exit(1)
	}

	_ = sm

	code := m.Run()
	ResetDatabase()
	os.Exit(code)
}

func TestPKSGetKey(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Internal
	c := database.GetConnection()

	z, err := ioutil.ReadFile("../testkey_privateTestKey.gpg")
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

	key := pks.PKSGetKey(gpgKey.FullFingerPrint)

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

	key = pks.PKSGetKey(remote_signer.ExternalKeyFingerprint)

	fp, err = remote_signer.GetFingerPrintFromKey(key)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !remote_signer.CompareFingerPrint(remote_signer.ExternalKeyFingerprint, fp) {
		t.Errorf("Expected %s got %s", remote_signer.ExternalKeyFingerprint, fp)
	}
}

func TestPKSSearchByName(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Panics
	remote_signer.EnableRethinkSKS = false
	assertPanic(t, func() {
		_ = pks.PKSSearchByName("", 0, 1)
	}, "SearchByName without RethinkSKS Should panic!")
}

func TestPKSSearchByFingerPrint(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Panics
	remote_signer.EnableRethinkSKS = false
	assertPanic(t, func() {
		_ = pks.PKSSearchByFingerPrint("", 0, 1)
	}, "SearchByFingerPrint without RethinkSKS Should panic!")
}

func TestPKSSearchByEmail(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Panics
	remote_signer.EnableRethinkSKS = false
	assertPanic(t, func() {
		_ = pks.PKSSearchByEmail("", 0, 1)
	}, "SearchByEmail without RethinkSKS Should panic!")
}

func TestPKSSearch(t *testing.T) {
	// TODO: Implement method and test
	// For now, should always panic

	assertPanic(t, func() {
		_ = pks.PKSSearch("", 0, 1)
	}, "Search should always panic (NOT IMPLEMENTED)")
}

func TestPKSAdd(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()
	// Test Internal
	z, err := ioutil.ReadFile("../testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fp, err := remote_signer.GetFingerPrintFromKey(string(z))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	o := pks.PKSAdd(string(z))

	if o != "OK" {
		t.Errorf("Expected %s got %s", "OK", o)
	}

	p := pks.PKSGetKey(fp)

	if p == "" {
		t.Errorf("Key was not found")
		t.FailNow()
	}

	fp2, err := remote_signer.GetFingerPrintFromKey(string(p))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !remote_signer.CompareFingerPrint(fp, fp2) {
		t.Errorf("FingerPrint does not match. Expected %s got %s", fp, fp2)
	}

	// Test External
	remote_signer.EnableRethinkSKS = false
	// TODO: How to be a good test without stuffying SKS?
}

func assertPanic(t *testing.T, f func(), message string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
		}
	}()
	f()
}
