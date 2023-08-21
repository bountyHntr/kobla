package chain

import (
	"errors"
	"fmt"
	"kobla/blockchain/core/types"

	log "github.com/sirupsen/logrus"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	ErrInvalidBlockNumber = errors.New("invalid block number")
	ErrUnknownTxHash      = errors.New("unknown tx hash")
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func (bc *Blockchain) MineBlock(txs []*types.Transaction, coibase types.Address) error {
	return bc.mineBlock(txs, coibase)
}

func (bc *Blockchain) SendTx(tx *types.Transaction) {
	if ok := bc.mempool.add(tx); !ok {
		return
	}
	bc.comm.broadcast(tx)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

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

func (bc *Blockchain) TxByHash(txHash types.Hash) (*types.Transaction, error) {
	blockHash, err := bc.db.TxToBlock(txHash)
	if err != nil {
		return nil, fmt.Errorf("get block by tx hash: %w", err)
	}

	block, err := bc.BlockByHash(blockHash)
	if err != nil {
		return nil, err
	}

	for _, tx := range block.Transactions {
		if tx.Hash == txHash {
			return tx, nil
		}
	}

	log.Panicf("inconsistent state")
	return nil, nil
}

func (bc *Blockchain) Balance(address types.Address) (uint64, error) {
	return bc.db.Balance(address)
}

func (bc *Blockchain) TopMempoolTxs(n int) []*types.Transaction {
	return bc.mempool.top(n)
}

func (bc *Blockchain) TxByHashFromMempool(txHash types.Hash) (*types.Transaction, error) {
	tx := bc.mempool.get(txHash)
	if tx == nil {
		return nil, ErrUnknownTxHash
	}

	return tx, nil
}

func (bc *Blockchain) MempoolSize() int {
	return bc.mempool.size()
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (bc *Blockchain) SubscribeNewBlocks(subCh chan *types.Block) SubscriptionID {
	subCh <- bc.lastBlock().Copy()
	return bc.blockSubs.subscribe(subCh)
}

func (bc *Blockchain) UnsubscribeBlocks(id SubscriptionID) {
	bc.blockSubs.unsubscribe(id)
}

func (bc *Blockchain) SubscribeMempoolUpdates(subCh chan *types.Void) SubscriptionID {
	return bc.mempool.subscribeUpdates(subCh)
}

func (bc *Blockchain) UnsubscribeMempoolUpdates(id SubscriptionID) {
	bc.mempool.unsubscribeUpdates(id)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
