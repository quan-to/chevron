package interfaces

import (
	"context"
)

// SMInterface is a interface for a encrypted secret password manager
type SMInterface interface {
	// PutKeyPassword stores the password for the specified key fingerprint in the key backend encrypted with the master key
	PutKeyPassword(ctx context.Context, fingerPrint, password string)
	// PutEncryptedPassword stores in memory a master key encrypted password for the specified fingerprint
	PutEncryptedPassword(ctx context.Context, fingerPrint, encryptedPassword string)
	// GetPasswords returns a list of master key encrypted passwords stored in memory
	GetPasswords(ctx context.Context) map[string]string
	// UnlockLocalKeys unlocks the local private keys using memory stored master key encrypted passwords
	UnlockLocalKeys(ctx context.Context, gpg PGPInterface)
	// GetMasterKeyFingerPrint returns the fingerprint of the master key
	GetMasterKeyFingerPrint(ctx context.Context) string
}
