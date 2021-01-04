package rql

import (
	"strings"
	"testing"

	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/pkg/models/testmodels"
)

func TestRethinkDBDriver_AddGPGKeys(t *testing.T) {
	h, mock, toUpdate := prepareAdd(t)

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
	_, _, err = h.AddGPGKeys([]models.GPGKey{{FullFingerprint: "ERR"}})
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error to be %q but got %q", "test error", err.Error())
	}

	mock.AssertExpectations(t)
}
