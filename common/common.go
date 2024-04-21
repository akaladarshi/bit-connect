package common

import "crypto/sha256"

const (
	Regtest               uint32 = 0xdab5bffa
	LatestProtocolVersion        = 70016
)

// PayloadHash calculates the double hash of the payload.
// we only require the first 4 bytes of the hash for the checksum
// ignoring the rest of the hash
func PayloadHash(payload []byte) [4]byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return [4]byte(secondHash[:])
}
