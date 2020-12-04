package keymagic

import (
	"context"
	"fmt"
	"github.com/bouk/monkey"
	"github.com/google/uuid"
	config "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keybackend"
	"github.com/quan-to/chevron/internal/vaultManager"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/chevron/pkg/database/memory"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/test"
	"os"
	"testing"

	"github.com/quan-to/slog"
)

var testData = []byte(test.TestSignatureData)

var pgpMan *pgpManager
var sm *secretsManager

func TestMain(m *testing.M) {
	slog.UnsetTestMode()
	var err error

	QuantoError.EnableStackTrace()
	slog.SetTestMode()
	u, _ := uuid.NewRandom()
	config.DatabaseName = "qrs_test_" + u.String()
	config.PrivateKeyFolder = "../../test/data/"
	config.KeyPrefix = "testkey_"
	config.KeysBase64Encoded = false

	config.MasterGPGKeyBase64Encoded = false
	config.MasterGPGKeyPath = "../../test/data/testkey_privateTestKey.gpg"
	config.MasterGPGKeyPasswordPath = "../../test/data/testprivatekeyPassword.txt"

	config.HttpPort = 40000
	config.SKSServer = fmt.Sprintf("http://localhost:%d/sks/", config.HttpPort)
	config.EnableRethinkSKS = true
	config.PushVariables()

	var kb interfaces.StorageBackend

	if config.VaultStorage {
		kb = vaultManager.MakeVaultManager(nil, config.KeyPrefix)
	} else {
		kb = keybackend.MakeSaveToDiskBackend(nil, config.PrivateKeyFolder, config.KeyPrefix)
	}

	ctx := context.Background()
	mem := memory.MakeMemoryDBDriver(nil)
	ctx = context.WithValue(ctx, "dbHandler", mem)
	pgpMan = MakePGPManager(nil, kb, MakeKeyRingManager(nil, mem)).(*pgpManager)
	pgpMan.LoadKeys(ctx)

	sm = MakeSecretsManager(nil, mem).(*secretsManager)

	err = pgpMan.UnlockKey(ctx, test.TestKeyFingerprint, test.TestKeyPassword)

	if err != nil {
		slog.SetError(true)
		slog.Error(err)
		os.Exit(1)
	}

	code := m.Run()
	slog.UnsetTestMode()
	os.Exit(code)
}

func assertPanic(t *testing.T, f func(), message string) {
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
		}
	}()
	f()
}
