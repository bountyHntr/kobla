package db

import (
	"errors"
	"fmt"
	"kobla/blockchain/core/common"
	"kobla/blockchain/core/types"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
)

var ErrNotFound = errors.New("not found")

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
		item, err := txn.Get(hash.Bytes())
		if err != nil {
			return fmt.Errorf("get value: %w", err)
		}

		data, err = item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("copy value: %w", err)
		}

		return nil
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

func (db *Database) BlockHash(number int64) (types.Hash, error) {
	data := make([]byte, types.HashBytes)

	err := db.cli.View(func(txn *badger.Txn) error {
		item, err := txn.Get(common.Int64ToBytes(number))
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

func (db *Database) Close() {
	db.cli.Close()
}

func checkNotFoundError(err error) error {
	if errors.Is(err, badger.ErrKeyNotFound) {
		return ErrNotFound
	}

	return err
}
