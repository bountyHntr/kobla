package utils

import (
	"crypto/sha256"
	"math/big"
)

const (
	HashBits = 256
)

var (
	BigZero = big.NewInt(0)
	BigOne  = big.NewInt(1)
)

func Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
