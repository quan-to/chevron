package cache

import (
	"testing"

	"github.com/quan-to/chevron/pkg/models/testmodels"

	"bou.ke/monkey"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/pkg/database/memory"
	"github.com/quan-to/chevron/pkg/uuid"
	"github.com/quan-to/slog"
)

func TestDriver_AddUser(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	monkey.Patch(uuid.EnsureUUID, func(log slog.Instance) string {
		return testmodels.User.ID
	})

	// Passthrough call, test if consistent
	_, err := h.AddUser(testmodels.User)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	user, err := mem.GetUser(testmodels.User.Username)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if diff := pretty.Compare(testmodels.User, user); diff != "" {
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
		return testmodels.User.ID
	})

	_, err := mem.AddUser(testmodels.User)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	// Passthrough call, test if consistent
	user, err := mem.GetUser(testmodels.User.Username)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if diff := pretty.Compare(testmodels.User, user); diff != "" {
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
		return testmodels.User.ID
	})

	_, err := mem.AddUser(testmodels.User)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	TestUserUpdate := testmodels.User
	TestUserUpdate.Password = "WOLOLO"
	TestUserUpdate.FullName = "JHUEBAUH#$IUH"

	// Passthrough call, test if consistent
	err = mem.UpdateUser(TestUserUpdate)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	user, err := mem.GetUser(testmodels.User.Username)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if diff := pretty.Compare(TestUserUpdate, user); diff != "" {
		t.Errorf("Expected key list to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}
