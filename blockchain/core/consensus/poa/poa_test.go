//go:build poa

package poa_test

import (
	"testing"

	"kobla/blockchain/core/consensus/poa"
	"kobla/blockchain/core/types"

	"github.com/stretchr/testify/assert"
)

func TestPow(t *testing.T) {
	acc1, err := types.AccountFromPrivKey("2Gd5MkGaHAm9nv4xyeCozZ3L4DmVngBQbj8R7FcLLgpitwZzTN7qw2CXqZFTYKGP6srUvsqdXdeoN4TQWg7dFE2X")
	assert.NoError(t, err)
	acc2, err := types.AccountFromPrivKey("5XFD59e8TqmotcGXNmqT4RfUvUq7qCCvadyCAC3NK8B5FEt2VQJqtvqKBTzyg6pwrLZ476iBreaTCRmq1UrXWnUq")
	assert.NoError(t, err)

	block := &types.Block{Number: 10, Coinbase: acc1.Address()}
	block.Hash = block.CalcHash()
	assert.Equal(t, "vmZQCuEsaBcuNE8kjU5MkfdPSGHLxwTC8tJkMgyKjL9rGEW8cpGgbcUiahcykdfie97VUFC3SjkYkqxhcbir7GQ", block.Hash.String())

	validators := []string{acc1.Address().String(), acc2.Address().String()}
	cons, err := poa.New(validators, acc1.PrivateKey())
	assert.NoError(t, err)

	err = cons.Run(block)
	assert.NoError(t, err)

	assert.NotEmpty(t, block.Signature)
	ok, err := acc1.Verify(block.Hash, block.Signature)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok = cons.Validate(block)
	assert.True(t, ok)

	acc, err := types.NewAccount()
	assert.NoError(t, err)
	block.Coinbase = acc.Address()
	cons, err = poa.New(validators, acc.PrivateKey())
	assert.NoError(t, err)
	err = cons.Run(block)
	assert.NoError(t, err)

	ok = cons.Validate(block)
	assert.False(t, ok)
}
