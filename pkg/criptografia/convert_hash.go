package criptografia

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncryptInSha256(txt string) string {

	hashBytes := sha256.Sum256([]byte(txt))

	return hex.EncodeToString(hashBytes[:])
}
