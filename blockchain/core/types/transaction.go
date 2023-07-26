package types

import (
	"bytes"
	"fmt"
	"kobla/blockchain/core/common"
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
	Sender    Address
	Receiver  Address
	Amount    uint64
	Data      []byte
	Status    TxStatus
	Hash      Hash
	Signature []byte
}

func NewTransaction(from, to Address, amount uint64, data []byte) *Transaction {
	tx := &Transaction{
		Sender:   from,
		Receiver: to,
		Amount:   amount,
		Data:     data,
	}
	tx.setHash()

	return tx
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
		Sender:    tx.Sender.Bytes(),
		Receiver:  tx.Receiver.Bytes(),
		Amount:    tx.Amount,
		Data:      tx.Data,
		Hash:      tx.Hash[:],
		Status:    pb.TxStatus(tx.Status),
		Signature: tx.Signature,
	}
}

func TransactionFromProto(pbTx *pb.Transaction) (*Transaction, error) {
	return &Transaction{
		Sender:    AddressFromBytes(pbTx.Sender),
		Receiver:  AddressFromBytes(pbTx.Receiver),
		Amount:    pbTx.Amount,
		Data:      pbTx.Data,
		Hash:      HashFromSlice(pbTx.Hash),
		Status:    TxStatus(pbTx.Status),
		Signature: pbTx.Signature,
	}, nil
}

func (tx *Transaction) Copy() *Transaction {
	txCopy := *tx
	txCopy.Data = make([]byte, len(tx.Data))
	copy(txCopy.Data, tx.Data)
	txCopy.Signature = make([]byte, len(tx.Signature))
	copy(txCopy.Signature, tx.Signature)

	return &txCopy
}

func (tx *Transaction) Cost() uint64 {
	return uint64(AddressLength*2+8+len(tx.Data)+HashBytes) * byteCost
}

func (tx *Transaction) WithSignature(signer Account) error {
	signature, err := signer.Sign(tx.Hash)
	if err != nil {
		return err
	}

	tx.Signature = signature
	return nil
}

func (tx *Transaction) setHash() {
	data := bytes.Join([][]byte{
		tx.Sender.Bytes(),
		tx.Receiver.Bytes(),
		common.Int64ToBytes(tx.Amount),
		tx.Data,
	}, nil)

	tx.Hash = NewHash(data)
}
