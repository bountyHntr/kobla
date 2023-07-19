package db

import (
	"errors"
	"fmt"
	"kobla/blockchain/core/common"
	"kobla/blockchain/core/types"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidBalance  = errors.New("invalid account balance")
	ErrZeroAddressScam = errors.New("zero address scam")
)

var lastBlockHashKey = []byte("l")

type Database struct {
	cli *badger.DB
}

func New(path string) (*Database, error) {
	opts := badger.DefaultOptions(path).
		WithLoggingLevel(badger.WARNING).
		WithCompression(options.None)

	cli, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return &Database{cli: cli}, nil
}

func (db *Database) LastBlockHash() (types.Hash, error) {
	data := make([]byte, types.HashBytes)

	err := db.cli.View(func(txn *badger.Txn) error {
		item, err := txn.Get(lastBlockHashKey)
		if err != nil {
			return fmt.Errorf("get value: %w", err)
		}

		_, err = item.ValueCopy(data)
		if err != nil {
			return fmt.Errorf("copy value: %w", err)
		}

		return nil
	})

	if err != nil {
		return types.EmptyHash, checkNotFoundError(err)
	}

	return types.HashFromSlice(data), nil
}

func (db *Database) SaveBlock(block *types.Block) error {
	return db.cli.Update(func(txn *badger.Txn) error {
		if err := txn.Set(lastBlockHashKey, block.Hash[:]); err != nil {
			return fmt.Errorf("update last block hash: %w", err)
		}

		for _, tx := range block.Transactions {
			if err := executeTx(txn, tx, block.Coinbase); err != nil {
				return fmt.Errorf("execute %s tx: %w", tx.Hash.String(), err)
			}

			if err := txn.Set(tx.Hash.Bytes(), block.Hash.Bytes()); err != nil {
				return fmt.Errorf("save reference from tx to block: %w", err)
			}
		}

		blockData, err := block.Serialize()
		if err != nil {
			return fmt.Errorf("serialize block: %w", err)
		}

		if err := txn.Set(block.Hash.Bytes(), blockData); err != nil {
			return fmt.Errorf("save block: %w", err)
		}

		if err := txn.Set(common.Int64ToBytes(block.Number), block.Hash.Bytes()); err != nil {
			return fmt.Errorf("save hash by number: %w", err)
		}

		return nil
	})
}

func (db *Database) Block(hash types.Hash) (*types.Block, error) {
	var data []byte

	err := db.cli.View(func(txn *badger.Txn) error {
		return copyValue(txn, hash.Bytes(), &data)
	})

	if err != nil {
		return nil, checkNotFoundError(err)
	}

	block, err := types.DeserializeBlock(data)
	if err != nil {
		return nil, fmt.Errorf("deserialize block: %w", err)
	}

	return block, nil
}

func (db *Database) TxToBlock(hash types.Hash) (types.Hash, error) {
	data := make([]byte, types.HashBytes)

	err := db.cli.View(func(txn *badger.Txn) error {
		return copyValue(txn, hash.Bytes(), &data)
	})

	if err != nil {
		return types.EmptyHash, checkNotFoundError(err)
	}

	return types.HashFromSlice(data), nil
}

func (db *Database) BlockHash(number int64) (types.Hash, error) {
	data := make([]byte, types.HashBytes)

	err := db.cli.View(func(txn *badger.Txn) error {
		return copyValue(txn, common.Int64ToBytes(number), &data)
	})

	if err != nil {
		return types.EmptyHash, checkNotFoundError(err)
	}

	return types.HashFromSlice(data), nil
}

func (db *Database) Balance(address types.Address) (uint64, error) {
	var data []byte

	err := db.cli.View(func(txn *badger.Txn) error {
		return copyValue(txn, address.Bytes(), &data)
	})

	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return 0, nil
		}

		return 0, err
	}

	return common.Int64FromBytes[uint64](data), nil
}

func (db *Database) Close() {
	db.cli.Close()
}

func executeTx(db *badger.Txn, tx *types.Transaction, coinbase types.Address) (err error) {
	defer func() {
		if err != nil {
			tx.Status = types.TxFail
		} else {
			tx.Status = types.TxSuccess
		}
	}()

	if tx.Sender == types.ZeroAddress && (tx.Receiver != coinbase || tx.Amount != types.BlockReward) {
		return ErrZeroAddressScam
	}

	txCost := tx.Cost()

	if err = changeBalance(db, tx.Sender, txCost+tx.Amount, -1); err != nil {
		return fmt.Errorf("change sender balance: %w", err)
	}

	if tx.Amount > 0 {
		if err = changeBalance(db, tx.Receiver, tx.Amount, 1); err != nil {
			return fmt.Errorf("change receiver balance: %w", err)
		}
	}

	if err = changeBalance(db, coinbase, txCost, 1); err != nil {
		return fmt.Errorf("change coinbase balance: %w", err)
	}

	return nil
}

func checkNotFoundError(err error) error {
	if errors.Is(err, badger.ErrKeyNotFound) {
		return ErrNotFound
	}

	return err
}

func copyValue(txn *badger.Txn, key []byte, data *[]byte) error {
	item, err := txn.Get(key)
	if err != nil {
		return fmt.Errorf("get value: %w", err)
	}

	*data, err = item.ValueCopy(*data)
	if err != nil {
		return fmt.Errorf("copy value: %w", err)
	}

	return nil
}

func changeBalance(db *badger.Txn, address types.Address, diff uint64, sign int8) error {
	if address == types.ZeroAddress {
		return nil
	}

	balanceItem, err := db.Get(address.Bytes())
	if err != nil {
		if !errors.Is(err, badger.ErrKeyNotFound) {
			return fmt.Errorf("get %s balance: %w", address.String(), err)
		}

		if err := db.Set(address.Bytes(), common.Int64ToBytes(types.InitBalance)); err != nil {
			return fmt.Errorf("set init balance: %w", err)
		}

		balanceItem, _ = db.Get(address.Bytes())
	}

	return balanceItem.Value(func(balance []byte) error {
		balanceU64 := common.Int64FromBytes[uint64](balance)

		if sign < 0 {

			if diff > balanceU64 {
				return fmt.Errorf("balance: %d; trying to spend %d: %w", balanceU64, diff, ErrInvalidBalance)
			}

			balanceU64 -= diff
		} else {
			balanceU64 += diff
		}

		if err := db.Set(address.Bytes(), common.Int64ToBytes(balanceU64)); err != nil {
			return fmt.Errorf("save %s balance: %w", address.String(), err)
		}

		return nil
	})
}
