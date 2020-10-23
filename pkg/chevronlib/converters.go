package chevronlib

import "github.com/quan-to/chevron/internal/tools"

// GPG2Quanto converts a GPG Signature to Quanto Format
func GPG2Quanto(signature, fingerprint, hash string) string {
	return tools.GPG2Quanto(signature, fingerprint, hash)
}

// GPG2Quanto converts a Quanto Signature to GPG
func Quanto2GPG(signature string) string {
	return tools.Quanto2GPG(signature)
}

// GetFingerprintFromKey returns the main fingerprint from key or error if key is invalid
func GetFingerprintFromKey(armored string) (string, error) {
	return tools.GetFingerPrintFromKey(armored)
}
