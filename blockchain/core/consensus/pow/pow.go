//go:build pow

package pow

import (
	"errors"
	"kobla/blockchain/core/common"
	"kobla/blockchain/core/types"
	"math"
	"math/big"
)

var (
	ErrNonceNotFound     = errors.New("nonce not found")
	ErrBlockAlreadyMined = errors.New("block already mined")
)

const (
	targetBits = 8
	maxNonce   = math.MaxUint64
)

var target = new(big.Int).Lsh(common.BigOne, types.HashBits-targetBits)

type ProofOfWork struct{}

func New() types.ConsesusProtocol {
	return &ProofOfWork{}
}

// updates the state of the block
func (ProofOfWork) Run(block *types.Block) error {
	if block.Nonce != 0 {
		return ErrBlockAlreadyMined
	}

	for block.Nonce < maxNonce {
		hash := block.CalcHash()

		if hashIsValid(hash) {
			block.Hash = hash
			return nil
		}

		block.Nonce++
	}

	return ErrNonceNotFound
}

func (ProofOfWork) Validate(block *types.Block) bool {
	return hashIsValid(block.CalcHash())
}

func hashIsValid(hash types.Hash) bool {
	hashBigInt := new(big.Int).SetBytes(hash.Bytes())
	return hashBigInt.Cmp(target) < 0
}
