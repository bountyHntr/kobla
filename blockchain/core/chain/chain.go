package chain

import (
	"errors"
	"fmt"
	"path2perpetuity/blockchain/core/db"
	"sync"
)

type Config struct {
	DBPath string
}

type Blockchain struct {
	mu   sync.RWMutex
	tail *Block // last block, use getter lastBlock()
	db   *db.Database
}

func New(cfg *Config) (*Blockchain, error) {
	database, err := db.New(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("create db: %w", err)
	}

	bc := Blockchain{
		db: database,
	}

	if _, err = database.LastBlockHash(); err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, fmt.Errorf("get last block: %w", err)
		}

		genesisBlock, err := NewGenesisBlock()
		if err != nil {
			return nil, fmt.Errorf("genesis block: %w", err)
		}

		if err = bc.saveNewBlock(genesisBlock); err != nil {
			return nil, fmt.Errorf("save genesis block: %w", err)
		}
	}

	return &bc, nil
}

func (bc *Blockchain) AddBlock(data []byte) error {
	newBlock, err := NewBlock(data, bc.lastBlock())
	if err != nil {
		return fmt.Errorf("create new block: %w", err)
	}

	if err = bc.saveNewBlock(newBlock); err != nil {
		return fmt.Errorf("save new block %d: %w", newBlock.Number, err)
	}

	return nil
}

func (bc *Blockchain) saveNewBlock(block *Block) error {
	data, err := block.Serialize()
	if err != nil {
		return fmt.Errorf("serialize: %w", err)
	}

	if err = bc.db.PutByHash(block.Hash, data); err != nil {
		return fmt.Errorf("block: %w", err)
	}

	if err = bc.db.UpdateLastBlockHash(block.Hash); err != nil {
		return fmt.Errorf("last block hash: %w", err)
	}

	bc.mu.Lock()
	bc.tail = block
	bc.mu.Unlock()

	return nil
}

// copy value
func (bc *Blockchain) lastBlock() *Block {
	bc.mu.RLock()
	block := *bc.tail
	bc.mu.RUnlock()

	return &block
}
