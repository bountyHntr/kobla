//go:build poa

package poa

import (
	"errors"
	"kobla/blockchain/core/types"

	log "github.com/sirupsen/logrus"
)

var ErrBlockAlreadyMined = errors.New("block already mined")

type ProofOfAuthority struct {
	validators      map[types.Address]struct{}
	coinbaseAccount types.Account
}

func New(validators []string, coinbasePrivKey string) (types.ConsesusProtocol, error) {
	poa := ProofOfAuthority{
		validators: make(map[types.Address]struct{}, len(validators)),
	}

	for _, v := range validators {
		poa.validators[types.AddressFromString(v)] = struct{}{}
	}

	coinbaseAcc, err := types.AccountFromPrivKey(coinbasePrivKey)
	if err != nil {
		return nil, err
	}
	poa.coinbaseAccount = coinbaseAcc

	return &poa, nil
}

// updates the state of the block
func (poa *ProofOfAuthority) Run(block *types.Block) error {
	block.SetHash()
	return block.SetSignature(poa.coinbaseAccount)
}

func (poa *ProofOfAuthority) Validate(block *types.Block) bool {
	if block.Coinbase == types.ZeroAddress && block.Number == 0 {
		return true
	}

	if _, ok := poa.validators[block.Coinbase]; !ok {
		return false
	}

	ok, err := block.Coinbase.Verify(block.Hash, block.Signature)
	if err != nil || !ok {
		if err != nil {
			log.WithError(err).Warn("can't verify signature")
		}
		return false
	}

	return true
}
