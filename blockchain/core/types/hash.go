package types

import (
	"crypto/sha256"
	"encoding/hex"
)

const (
	HashBits  = 256
	HashBytes = 32
)

var EmptyHash = Hash{}

type Hash [HashBytes]byte

func NewHash(data []byte) Hash {
	hash := sha256.Sum256(data)
	return hash
}

func HashFromSlice(hashData []byte) (hash Hash) {
	if len(hashData) != HashBytes {
		panic("invalid hash length")
	}

	copy(hash[:], hashData)
	return
}

func HashFromHex(hexData string) (hash Hash) {
	data, err := hex.DecodeString(hexData)
	if err != nil {
		panic("invalid hex data")
	}

	return HashFromSlice(data)
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}
