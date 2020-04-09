package keymagic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	config "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/etc"
	"github.com/quan-to/chevron/internal/keyBackend"
	"github.com/quan-to/chevron/internal/vaultManager"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/testdata"
	"os"
	"testing"

	"github.com/quan-to/slog"
)

var testData = []byte(testdata.TestSignatureData)

var pgpMan *PGPManager
var sm *SecretsManager

func TestMain(m *testing.M) {
	slog.UnsetTestMode()
	var err error

	QuantoError.EnableStackTrace()
	slog.SetTestMode()
	u, _ := uuid.NewRandom()
	config.DatabaseName = "qrs_test_" + u.String()
	config.PrivateKeyFolder = "../../testdata/"
	config.KeyPrefix = "testkey_"
	config.KeysBase64Encoded = false

	config.MasterGPGKeyBase64Encoded = false
	config.MasterGPGKeyPath = "../../testdata/testkey_privateTestKey.gpg"
	config.MasterGPGKeyPasswordPath = "../../testdata/testprivatekeyPassword.txt"

	config.HttpPort = 40000
	config.SKSServer = fmt.Sprintf("http://localhost:%d/sks/", config.HttpPort)
	config.EnableRethinkSKS = true
	config.PushVariables()

	slog.UnsetTestMode()
	etc.DbSetup()
	etc.ResetDatabase()
	etc.InitTables()
	slog.SetTestMode()

	var kb interfaces.Backend

	if config.VaultStorage {
		kb = vaultManager.MakeVaultManager(nil, config.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(nil, config.PrivateKeyFolder, config.KeyPrefix)
	}

	ctx := context.Background()
	pgpMan = MakePGPManagerWithKRM(nil, kb, MakeKeyRingManager(nil)).(*PGPManager)
	pgpMan.LoadKeys(ctx)

	sm = MakeSecretsManager(nil)

	err = pgpMan.UnlockKey(ctx, testdata.TestKeyFingerprint, testdata.TestKeyPassword)

	if err != nil {
		slog.SetError(true)
		slog.Error(err)
		os.Exit(1)
	}

	code := m.Run()
	etc.ResetDatabase()
	slog.UnsetTestMode()
	etc.Cleanup()
	os.Exit(code)
}
