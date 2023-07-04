package pow

import (
	"encoding/binary"
	"errors"
	"math"
	"math/big"
	"path2perpetuity/blockchain/core/utils"
)

var ErrNonceNotFound = errors.New("nonce not found")

const (
	targetBits = 24
	maxNonce   = math.MaxUint64
)

var target = new(big.Int).Lsh(utils.BigOne, utils.HashBits-targetBits)

func Run(data []byte) (uint64, []byte, error) {
	var nonce uint64

	for nonce < maxNonce {
		hash := utils.Hash(appendNonceBytes(data, nonce))

		if hashIsValid(hash) {
			return nonce, hash[:], nil
		}

		nonce++
	}

	return 0, nil, ErrNonceNotFound
}

func Validate(data []byte, nonce uint64) bool {
	return hashIsValid(utils.Hash(appendNonceBytes(data, nonce)))
}

func hashIsValid(hash []byte) bool {
	hashBigInt := new(big.Int).SetBytes(hash[:])
	return hashBigInt.Cmp(target) < 0
}

func appendNonceBytes(data []byte, nonce uint64) []byte {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	return append(data, nonceBytes...)
}
