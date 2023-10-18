//go:build pow

package types_test

import (
	"crypto/rand"
	"testing"

	"kobla/blockchain/core/consensus/pow"
	"kobla/blockchain/core/types"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	account, err := types.NewAccount()
	assert.NoError(t, err)

	data := make([]byte, 100)
	_, err = rand.Read(data)
	assert.NoError(t, err)
	prevBlock := types.Block{Number: 10, Hash: types.NewHash(data)}

	block, err := types.NewBlock(pow.New(), []*types.Transaction{}, &prevBlock, account.Address())
	assert.NoError(t, err)

	assert.Equal(t, block.Hash, block.CalcHash())
	assert.NotEmpty(t, block.Timestamp)
	assert.Equal(t, block.Number, prevBlock.Number+1)
	assert.Equal(t, block.PrevBlockHash, prevBlock.Hash)
	assert.Equal(t, block.Coinbase, account.Address())

	data, err = block.Serialize()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	block2, err := types.DeserializeBlock(data)
	assert.NoError(t, err)
	assert.NotNil(t, block2)
	assert.Equal(t, *block, *block2)
}
