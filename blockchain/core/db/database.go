package db

import (
	"errors"
	"fmt"
	"path2perpetuity/blockchain/core/types"

	badger "github.com/dgraph-io/badger/v3"
)

var ErrNotFound = errors.New("not found")

var lastBlockHashKey = []byte("l")

type Database struct {
	cli *badger.DB
}

func New(path string) (*Database, error) {
	cli, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return &Database{cli: cli}, nil
}

func (db *Database) LastBlockHash() (types.Hash, error) {
	var data []byte

	err := db.cli.View(func(txn *badger.Txn) error {
		item, err := txn.Get(lastBlockHashKey)
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
		if errors.Is(err, badger.ErrKeyNotFound) {
			err = ErrNotFound
		}

		return types.EmptyHash, err
	}

	return types.HashFromSlice(data), nil
}

func (db *Database) SaveBlock(block *types.Block) error {
	return nil
}

func (db *Database) Close() {
	db.cli.Close()
}
