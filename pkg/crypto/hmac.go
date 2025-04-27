// pkg/crypto/hmac.go
package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateHMAC создает HMAC-SHA256 подпись для данных
func GenerateHMAC(data []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	return h.Sum(nil)
}

// VerifyHMAC проверяет соответствие HMAC подписи
func VerifyHMAC(data []byte, receivedMAC []byte, secret []byte) bool {
	expectedMAC := GenerateHMAC(data, secret)
	return hmac.Equal(receivedMAC, expectedMAC)
}

// GenerateHMACHex возвращает HMAC в виде hex-строки
func GenerateHMACHex(data []byte, secret []byte) string {
	return hex.EncodeToString(GenerateHMAC(data, secret))
}
