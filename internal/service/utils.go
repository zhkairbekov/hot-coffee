// internal/service/utils.go
package service

import (
	"crypto/rand"
	"encoding/hex"
)

func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
