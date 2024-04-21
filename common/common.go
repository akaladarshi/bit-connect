package common

import "crypto/sha256"

const (
	Regtest               uint32 = 0xdab5bffa
	LatestProtocolVersion        = 70016
)

// PayloadHash calculates the double hash of the payload.
func PayloadHash(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:]
}
