package poa

import (
	"errors"
	"kobla/blockchain/core/types"
)

var ErrBlockAlreadyMined = errors.New("block already mined")

type Validator struct {
	Url     string
	Address types.Address
}

type ProofOfAuthority struct {
	validators map[Validator]struct{}
}

func New(validators []Validator) types.ConsesusProtocol {
	poa := ProofOfAuthority{
		validators: make(map[Validator]struct{}, len(validators)),
	}

	for _, validator := range validators {
		poa.validators[validator] = struct{}{}
	}

	return &poa
}

// updates the state of the block
func (poa *ProofOfAuthority) Run(block *types.Block, _ any) error {
	block.SetHash()
	return nil
}

func (poa *ProofOfAuthority) Validate(block *types.Block, meta any) bool {
	if block.Coinbase == types.ZeroAddress && block.Number == 0 {
		return true
	}

	validator := Validator{
		Url:     meta.(string),
		Address: block.Coinbase,
	}

	_, ok := poa.validators[validator]
	return ok
}
