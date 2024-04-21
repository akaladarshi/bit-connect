package common

import "crypto/sha256"

// PayloadHash calculates the double hash of the payload.
func PayloadHash(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:]
}
