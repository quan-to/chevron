package keymagic

import (
	"fmt"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/QuantoError"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/rstest"
	"github.com/quan-to/chevron/vaultManager"
	"github.com/quan-to/slog"
	"os"
	"os/exec"
	"testing"
)

var testData = []byte(rstest.TestSignatureData)

var pgpMan *PGPManager
var sm *SecretsManager

func TestMain(m *testing.M) {
	slog.UnsetTestMode()
	var rql *exec.Cmd
	var err error
	var port int
	rql, port, err = rstest.RQLStart()
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
	remote_signer.RethinkDBPort = port
	remote_signer.PushVariables()

	slog.UnsetTestMode()
	etc.DbSetup()
	etc.ResetDatabase()
	etc.InitTables()
	slog.SetTestMode()

	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(nil, remote_signer.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(nil, remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
	}

	pgpMan = MakePGPManagerWithKRM(nil, kb, MakeKeyRingManager(nil)).(*PGPManager)
	pgpMan.LoadKeys()

	sm = MakeSecretsManager(nil)

	err = pgpMan.UnlockKey(rstest.TestKeyFingerprint, rstest.TestKeyPassword)

	if err != nil {
		slog.SetError(true)
		slog.Error(err)
		rstest.RQLStop(rql)
		os.Exit(1)
	}

	code := m.Run()
	etc.ResetDatabase()
	slog.UnsetTestMode()
	etc.Cleanup()
	rstest.RQLStop(rql)
	os.Exit(code)
}
