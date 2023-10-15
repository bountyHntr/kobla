package types_test

import (
	"crypto/rand"
	"kobla/blockchain/core/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	acc, err := types.NewAccount()
	assert.NoError(t, err)

	address := acc.Address()
	assert.NotNil(t, address)
	assert.NotNil(t, address.String())

	privKey := acc.PrivateKey()
	assert.NotNil(t, privKey)

	_, err = types.AccountFromPrivKey("asfd")
	assert.Error(t, err)
	acc2, err := types.AccountFromPrivKey(privKey)
	assert.NoError(t, err)

	assert.Equal(t, acc, acc2)

	var h types.Hash
	n, err := rand.Read(h[:])
	assert.NoError(t, err)
	assert.Equal(t, types.HashBytes, n)

	signature, err := acc.Sign(h)
	assert.NoError(t, err)
	assert.NotEmpty(t, signature)

	ok, err := acc.Verify(h, signature)
	assert.NoError(t, err)
	assert.True(t, ok)

	signature[0] = signature[0] + 1
	ok, err = acc.Verify(h, signature)
	assert.NoError(t, err)
	assert.False(t, ok)
}
