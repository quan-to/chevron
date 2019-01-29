package vaultManager

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
)

var slog = SLog.Scope("Vault")

const qrsVaultPrefix = "qrs"

type VaultManager struct {
	//client *api.Client
	prefix string
	log    *SLog.Instance
}

func MakeVaultManager(prefix string) *VaultManager {
	slog.Info("Initialized Vault Backend at %s with prefix %s", remote_signer.VaultAddress, prefix)
	//client, err := api.NewClient(&api.Config{
	//	Address: remote_signer.VaultAddress,
	//})
	//
	//if err != nil {
	//	slog.Error(err)
	//	return nil
	//}
	//
	//client.SetToken(remote_signer.VaultRootToken)
	//
	//if err != nil {
	//	panic(err)
	//}

	return &VaultManager{
		//client: client,
		prefix: prefix,
		log:    SLog.Scope(fmt.Sprintf("Vault (%s)", prefix)),
	}
}

func vaultPath(key string) string {
	return fmt.Sprintf("secret/data/%s", getVaultFullPrefix(key))
}

func (vm *VaultManager) putSecret(key string, data map[string]string) error {
	//_, err := vm.client.Logical().Write(vaultPath(key), map[string]interface{}{
	//	"data": data,
	//})
	//
	//return err
	return nil
}

func (vm *VaultManager) getSecret(key string) (string, string, error) {
	//s, err := vm.client.Logical().Read(vaultPath(key))
	//if err != nil {
	//	return "", "", err
	//}
	//
	//data := s.Data["data"].(map[string]interface{})
	//
	//if data["data"] == nil {
	//	return "", "", fmt.Errorf("corrupted data")
	//}
	//
	//d := data["data"].(string)
	//m := ""
	//if data["metadata"] != nil {
	//	m = data["metadata"].(string)
	//}
	//
	//return d, m, nil
	return "", "", nil
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
	//vm.log.Debug("Listing keys")
	//s, err := vm.client.Logical().List("secret/metadata")
	//if err != nil {
	//	return nil, err
	//}
	//
	//if s == nil {
	//	return make([]string, 0), nil
	//}
	//
	//keys := make([]string, 0)
	//data := s.Data["keys"].([]interface{})
	//basePrefix := getVaultFullPrefix(vm.prefix)
	//
	//for _, v := range data {
	//	v2 := v.(string)
	//	if len(v2) > len(basePrefix) && v2[:len(basePrefix)] == basePrefix {
	//		keys = append(keys, v2[len(basePrefix):])
	//	}
	//}
	//
	//return keys, nil
	return nil, nil
}

func (vm *VaultManager) Name() string {
	return "Vault Backend"
}

func (vm *VaultManager) Path() string {
	return getVaultFullPrefix(vm.prefix + "*")
}
