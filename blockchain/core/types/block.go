package types

import (
	"fmt"
	"path2perpetuity/blockchain/core/pb"
	"time"

	"google.golang.org/protobuf/proto"
)

type Block struct {
	Timestamp int64
	Nonce     uint64
	Number    int64

	Data          []byte
	PrevBlockHash Hash
	Hash          Hash
}

func NewBlock(cons ConsesusProtocol, data []byte, prevBlock *Block) (block *Block, err error) {

	block = &Block{
		Timestamp:     time.Now().Unix(),
		Number:        prevBlock.Number + 1,
		Data:          data,
		PrevBlockHash: prevBlock.Hash,
	}

	if err = cons.Run(block); err != nil {
		return nil, fmt.Errorf("Proof-Of-Work: %w", err)
	}

	return block, nil
}

var genesisData = []byte("Genesis")

func NewGenesisBlock(cons ConsesusProtocol) (*Block, error) {
	return NewBlock(cons, genesisData, &Block{Number: -1})
}

func (b *Block) Serialize() ([]byte, error) {
	pbBlock := pb.Block{
		Timestamp:     b.Timestamp,
		Nonce:         b.Nonce,
		Number:        b.Number,
		Data:          b.Data,
		PrevBlockHash: b.PrevBlockHash[:],
		Hash:          b.Hash[:],
	}

	return proto.Marshal(&pbBlock)
}

func DeserializeBlock(data []byte) (*Block, error) {
	var pbBlock pb.Block
	if err := proto.Unmarshal(data, &pbBlock); err != nil {
		return nil, err
	}

	return &Block{
		Timestamp:     pbBlock.Timestamp,
		Nonce:         pbBlock.Nonce,
		Number:        pbBlock.Number,
		Data:          pbBlock.Data,
		PrevBlockHash: HashFromSlice(pbBlock.PrevBlockHash),
		Hash:          HashFromSlice(pbBlock.Hash),
	}, nil
}
