//go:build pow

package db_test

import (
	"crypto/rand"
	"os"
	"testing"

	"kobla/blockchain/core/consensus/pow"
	"kobla/blockchain/core/db"
	"kobla/blockchain/core/types"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	const dbPath = "./test_db_"
	defer func() {
		err := os.RemoveAll(dbPath)
		assert.NoError(t, err)
	}()

	db, err := db.New(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	acc1, err := types.NewAccount()
	assert.NoError(t, err)
	acc2, err := types.NewAccount()
	assert.NoError(t, err)
	acc3, err := types.NewAccount()
	assert.NoError(t, err)

	tx1 := types.NewTransaction(acc1.Address(), acc2.Address(), 100, []byte("Hello"))
	tx2 := types.NewTransaction(acc3.Address(), acc1.Address(), 20, nil)

	data := make([]byte, 100)
	_, err = rand.Read(data)
	assert.NoError(t, err)
	prevBlock := types.Block{Number: 10, Hash: types.NewHash(data)}

	block, err := types.NewBlock(pow.New(), []*types.Transaction{tx1, tx2}, &prevBlock, acc1.Address())
	assert.NoError(t, err)

	err = db.SaveBlock(&prevBlock)
	assert.NoError(t, err)
	err = db.SaveBlock(block)
	assert.NoError(t, err)

	hash, err := db.LastBlockHash()
	assert.NoError(t, err)
	assert.Equal(t, block.Hash, hash)
	hash, err = db.BlockHash(block.Number)
	assert.NoError(t, err)
	assert.Equal(t, block.Hash, hash)

	lastBlock, err := db.Block(hash)
	assert.NoError(t, err)
	assert.Equal(t, block, lastBlock)

	blockHash, err := db.TxToBlock(tx1.Hash)
	assert.NoError(t, err)
	assert.Equal(t, block.Hash, blockHash)

	balance1, err := db.Balance(acc1.Address())
	assert.NoError(t, err)
	balance2, err := db.Balance(acc2.Address())
	assert.NoError(t, err)
	assert.NotEmpty(t, balance1)
	assert.NotEmpty(t, balance2)
	assert.NotEqual(t, balance1, balance2)
}
