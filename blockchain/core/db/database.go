package db

import (
	"errors"
	"fmt"
	"path2perpetuity/blockchain/core/common"
	"sync"

	"github.com/prologic/bitcask"
)

var ErrNotFound = errors.New("not found")

var lastBlockHashKey = []byte("l")

type Database struct {
	mu sync.RWMutex
	db *bitcask.Bitcask
}

func New(path string) (*Database, error) {
	db, err := bitcask.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return &Database{db: db}, nil
}

func (db *Database) LastBlockHash() (common.Hash, error) {
	data, err := db.get(lastBlockHashKey)
	if err != nil {
		return common.EmptyHash, nil
	}

	return common.HashFromSlice(data), nil
}

func (db *Database) UpdateLastBlockHash(hash common.Hash) error {
	return db.put(lastBlockHashKey, hash[:])
}

func (db *Database) GetByHash(hash common.Hash) ([]byte, error) {
	return db.get(hash[:])
}

func (db *Database) PutByHash(hash common.Hash, data []byte) error {
	return db.put(hash[:], data)
}

func (db *Database) get(key []byte) ([]byte, error) {
	db.mu.RLock()
	data, err := db.db.Get(key)
	db.mu.RUnlock()

	if err != nil {
		if errors.Is(err, bitcask.ErrKeyNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get from db: %w", err)
	}

	return data, nil
}

func (db *Database) put(key, value []byte) error {
	db.mu.Lock()
	err := db.db.Put(key, value)
	db.mu.Unlock()

	return err
}
