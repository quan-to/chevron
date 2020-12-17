package cache

import (
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/pkg/database/memory"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/pkg/uuid"
	"github.com/quan-to/slog"
)

var testUser = models.User{
	ID:          "abcd",
	Username:    "johnhuebr",
	FullName:    "John HUEBR",
	Fingerprint: "DEADBEEFDEADBEEF",
	Password:    "I think you will never guess",
	CreatedAt:   time.Now().Truncate(time.Second),
}

func TestDriver_AddUser(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	monkey.Patch(uuid.EnsureUUID, func(log slog.Instance) string {
		return testUser.ID
	})

	// Passthrough call, test if consistent
	_, err := h.AddUser(testUser)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	user, err := mem.GetUser(testUser.Username)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if diff := pretty.Compare(testUser, user); diff != "" {
		t.Errorf("Expected key list to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_GetUser(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	monkey.Patch(uuid.EnsureUUID, func(log slog.Instance) string {
		return testUser.ID
	})

	_, err := mem.AddUser(testUser)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	// Passthrough call, test if consistent
	user, err := mem.GetUser(testUser.Username)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if diff := pretty.Compare(testUser, user); diff != "" {
		t.Errorf("Expected key list to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_UpdateUser(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	monkey.Patch(uuid.EnsureUUID, func(log slog.Instance) string {
		return testUser.ID
	})

	_, err := mem.AddUser(testUser)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	testUserUpdate := testUser
	testUserUpdate.Password = "WOLOLO"
	testUserUpdate.FullName = "JHUEBAUH#$IUH"

	// Passthrough call, test if consistent
	err = mem.UpdateUser(testUserUpdate)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	user, err := mem.GetUser(testUser.Username)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if diff := pretty.Compare(testUserUpdate, user); diff != "" {
		t.Errorf("Expected key list to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}
