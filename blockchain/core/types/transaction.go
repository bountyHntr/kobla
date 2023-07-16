package types

import (
	"fmt"
	"kobla/blockchain/core/pb"

	"google.golang.org/protobuf/proto"
)

const byteCost = 10

type TxStatus int32

const (
	TxFail TxStatus = iota
	TxSuccess
)

type Transaction struct {
	Sender   Address
	Receiver Address
	Amount   uint64
	Data     []byte
	Hash     Hash
	Status   TxStatus
}

func NewTransaction(from, to Address, amount uint64, data []byte) (*Transaction, error) {
	tx := Transaction{
		Sender:   from,
		Receiver: to,
		Amount:   amount,
		Data:     data,
	}

	if err := tx.setHash(); err != nil {
		return nil, err
	}

	return &tx, nil
}

func (tx *Transaction) Serialize() ([]byte, error) {
	return proto.Marshal(tx.ToProto())
}

func DeserializeTx(data []byte) (*Transaction, error) {
	var pbTx pb.Transaction
	if err := proto.Unmarshal(data, &pbTx); err != nil {
		return nil, fmt.Errorf("proto unmarshal: %w", err)
	}

	return TransactionFromProto(&pbTx)
}

func (tx *Transaction) ToProto() *pb.Transaction {
	return &pb.Transaction{
		Sender:   tx.Sender.Bytes(),
		Receiver: tx.Receiver.Bytes(),
		Amount:   tx.Amount,
		Data:     tx.Data,
		Hash:     tx.Hash[:],
		Status:   pb.TxStatus(tx.Status),
	}
}

func TransactionFromProto(pbTx *pb.Transaction) (*Transaction, error) {
	return &Transaction{
		Sender:   AddressFromBytes(pbTx.Sender),
		Receiver: AddressFromBytes(pbTx.Receiver),
		Data:     pbTx.Data,
		Hash:     HashFromSlice(pbTx.Hash),
		Status:   TxStatus(pbTx.Status),
	}, nil
}

func (tx *Transaction) Copy() *Transaction {
	data := make([]byte, len(tx.Data))
	copy(data, tx.Data)

	txCopy := *tx
	txCopy.Data = data
	return &txCopy
}

func (tx *Transaction) Cost() uint64 {
	return uint64(AddressLength*2+8+len(tx.Data)+HashBytes) * byteCost
}

func (tx *Transaction) setHash() error {
	data, err := tx.Serialize()
	if err != nil {
		return fmt.Errorf("serialize tx: %w", err)
	}

	tx.Hash = NewHash(data)
	return nil
}
