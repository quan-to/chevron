package vault

import (
    " github.com/hashicorp/vault/api"
)

type VaultManager struct {

}

func MakeVaultManager() *VaultManager {
    client, err := api.NewClient(&api.Config{
        Address: vaultAddress,
    })
}