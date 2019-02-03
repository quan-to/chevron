package vaultManager

import (
	"github.com/quan-to/remote-signer"
	"os"
	"testing"
)

var vm *VaultManager

func TestMain(m *testing.M) {
	vm = MakeVaultManager("test_")

	code := m.Run()

	os.Exit(code)
}

func TestVaultManager_Make(t *testing.T) {
	remote_signer.PushVariables()
	// Test Vault SkipVerify
	remote_signer.VaultSkipVerify = true
	tmpVm := MakeVaultManager("test_")
	if tmpVm == nil {
		t.Errorf("Expected to get a vaultManager instance, got nil")
	}

	// Test With Root Token
	remote_signer.VaultUseUserpass = false
	t.Logf("Root Token: %s", remote_signer.VaultRootToken)
	tmpVm = MakeVaultManager("test_")
	if tmpVm == nil {
		t.Errorf("Expected to get a vaultManager instance, got nil")
	}

	_, err := tmpVm.List()

	if err != nil {
		t.Errorf("Got error listing: %s", err)
	}

	remote_signer.PopVariables()
}

func TestVaultManager_List(t *testing.T) {
	_ = vm.Save("__list__", "")
	entries, err := vm.List()
	if err != nil {
		t.Errorf("Error listing entries: %s", err)
		t.FailNow()
	}
	found := false
	for _, v := range entries {
		if v == "__list__" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find __list__ in entries.")
	}
	remote_signer.PushVariables()
	// Test error cases
	remote_signer.VaultNamespace = "huebr__123123123123"
	_, err = vm.List()
	if err == nil {
		t.Errorf("Expected error with unknown namespace")
	}
	remote_signer.PopVariables()
}

func TestVaultManager_Read(t *testing.T) {
	_ = vm.Save("testkey", "testdata")

	// Test Read

	data, _, err := vm.Read("testkey")
	if err != nil {
		t.Errorf("Error loading key: %s", err)
		t.FailNow()
	}

	if data != "testdata" {
		t.Errorf("Expected %s got %s", "testdata", data)
	}

	// Test error

	_, _, err = vm.Read("huebr123123012731923")
	if err == nil {
		t.Errorf("Expected error for unknown key")
	}
}

func TestVaultManager_Path(t *testing.T) {
	if vm.Path() != getVaultFullPrefix(vm.prefix+"*") {
		t.Errorf("Expected %s got %s", getVaultFullPrefix(vm.prefix+"*"), vm.Path())
	}
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
