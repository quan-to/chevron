package etc

import "context"

type SMInterface interface {
	PutKeyPassword(ctx context.Context, fingerPrint, password string)
	PutEncryptedPassword(ctx context.Context, fingerPrint, encryptedPassword string)
	GetPasswords(ctx context.Context) map[string]string
	UnlockLocalKeys(ctx context.Context, gpg PGPInterface)
	GetMasterKeyFingerPrint(ctx context.Context) string
}
