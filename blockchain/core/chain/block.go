package chain

import (
	"fmt"
	"path2perpetuity/blockchain/core/consensus/pow"
	"time"
)

type Block struct {
	Timestamp int64
	Nonce     uint64
	Number    int64

	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func NewBlock(data []byte, prevBlock *Block) (block *Block, err error) {

	block = &Block{
		Timestamp:     time.Now().Unix(),
		Number:        prevBlock.Number + 1,
		Data:          data,
		PrevBlockHash: prevBlock.Hash,
	}

	block.Nonce, block.Hash, err = pow.Run(data)
	if err != nil {
		return nil, fmt.Errorf("Proof-Of-Work: %w", err)
	}

	return block, nil
}

var genesisData = []byte("Genesis")

func NewGenesisBlock() (*Block, error) {
	return NewBlock(genesisData, &Block{Number: -1})
}
