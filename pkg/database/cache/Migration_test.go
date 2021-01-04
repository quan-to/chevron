package cache

import (
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/quan-to/chevron/pkg/models/testmodels"
	"github.com/quan-to/chevron/pkg/uuid"
	"github.com/quan-to/slog"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/quan-to/chevron/pkg/database/memory"
	"github.com/quan-to/chevron/pkg/models"
)

func TestDriver_InitCursor(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	// Passthrough test. Memory driver doesnt support migrations
	// so it will return error
	gotErr := h.InitCursor()
	expectedErr := h.proxy.InitCursor()

	if gotErr == nil {
		t.Fatal("expected gotError got nil")
	}
	if expectedErr == nil {
		t.Fatal("expected expectedError got nil")
	}

	if !strings.EqualFold(expectedErr.Error(), gotErr.Error()) {
		t.Fatalf("expected error to be %q got %q", expectedErr.Error(), gotErr.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_FinishCursor(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	// Passthrough test. Memory driver doesnt support migrations
	// so it will return error
	gotErr := h.FinishCursor()
	expectedErr := h.proxy.FinishCursor()

	if gotErr == nil {
		t.Fatal("expected gotError got nil")
	}
	if expectedErr == nil {
		t.Fatal("expected expectedError got nil")
	}

	if !strings.EqualFold(expectedErr.Error(), gotErr.Error()) {
		t.Fatalf("expected error to be %q got %q", expectedErr.Error(), gotErr.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_NextGPGKey(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	// Passthrough test. Memory driver doesnt support migrations
	// so it will return error
	gpgKey := models.GPGKey{}
	gotErr := h.NextGPGKey(&gpgKey)
	expectedErr := h.proxy.NextGPGKey(&gpgKey)

	if gotErr == true {
		t.Fatal("expected gotError got nil")
	}
	if expectedErr == true {
		t.Fatal("expected expectedError got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_NextUser(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	// Passthrough test. Memory driver doesnt support migrations
	// so it will return error
	user := models.User{}
	gotErr := h.NextUser(&user)
	expectedErr := h.proxy.NextUser(&user)

	if gotErr == true {
		t.Fatal("expected gotError got nil")
	}
	if expectedErr == true {
		t.Fatal("expected expectedError got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_NumGPGKeys(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	// Passthrough test. Memory driver doesnt support migrations
	// so it will return error
	_, gotErr := h.NumGPGKeys()
	_, expectedErr := h.proxy.NumGPGKeys()

	if gotErr == nil {
		t.Fatal("expected gotError got nil")
	}
	if expectedErr == nil {
		t.Fatal("expected expectedError got nil")
	}

	if !strings.EqualFold(expectedErr.Error(), gotErr.Error()) {
		t.Fatalf("expected error to be %q got %q", expectedErr.Error(), gotErr.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_AddGPGKeys(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	testKeyToAdd := testmodels.GpgKey
	testKeyToAdd.ID = "0000"

	monkey.Patch(uuid.EnsureUUID, func(log slog.Instance) string {
		return testKeyToAdd.ID
	})

	data, err := h.cache.Marshal(&testKeyToAdd)

	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	mock.ExpectSet(gpgKeyByIDPrefix+testKeyToAdd.ID, data, gpgKeyExpiration).
		SetVal("")
	mock.ExpectSet(gpgKeyByFingerprintPrefix+testKeyToAdd.GetShortFingerPrint(), data, gpgKeyExpiration).
		SetVal("")

	id, added, err := h.AddGPGKeys([]models.GPGKey{testmodels.GpgKey})
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if len(added) == 0 || !added[0] {
		t.Fatal("expected key to be added")
	}

	// Check if key has been added
	key, err := mem.FetchGPGKeyByFingerprint(testmodels.GpgKey.FullFingerprint)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if key.ID != id[0] {
		t.Fatalf("expected added key ID to be %s got %s", id, key.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}
