//go:build pow

package pow_test

import (
	"testing"

	"kobla/blockchain/core/consensus/pow"
	"kobla/blockchain/core/types"

	"github.com/stretchr/testify/assert"
)

func TestPow(t *testing.T) {
	block := &types.Block{Number: 10, Hash: types.NewHash([]byte("hello"))}
	cons := pow.New()
	err := cons.Run(block)
	assert.NoError(t, err)

	assert.EqualValues(t, 25, block.Nonce)
	assert.Equal(t, "1CbPz8eqJSXrGPVFg1Bf6LR7ysnbm6x6zmZQU16x5RUqKh8bvUFFFWApkwusCebcyeEr6dDJo9hk9a3rc9peQv8", block.Hash.String())

	ok := cons.Validate(block)
	assert.True(t, ok)
	block.Nonce += 1
	ok = cons.Validate(block)
	assert.False(t, ok)
}
