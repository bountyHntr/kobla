package types_test

import (
	"testing"

	"kobla/blockchain/core/types"

	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	from, err := types.NewAccount()
	assert.NoError(t, err)
	to, err := types.NewAccount()
	assert.NoError(t, err)

	tx := types.NewTransaction(from.Address(), to.Address(), 100, []byte("Hello"))
	assert.EqualValues(t, 100, tx.Amount)
	assert.Equal(t, from.Address(), tx.Sender)
	assert.Equal(t, to.Address(), tx.Receiver)
	assert.Equal(t, []byte("Hello"), tx.Data)
	assert.Zero(t, tx.Status)
	assert.NotEmpty(t, tx.Hash)
	assert.Empty(t, tx.Signature)

	err = tx.WithSignature(from)
	assert.NoError(t, err)
	assert.NotEmpty(t, tx.Signature)

	data, err := tx.Serialize()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	tx2, err := types.DeserializeTx(data)
	assert.NoError(t, err)
	assert.Equal(t, tx, tx2)

	msg := tx.ToProto()
	assert.NoError(t, err)
	assert.NotEmpty(t, msg)
	tx2, err = types.TransactionFromProto(msg)
	assert.NoError(t, err)
	assert.Equal(t, tx, tx2)

	assert.EqualValues(t, 3330, tx.Cost())
	tx2 = tx.Copy()
	tx2.Sender = to.Address()
	assert.NotEqual(t, tx.Sender, tx2.Sender)
}
