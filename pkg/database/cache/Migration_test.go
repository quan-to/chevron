package cache

import (
	"strings"
	"testing"

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
