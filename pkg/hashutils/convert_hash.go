package hashutils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSHA256(txt string) string {

	hashBytes := sha256.Sum256([]byte(txt))

	return hex.EncodeToString(hashBytes[:])
}
