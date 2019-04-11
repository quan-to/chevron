package keymagic

import (
	"fmt"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/QuantoError"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/vaultManager"
	"github.com/quan-to/slog"
	"os"
	"os/exec"
	"testing"
)

var testData = []byte(remote_signer.TestSignatureData)

var pgpMan *PGPManager
var sm *SecretsManager

func TestMain(m *testing.M) {
	slog.UnsetTestMode()
	var rql *exec.Cmd
	var err error
	rql, err = remote_signer.RQLStart()
	if err != nil {
		slog.Error(err)
		os.Exit(1)
	}

	QuantoError.EnableStackTrace()
	slog.SetTestMode()

	remote_signer.DatabaseName = "qrs_test"
	remote_signer.PrivateKeyFolder = "../tests/"
	remote_signer.KeyPrefix = "testkey_"
	remote_signer.KeysBase64Encoded = false

	remote_signer.MasterGPGKeyBase64Encoded = false
	remote_signer.MasterGPGKeyPath = "../tests/testkey_privateTestKey.gpg"
	remote_signer.MasterGPGKeyPasswordPath = "../tests/testprivatekeyPassword.txt"

	remote_signer.HttpPort = 40000
	remote_signer.SKSServer = fmt.Sprintf("http://localhost:%d/sks/", remote_signer.HttpPort)
	remote_signer.EnableRethinkSKS = true
	remote_signer.PushVariables()

	slog.UnsetTestMode()
	etc.DbSetup()
	etc.ResetDatabase()
	etc.InitTables()
	slog.SetTestMode()

	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(remote_signer.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
	}

	pgpMan = MakePGPManagerWithKRM(kb, MakeKeyRingManager()).(*PGPManager)
	pgpMan.LoadKeys()

	sm = MakeSecretsManager()

	err = pgpMan.UnlockKey(remote_signer.TestKeyFingerprint, remote_signer.TestKeyPassword)

	if err != nil {
		slog.SetError(true)
		slog.Error(err)
		os.Exit(1)
	}

	code := m.Run()
	etc.ResetDatabase()
	slog.UnsetTestMode()
	etc.Cleanup()
	slog.Warn("STOPPING RETHINKDB")
	remote_signer.RQLStop(rql)
	os.Exit(code)
}
