package vaultManager

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"net/http"
)

var slog = SLog.Scope("Vault")

const qrsVaultPrefix = "qrs"

type VaultManager struct {
	client *api.Client
	prefix string
	log    *SLog.Instance
}

func MakeVaultManager(prefix string) *VaultManager {
	slog.Info("Initialized Vault Backend at %s with prefix %s", remote_signer.VaultAddress, prefix)
	var httpClient *http.Client
	if remote_signer.VaultSkipVerify {
		slog.Warn("WARNING: Vault Skip Verify is enable. We will not check for SSL Certs in Vault!")
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
		slog.Error(err)
		return nil
	}
	if !remote_signer.VaultUseUserpass {
		slog.Info("Token Mode enabled.")
		client.SetToken(remote_signer.VaultRootToken)
	} else {
		// to pass the password
		options := map[string]interface{}{
			"password": remote_signer.VaultPassword,
		}
		slog.Info("Userpass mode enabled. Logging with %s", remote_signer.VaultUsername)
		// PUT call to get a token
		secret, err := client.Logical().Write(fmt.Sprintf("auth/userpass/login/%s", remote_signer.VaultUsername), options)

		if err != nil {
			slog.Error(err)
			return nil
		}

		slog.Info("Logged in successfully.")
		client.SetToken(secret.Auth.ClientToken)
	}

	return &VaultManager{
		client: client,
		prefix: prefix,
		log:    SLog.Scope(fmt.Sprintf("Vault (%s)", prefix)),
	}
}

func vaultPath(key string) string {
	return fmt.Sprintf("%s/data/%s", remote_signer.VaultNamespace, getVaultFullPrefix(key))
}

func (vm *VaultManager) putSecret(key string, data map[string]string) error {
	_, err := vm.client.Logical().Write(vaultPath(key), map[string]interface{}{
		"data": data,
	})

	return err
}

func (vm *VaultManager) getSecret(key string) (string, string, error) {
	s, err := vm.client.Logical().Read(vaultPath(key))
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

func getVaultFullPrefix(key string) string {
	return remote_signer.VaultPathPrefix + qrsVaultPrefix + "_" + key
}

func (vm *VaultManager) Save(key, data string) error {
	vm.log.Debug("Saving %s", key)
	return vm.putSecret(vm.prefix+key, map[string]string{
		"data": data,
	})
}

func (vm *VaultManager) SaveWithMetadata(key, data, metadata string) error {
	vm.log.Debug("Saving %s", key)
	return vm.putSecret(vm.prefix+key, map[string]string{
		"data":     data,
		"metadata": metadata,
	})
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
	vm.log.Debug("Listing keys")
	s, err := vm.client.Logical().List(fmt.Sprintf("%s/metadata", remote_signer.VaultNamespace))
	if err != nil {
		return nil, err
	}

	if s == nil {
		return make([]string, 0), nil
	}

	keys := make([]string, 0)
	data := s.Data["keys"].([]interface{})
	basePrefix := getVaultFullPrefix(vm.prefix)

	for _, v := range data {
		v2 := v.(string)
		if len(v2) > len(basePrefix) && v2[:len(basePrefix)] == basePrefix {
			keys = append(keys, v2[len(basePrefix):])
		}
	}

	return keys, nil
}

func (vm *VaultManager) Name() string {
	return "Vault Backend"
}

func (vm *VaultManager) Path() string {
	return getVaultFullPrefix(vm.prefix + "*")
}
