package pow

import (
	"errors"
	"fmt"
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
	targetBits = 24
	maxNonce   = math.MaxUint64
)

var target = new(big.Int).Lsh(common.BigOne, types.HashBits-targetBits)

type ProofOfWork struct{}

func New() *ProofOfWork {
	return &ProofOfWork{}
}

var _ types.ConsesusProtocol = ProofOfWork{}

// updates the state of the block
func (ProofOfWork) Run(block *types.Block) error {
	if block.Nonce != 0 {
		return ErrBlockAlreadyMined
	}

	for block.Nonce < maxNonce {
		data, err := block.Serialize()
		if err != nil {
			return fmt.Errorf("serialize block: %w", err)
		}

		hash := types.NewHash(data)

		if hashIsValid(hash) {
			block.Hash = hash
			return nil
		}

		block.Nonce++
	}

	return ErrNonceNotFound
}

func (ProofOfWork) Validate(block *types.Block) bool {
	data, err := block.Serialize()
	if err != nil {
		return false
	}

	return hashIsValid(types.NewHash(data))
}

func hashIsValid(hash types.Hash) bool {
	hashBigInt := new(big.Int).SetBytes(hash.Bytes())
	return hashBigInt.Cmp(target) < 0
}
