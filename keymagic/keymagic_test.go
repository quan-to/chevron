package keymagic

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/keyBackend"
	"github.com/quan-to/remote-signer/vaultManager"
	"gopkg.in/rethinkdb/rethinkdb-go.v5"
	"os"
	"sync"
	"testing"
	"time"
)

var testData = []byte(remote_signer.TestSignatureData)

var pgpMan etc.PGPInterface

var dbLock = sync.Mutex{}

func ResetDatabase() {
	dbLock.Lock()
	defer dbLock.Unlock()
	c := etc.GetConnection()
	dbs := etc.GetDatabases()
	if remote_signer.StringIndexOf(remote_signer.DatabaseName, dbs) > -1 {
		_, err := rethinkdb.DBDrop(remote_signer.DatabaseName).Run(c)
		if err != nil {
			panic(err)
		}
	}
	time.Sleep(3 * time.Second)
}

func TestMain(m *testing.M) {
	QuantoError.EnableStackTrace()
	SLog.SetTestMode()

	remote_signer.PrivateKeyFolder = "../tests/"
	remote_signer.KeyPrefix = "testkey_"
	remote_signer.KeysBase64Encoded = false

	remote_signer.MasterGPGKeyBase64Encoded = false
	remote_signer.MasterGPGKeyPath = "../tests/testkey_privateTestKey.gpg"
	remote_signer.MasterGPGKeyPasswordPath = "../tests/testprivatekeyPassword.txt"

	remote_signer.DatabaseName = "qrs_test"
	remote_signer.HttpPort = 40000
	remote_signer.SKSServer = fmt.Sprintf("http://localhost:%d/sks/", remote_signer.HttpPort)
	remote_signer.EnableRethinkSKS = true

	ResetDatabase()
	etc.DbSetup()
	etc.InitTables()

	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(remote_signer.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
	}

	pgpMan = MakePGPManagerWithKRM(kb, MakeKeyRingManager())
	pgpMan.LoadKeys()

	err := pgpMan.UnlockKey(remote_signer.TestKeyFingerprint, remote_signer.TestKeyPassword)

	if err != nil {
		SLog.SetError(true)
		SLog.Error(err)
		os.Exit(1)
	}

	code := m.Run()
	ResetDatabase()
	os.Exit(code)
}
