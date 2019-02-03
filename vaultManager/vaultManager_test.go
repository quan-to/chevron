package vaultManager

import (
	"os"
	"testing"
)

var vm *VaultManager

func TestMain(m *testing.M) {
	vm = MakeVaultManager("test_")

	code := m.Run()

	os.Exit(code)
}

func TestVaultManager_Name(t *testing.T) {
	if vm.Name() != "Vault Backend" {
		t.Errorf("Expected Vault Backend got %s", vm.Name())
		t.FailNow()
	}
}

func TestVaultManager_Save(t *testing.T) {
	err := vm.Save("testkey", "testdata")
	if err != nil {
		t.Errorf("Error saving key: %s", err)
		t.FailNow()
	}

	data, _, err := vm.Read("testkey")
	if err != nil {
		t.Errorf("Error loading key: %s", err)
		t.FailNow()
	}

	if data != "testdata" {
		t.Errorf("Expected %s got %s", "testdata", data)
	}
}

func TestVaultManager_SaveWithMetadata(t *testing.T) {
	err := vm.SaveWithMetadata("testkey_meta", "testdata", "testmetadata")
	if err != nil {
		t.Errorf("Error saving key: %s", err)
		t.FailNow()
	}

	data, metadata, err := vm.Read("testkey_meta")
	if err != nil {
		t.Errorf("Error loading key: %s", err)
		t.FailNow()
	}

	if data != "testdata" {
		t.Errorf("Expected %s got %s", "testdata", data)
	}

	if metadata != "testmetadata" {
		t.Errorf("Expected %s got %s", "testmetadata", metadata)
	}
}
