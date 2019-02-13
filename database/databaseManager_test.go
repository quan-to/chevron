package database

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/SLog"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	QuantoError.EnableStackTrace()
	SLog.UnsetTestMode()

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
	DbSetup()

	SLog.SetTestMode()
	code := m.Run()
	SLog.UnsetTestMode()

	ResetDatabase()
	time.Sleep(120 * time.Second)
	os.Exit(code)
}

func TestInitTable(t *testing.T) {
	//ResetDatabase()
	//time.Sleep(5 * time.Second)
	// Breaks the test due rethink non atomic operations
	InitTables()
}
