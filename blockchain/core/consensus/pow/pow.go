package pow

import (
	"encoding/binary"
	"errors"
	"math"
	"math/big"
	"path2perpetuity/blockchain/core/common"
)

var ErrNonceNotFound = errors.New("nonce not found")

const (
	targetBits = 24
	maxNonce   = math.MaxUint64
)

var target = new(big.Int).Lsh(common.BigOne, common.HashBits-targetBits)

func Run(data []byte) (uint64, common.Hash, error) {
	var nonce uint64

	for nonce < maxNonce {
		hash := common.CalculateHash(appendNonceBytes(data, nonce))

		if hashIsValid(hash) {
			return nonce, hash, nil
		}

		nonce++
	}

	return 0, common.EmptyHash, ErrNonceNotFound
}

func Validate(data []byte, nonce uint64) bool {
	return hashIsValid(common.CalculateHash(appendNonceBytes(data, nonce)))
}

func hashIsValid(hash common.Hash) bool {
	hashBigInt := new(big.Int).SetBytes(hash[:])
	return hashBigInt.Cmp(target) < 0
}

func appendNonceBytes(data []byte, nonce uint64) []byte {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	return append(data, nonceBytes...)
}
