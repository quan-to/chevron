// +build !js,!wasm

package vaultManager

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/quan-to/chevron"
	"github.com/quan-to/slog"
	"net/http"
)

const VaultData = "data"
const VaultMetadata = "metadata"

type VaultManager struct {
	client *api.Client
	prefix string
	log    slog.Instance
}

// MakeVaultManager creates an instance of VaultManager
func MakeVaultManager(log slog.Instance, prefix string) *VaultManager {
	if log == nil {
		log = slog.Scope("Vault")
	} else {
		log = log.SubScope("Vault")
	}

	log.Info("Initialized Vault Backend at %s with prefix %s", remote_signer.VaultAddress, prefix)
	var httpClient *http.Client
	if remote_signer.VaultSkipVerify {
		log.Warn("WARNING: Vault Skip Verify is enable. We will not check for SSL Certs in Vault!")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: tr}
	}

	client, err := api.NewClient(&api.Config{
		Address:    remote_signer.VaultAddress,
		HttpClient: httpClient,
	})

	if err != nil {
		log.Error(err)
		return nil
	}
	if !remote_signer.VaultUseUserpass {
		log.Info("Token Mode enabled.")
		client.SetToken(remote_signer.VaultRootToken)
	} else {
		// to pass the password
		options := map[string]interface{}{
			"password": remote_signer.VaultPassword,
		}
		log.Info("Userpass mode enabled. Logging with %s", remote_signer.VaultUsername)
		// PUT call to get a token
		secret, err := client.Logical().Write(fmt.Sprintf("auth/userpass/login/%s", remote_signer.VaultUsername), options)

		if err != nil {
			log.Error(err)
			return nil
		}

		log.Info("Logged in successfully.")
		client.SetToken(secret.Auth.ClientToken)
	}

	return &VaultManager{
		client: client,
		prefix: prefix,
		log:    slog.Scope(fmt.Sprintf("Vault (%s)", prefix)),
	}
}

func baseVaultPath(dataType string) string {
	if remote_signer.VaultSkipDataType {
		return fmt.Sprintf("%s/%s", remote_signer.VaultBackend, remote_signer.VaultNamespace)
	}
	return fmt.Sprintf("%s/%s/%s", remote_signer.VaultBackend, dataType, remote_signer.VaultNamespace)
}

func (vm *VaultManager) vaultPath(dataType, key string) string {
	return fmt.Sprintf("%s/%s", baseVaultPath(dataType), key)
}

func (vm *VaultManager) putSecret(key string, data map[string]string) error {
	vm.log.DebugAwait("Saving %s to %s", key, vm.vaultPath(VaultData, key))
	_, err := vm.client.Logical().Write(vm.vaultPath(VaultData, key), map[string]interface{}{
		"data": data,
	})

	if err != nil {
		vm.log.ErrorDone("Error saving %s to %s: %s", err)
	}

	return err
}

func (vm *VaultManager) deleteSecret(key string) error {
	vm.log.DebugAwait("Deleting %s from %s", key, vm.vaultPath(VaultData, key))
	_, err := vm.client.Logical().Read(vm.vaultPath(VaultData, key))
	if err != nil {
		vm.log.ErrorDone("Error reading to %s: %s, file not exist to delete", vm.vaultPath(VaultData, key), err)
		return err
	}

	_, err = vm.client.Logical().Delete(vm.vaultPath(VaultData, key))
	if err != nil {
		vm.log.ErrorDone("Error deleting from %s: %s", vm.vaultPath(VaultData, key), err)
	}

	return err
}

func (vm *VaultManager) getSecret(key string) (string, string, error) {
	//vm.log.Debug("getSecret(%s)", key)
	s, err := vm.client.Logical().Read(vm.vaultPath(VaultData, key))
	if err != nil {
		return "", "", err
	}

	if s == nil {
		return "", "", fmt.Errorf("not found")
	}

	data := s.Data["data"].(map[string]interface{})

	if data["data"] == nil {
		return "", "", fmt.Errorf("corrupted data")
	}

	d := data["data"].(string)
	m := ""
	if data["metadata"] != nil {
		m = data["metadata"].(string)
	}

	return d, m, nil
}

func (vm *VaultManager) HealthStatus() (*api.HealthResponse, error) {
	return vm.client.Sys().Health()
}

func (vm *VaultManager) Save(key, data string) error {
	vm.log.DebugAwait("Saving %s", key)
	return vm.putSecret(vm.prefix+key, map[string]string{
		"data": data,
	})
}

func (vm *VaultManager) SaveWithMetadata(key, data, metadata string) error {
	return vm.putSecret(vm.prefix+key, map[string]string{
		"data":     data,
		"metadata": metadata,
	})
}

// Delete deletes a key from the vault
func (vm *VaultManager) Delete(key string) error {
	vm.log.DebugAwait("Deleting %s", key)
	return vm.deleteSecret(vm.prefix+key)
}

func (vm *VaultManager) Read(key string) (data string, metadata string, err error) {
	vm.log.Debug("Reading %s", key)
	d, m, err := vm.getSecret(vm.prefix + key)
	if err != nil {
		return "", "", err
	}

	return d, m, nil
}

func (vm *VaultManager) List() ([]string, error) {
	vm.log.Debug("Listing keys on %s", baseVaultPath(VaultMetadata))
	s, err := vm.client.Logical().List(baseVaultPath(VaultMetadata))
	if err != nil {
		return nil, err
	}

	if s == nil {
		return make([]string, 0), nil
	}

	keys := make([]string, 0)
	data := s.Data["keys"].([]interface{})

	for _, v := range data {
		v2 := v.(string)
		if len(v2) > len(vm.prefix) && v2[:len(vm.prefix)] == vm.prefix {
			keys = append(keys, v2[len(vm.prefix):])
		}
	}

	vm.log.Debug("Found %d keys", len(keys))

	return keys, nil
}

func (vm *VaultManager) Name() string {
	return "Vault Backend"
}

func (vm *VaultManager) Path() string {
	return vm.vaultPath(VaultData, vm.prefix+"*")
}
