package rql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/pkg/models/testmodels"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func TestRethinkDBDriver_AddGPGKeys(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	toUpdate := testmodels.GpgKey
	toUpdate.FullFingerprint = "HUEBR"
	toUpdate.ID = "A123"

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.
		Table(gpgKeyTableInit.TableName).
		GetAllByIndex("FullFingerprint", testmodels.GpgKey.FullFingerprint)).
		Return(nil, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.
		Table(gpgKeyTableInit.TableName).
		GetAllByIndex("FullFingerprint", toUpdate.FullFingerprint)).
		Return([]map[string]interface{}{
			{"id": toUpdate.ID},
		}, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.
		Table(gpgKeyTableInit.TableName).
		GetAllByIndex("FullFingerprint", "ERR")).
		Return(nil, fmt.Errorf("test error")))

	m, _ := convertToRethinkDB(testmodels.GpgKey)
	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).Insert(m)).Return(r.WriteResponse{
		GeneratedKeys: []string{testmodels.GpgKey.ID}, Inserted: 1,
	}, nil))

	m2, _ := convertToRethinkDB(toUpdate)
	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).Get(toUpdate.ID).Update(m2)).Return(r.WriteResponse{
		Updated: 1,
	}, nil))

	// Test Create
	id, added, err := h.AddGPGKeys([]models.GPGKey{testmodels.GpgKey})

	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if len(added) == 0 || !added[0] {
		t.Fatalf("expected item to be added")
	}

	if id[0] != testmodels.GpgKey.ID {
		t.Fatalf("expected id to be %q but got %q", testmodels.GpgKey.ID, id)
	}

	// Test Update

	id, added, err = h.AddGPGKeys([]models.GPGKey{toUpdate})
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if len(added) == 0 || added[0] {
		t.Fatalf("expected item to be updated not added")
	}

	if id[0] != toUpdate.ID {
		t.Fatalf("expected id to be %q got %q", toUpdate.ID, id)
	}

	// Test Error
	_, _, err = h.AddGPGKey(models.GPGKey{FullFingerprint: "ERR"})
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error to be %q but got %q", "test error", err.Error())
	}

	mock.AssertExpectations(t)
}
