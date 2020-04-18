// Package wtools implements an additional function to analyze XMLs in Remote API request body
package wtools

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

// CalculateSHA1Checksum returns checksum with SHA1 algo in string data type.
func CalculateSHA1Checksum(pk, sth string) string {
	digester := hmac.New(sha1.New, []byte(pk))
	digester.Write([]byte(sth))
	checksum := hex.EncodeToString(digester.Sum(nil))
	return checksum
}

// CalculateSHA256Checksum returns checksum with SHA256 algo in string data type.
func CalculateSHA256Checksum(pk, sth string) string {
	digester := hmac.New(sha256.New, []byte(pk))
	digester.Write([]byte(sth))
	checksum := hex.EncodeToString(digester.Sum(nil))
	return checksum
}
