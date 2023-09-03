//go:build pow

package types

import (
	"bytes"
	"fmt"
	"kobla/blockchain/core/common"
	"kobla/blockchain/core/pb"
	"time"

	"google.golang.org/protobuf/proto"
)

const BlockReward = 10240

type Block struct {
	Timestamp     int64
	Nonce         uint64
	Number        int64
	Transactions  []*Transaction
	PrevBlockHash Hash
	Coinbase      Address
	Hash          Hash
}

var (
	_ Serializable  = &Block{}
	_ Copier[Block] = &Block{}
)

func NewBlock(
	cons ConsesusProtocol,
	txs []*Transaction,
	prevBlock *Block,
	coinbase Address,
) (block *Block, err error) {

	block = &Block{
		Timestamp:     time.Now().Unix(),
		Number:        prevBlock.Number + 1,
		Transactions:  txs,
		PrevBlockHash: prevBlock.Hash,
		Coinbase:      coinbase,
	}

	if err = cons.Run(block); err != nil {
		return nil, fmt.Errorf("Proof-Of-Work: %w", err)
	}

	return block, nil
}

func (b *Block) SetHash() {

	txs := make([]byte, 0, len(b.Transactions)*HashBytes)
	for _, tx := range b.Transactions {
		txs = append(txs, tx.Hash.Bytes()...)
	}

	data := bytes.Join([][]byte{
		common.Int64ToBytes(b.Timestamp),
		common.Int64ToBytes(b.Nonce),
		common.Int64ToBytes(b.Number),
		txs,
		b.PrevBlockHash.Bytes(),
		b.Coinbase.Bytes(),
	}, nil)

	b.Hash = NewHash(data)
}

func (b *Block) Copy() *Block {
	txs := make([]*Transaction, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		txs = append(txs, tx.Copy())
	}

	blockCopy := *b
	blockCopy.Transactions = txs

	return &blockCopy
}

func (b *Block) Serialize() ([]byte, error) {
	txs := make([]*pb.Transaction, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		txs = append(txs, tx.ToProto())
	}

	pbBlock := pb.Block{
		Timestamp:     b.Timestamp,
		Nonce:         b.Nonce,
		Number:        b.Number,
		Transactions:  txs,
		PrevBlockHash: b.PrevBlockHash.Bytes(),
		Hash:          b.Hash.Bytes(),
		Coinbase:      b.Coinbase.Bytes(),
	}

	return proto.Marshal(&pbBlock)
}

func DeserializeBlock(data []byte) (*Block, error) {
	var pbBlock pb.Block
	if err := proto.Unmarshal(data, &pbBlock); err != nil {
		return nil, fmt.Errorf("unmarshal block: %w", err)
	}

	txs := make([]*Transaction, 0, len(pbBlock.Transactions))
	for _, pbTx := range pbBlock.Transactions {
		tx, err := TransactionFromProto(pbTx)
		if err != nil {
			return nil, fmt.Errorf("tx from proto: %w", err)
		}

		txs = append(txs, tx)
	}

	return &Block{
		Timestamp:     pbBlock.Timestamp,
		Nonce:         pbBlock.Nonce,
		Number:        pbBlock.Number,
		Transactions:  txs,
		PrevBlockHash: HashFromSlice(pbBlock.PrevBlockHash),
		Hash:          HashFromSlice(pbBlock.Hash),
		Coinbase:      AddressFromBytes(pbBlock.Coinbase),
	}, nil
}
