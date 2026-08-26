package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

func HashPassword(p string) string {
	h := sha256.Sum256([]byte("vocational-admission:" + p))
	return hex.EncodeToString(h[:])
}
func CheckPassword(hash, p string) bool {
	got := HashPassword(p)
	return subtle.ConstantTimeCompare([]byte(hash), []byte(got)) == 1
}
