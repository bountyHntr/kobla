package common

import (
	"crypto/sha256"
	"math/big"
)

const (
	HashBits  = 256
	HashBytes = 32
)

var (
	BigZero = big.NewInt(0)
	BigOne  = big.NewInt(1)

	EmptyHash = Hash{}
)

type Hash [HashBytes]byte

func CalculateHash(data []byte) Hash {
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
