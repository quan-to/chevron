package database

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/etc/pgpBuilder"
	"github.com/quan-to/remote-signer/etc/smBuilder"
	"gopkg.in/rethinkdb/rethinkdb-go.v5"
	"os"
	"testing"
)

var slog *SLog.Instance

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
	slog = SLog.Scope("TestLog")

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

	etc.DbSetup()
	ResetDatabase()

	sm := smBuilder.MakeSM()
	gpg := pgpBuilder.MakePGP()
	gpg.LoadKeys()

	err := gpg.UnlockKey(remote_signer.TestKeyFingerprint, remote_signer.TestKeyPassword)

	if err != nil {
		SLog.SetError(true)
		slog.Error(err)
		os.Exit(1)
	}

	_ = sm

	code := m.Run()
	ResetDatabase()
	os.Exit(code)
}

func TestInitTable(t *testing.T) {
	etc.DbSetup()
	etc.InitTables()
	etc.InitTables() // Test if already initialized
	ResetDatabase()
}
