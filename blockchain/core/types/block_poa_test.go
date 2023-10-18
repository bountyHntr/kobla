//go:build poa

package types_test

import (
	"crypto/rand"
	"testing"

	"kobla/blockchain/core/consensus/poa"
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

	var validators []string
	for i := 0; i < 3; i++ {
		acc, err := types.NewAccount()
		assert.NoError(t, err)
		validators = append(validators, acc.Address().String())
	}

	cons, err := poa.New(validators, account.PrivateKey())
	assert.NoError(t, err)

	block, err := types.NewBlock(cons, []*types.Transaction{}, &prevBlock, account.Address())
	assert.NoError(t, err)

	assert.Equal(t, block.Hash, block.CalcHash())
	assert.NotEmpty(t, block.Timestamp)
	assert.Equal(t, block.Number, prevBlock.Number+1)
	assert.Equal(t, block.PrevBlockHash, prevBlock.Hash)
	assert.Equal(t, block.Coinbase, account.Address())
	assert.NotEmpty(t, block.Signature)
	ok, err := account.Verify(block.Hash, block.Signature)
	assert.NoError(t, err)
	assert.True(t, ok)

	data, err = block.Serialize()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	block2, err := types.DeserializeBlock(data)
	assert.NoError(t, err)
	assert.NotNil(t, block2)
	assert.Equal(t, *block, *block2)
}
