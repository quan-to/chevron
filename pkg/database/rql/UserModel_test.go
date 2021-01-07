package rql

import (
	"strings"
	"testing"

	"github.com/quan-to/chevron/pkg/models/testmodels"

	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func TestRethinkDBDriver_AddUser(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	userToAdd := testmodels.User
	userToAdd.ID = ""

	mock.On(r.Table(userModelTableInit.TableName).Insert(map[string]interface{}{
		"Fingerprint": userToAdd.Fingerprint,
		"Username":    userToAdd.Username,
		"Password":    userToAdd.Password,
		"FullName":    userToAdd.FullName,
		"CreatedAt":   r.MockAnything(),
	})).
		Return(r.WriteResponse{
			Inserted:      1,
			GeneratedKeys: []string{"abcd"},
		}, nil).Times(1)

	mock.On(r.Table(userModelTableInit.TableName).
		GetAllByIndex("Username", userToAdd.Username)).
		Return(nil, nil)

	id, err := h.AddUser(userToAdd)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if id != "abcd" {
		t.Errorf("Expected ID to be %s but got %s", "abcd", id)
	}

	mock.AssertExpectations(t)

	// Test error
	mock = r.NewMock()
	h.conn = mock

	mock.On(r.Table(userModelTableInit.TableName).
		GetAllByIndex("Username", userToAdd.Username)).
		Return([]map[string]interface{}{
			{"id": "abcd"},
		}, nil)

	_, err = h.AddUser(userToAdd)

	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "already exists") {
		t.Fatalf("expected error: %q got %q", "already exists", err)
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_GetUser(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	expectedUser := testmodels.User
	expectedUser.ID = "abcd1234"

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(
		r.Table(userModelTableInit.TableName).
			GetAllByIndex("Username", expectedUser.Username).
			Limit(1).
			CoerceTo("array"),
	).
		Return([]map[string]interface{}{
			{
				"id":          expectedUser.ID,
				"Fingerprint": expectedUser.Fingerprint,
				"Username":    expectedUser.Username,
				"FullName":    expectedUser.FullName,
				"Password":    expectedUser.Password,
				"CreatedAt":   expectedUser.CreatedAt,
			},
		}, nil))

	user, err := h.GetUser(expectedUser.Username)

	if err != nil {
		t.Fatalf("expected no error but got %q", err)
	}

	if diff := pretty.Compare(expectedUser, user); diff != "" {
		t.Errorf("Expected user to be the same. (-got +want)\\n%s", diff)
	}

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userModelTableInit.TableName).
		GetAllByIndex("Username", "huebr").
		Limit(1).
		CoerceTo("array")).
		Return([]map[string]interface{}{}, nil))

	_, err = h.GetUser("huebr")

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "not found") {
		t.Fatalf("Expected error to be %q but got %q", "not found", err.Error())
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_UpdateUser(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testmodels.User)

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userModelTableInit.TableName).
		GetAllByIndex("Username", testmodels.User.Username).
		Update(m)).Return(r.WriteResponse{
		Replaced: 1,
	}, nil))

	m2, _ := convertToRethinkDB(testmodels.User)
	m2["Username"] = testmodels.User.Username + "HUEBR"

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userModelTableInit.TableName).
		GetAllByIndex("Username", m2["Username"]).
		Update(m2)).Return(r.WriteResponse{
		Replaced: 0,
	}, nil))

	err := h.UpdateUser(testmodels.User)

	if err != nil {
		t.Fatalf("Unexpected error %q", err)
	}

	errorUser := testmodels.User
	errorUser.Username += "HUEBR"

	err = h.UpdateUser(errorUser)

	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "not found") {
		t.Fatalf("Expected error to be %q but got %q", "not found", err.Error())
	}

	mock.AssertExpectations(t)
}
