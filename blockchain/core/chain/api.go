package chain

import (
	"errors"
	"fmt"
	"kobla/blockchain/core/types"
)

var ErrInvalidBlockNumber = errors.New("invalid block number")

func (bc *Blockchain) BlockByHash(hash types.Hash) (*types.Block, error) {
	return bc.db.Block(hash)
}

func (bc *Blockchain) BlockByNumber(number int64) (*types.Block, error) {
	if number == -1 {
		return bc.lastBlock().Copy(), nil
	}

	if number < 0 {
		return nil, ErrInvalidBlockNumber
	}

	hash, err := bc.db.BlockHash(number)
	if err != nil {
		return nil, fmt.Errorf("get block hash: %w", err)
	}

	return bc.db.Block(hash)
}

func (bc *Blockchain) SubscribeNewBlocks(subCh chan *types.Block) SubscriptionID {
	subCh <- bc.lastBlock().Copy()
	return bc.blockSubs.subscribe(subCh)
}

func (bc *Blockchain) UnsubscribeBlocks(id SubscriptionID) {
	bc.blockSubs.unsubscribe(id)
}
