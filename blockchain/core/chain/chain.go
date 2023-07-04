package chain

import (
	"fmt"
	"sync"
)

type Blockchain struct {
	mu     sync.RWMutex
	blocks []*Block
}

func New() (*Blockchain, error) {
	genesisBlock, err := NewGenesisBlock()
	if err != nil {
		return nil, fmt.Errorf("genesis block: %w", err)
	}

	return &Blockchain{
		blocks: []*Block{genesisBlock},
	}, nil
}

func (bc *Blockchain) AddBlock(data []byte) error {
	newBlock, err := NewBlock(data, bc.lastBlock())
	if err != nil {
		return err
	}

	bc.mu.Lock()
	bc.blocks = append(bc.blocks, newBlock)
	bc.mu.Unlock()

	return nil
}

func (bc *Blockchain) lastBlock() *Block {
	bc.mu.RLock()
	block := bc.blocks[len(bc.blocks)-1]
	bc.mu.RUnlock()

	return block
}
