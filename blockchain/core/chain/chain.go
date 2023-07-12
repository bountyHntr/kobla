package chain

import (
	"errors"
	"fmt"
	"kobla/blockchain/core/db"
	"kobla/blockchain/core/types"
	"sync"

	log "github.com/sirupsen/logrus"
)

var ErrInvalidParentBlock = errors.New("invalid parent block")

type Config struct {
	DBPath    string
	Consensus types.ConsesusProtocol
}

type Blockchain struct {
	mu sync.RWMutex

	cons types.ConsesusProtocol

	tail *types.Block // last block, use getter lastBlock()
	db   *db.Database

	blockSubs *subscriptionManager[types.Block]
}

func New(cfg *Config) (*Blockchain, error) {
	logCtx := log.WithField("service", "blockchain")

	logCtx.Info("sync blockchain")
	logCtx.Info("init database")

	database, err := db.New(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("create db: %w", err)
	}

	bc := Blockchain{
		cons:      cfg.Consensus,
		db:        database,
		blockSubs: newSubscription(),
	}

	logCtx.Info("get last block")

	hash, err := database.LastBlockHash()
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, fmt.Errorf("get last block hash: %w", err)
		}

		if err := bc.addGenesisBlock(); err != nil {
			return nil, fmt.Errorf("genesis block: %w", err)
		}
	} else {
		lastBlock, err := bc.BlockByHash(hash)
		if err != nil {
			return nil, fmt.Errorf("get last block: %w", err)
		}
		bc.tail = lastBlock
	}

	return &bc, nil
}

func (bc *Blockchain) AddBlock(txs []*types.Transaction, coinbase types.Address) error {
	newBlock, err := types.NewBlock(bc.cons, txs, bc.lastBlock(), coinbase)
	if err != nil {
		return fmt.Errorf("create new block: %w", err)
	}

	if err = bc.saveNewBlock(newBlock); err != nil {
		return fmt.Errorf("save new block %d: %w", newBlock.Number, err)
	}

	bc.blockSubs.notify(newBlock)
	return nil
}

func (bc *Blockchain) ValidateTx(tx *types.Transaction) bool {
	return false
}

var genesisData = []byte("Genesis")

func (bc *Blockchain) addGenesisBlock() error {
	bc.tail = &types.Block{Number: -1}
	return bc.AddBlock(genesisData)
}

func (bc *Blockchain) saveNewBlock(block *types.Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.tail.Hash != block.PrevBlockHash {
		return fmt.Errorf("parent: %d, block: %d: %w",
			bc.tail.Number, block.Number, ErrInvalidParentBlock)
	}

	if err := bc.db.SaveBlock(block); err != nil {
		return err
	}

	bc.tail = block
	return nil
}

func (bc *Blockchain) lastBlock() *types.Block {
	bc.mu.RLock()
	block := bc.tail
	bc.mu.RUnlock()

	return block
}
