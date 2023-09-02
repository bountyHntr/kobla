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

func (ProofOfWork) NodesAreFixed() bool {
	return false
}

// updates the state of the block
func (ProofOfWork) Run(block *types.Block, _ any) error {
	if block.Nonce != 0 {
		return ErrBlockAlreadyMined
	}

	for block.Nonce < maxNonce {
		block.SetHash()

		if hashIsValid(block.Hash) {
			return nil
		}

		block.Nonce++
	}

	return ErrNonceNotFound
}

func (ProofOfWork) Validate(block *types.Block, _ any) bool {
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
