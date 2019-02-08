package database

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/SLog"
	"gopkg.in/rethinkdb/rethinkdb-go.v5"
	"os"
	"testing"
	"time"
)

func ResetDatabase() {
	dbLog.Info("Reseting Database")
	c := GetConnection()
	dbs := GetDatabases()
	dbLog.Info("Dropping test database %s", remote_signer.DatabaseName)
	if remote_signer.StringIndexOf(remote_signer.DatabaseName, dbs) > -1 {
		dbLog.Info("Test Database already exists, dropping.")
		_, _ = rethinkdb.DBDrop(remote_signer.DatabaseName).RunWrite(c)
	}
	WaitDatabaseDrop(remote_signer.DatabaseName)
	time.Sleep(5 * time.Second)
	dbLog.Info("Database reseted")
}

func TestMain(m *testing.M) {
	QuantoError.EnableStackTrace()
	SLog.SetTestMode()

	remote_signer.PrivateKeyFolder = ".."
	remote_signer.KeyPrefix = "testkey_"
	remote_signer.KeysBase64Encoded = false
	remote_signer.RethinkDBPoolSize = 1

	remote_signer.MasterGPGKeyBase64Encoded = false
	remote_signer.MasterGPGKeyPath = "../testkey_privateTestKey.gpg"
	remote_signer.MasterGPGKeyPasswordPath = "../testprivatekeyPassword.txt"

	remote_signer.DatabaseName = "qrs_test"
	remote_signer.HttpPort = 40000
	remote_signer.SKSServer = fmt.Sprintf("http://localhost:%d/sks/", remote_signer.HttpPort)
	remote_signer.EnableRethinkSKS = true

	code := m.Run()
	ResetDatabase()
	os.Exit(code)
}

func TestInitTable(t *testing.T) {
	SLog.UnsetTestMode()
	DbSetup()
	ResetDatabase()
	InitTables()
}
