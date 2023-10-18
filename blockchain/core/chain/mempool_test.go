package chain

import (
	"strings"
	"testing"

	"kobla/blockchain/core/types"

	"github.com/stretchr/testify/assert"
)

func TestMempool(t *testing.T) {
	m := newMempool()

	var txs []*types.Transaction
	for i := 0; i < 5; i++ {
		from, err := types.NewAccount()
		assert.NoError(t, err)
		to, err := types.NewAccount()
		assert.NoError(t, err)

		tx := types.NewTransaction(
			from.Address(),
			to.Address(),
			uint64(100*i),
			[]byte(strings.Repeat("test", i)),
		)
		err = tx.WithSignature(from)
		assert.NoError(t, err)

		txs = append(txs, tx)
		m.add(tx)
	}

	assert.Equal(t, m.size(), 5)
	assert.Subset(t, txs[3:], m.top(2))

	m.remove(txs[0].Hash)
	assert.Equal(t, m.size(), 4)
	assert.NotContains(t, m.top(4), txs[0].Hash)

	assert.Equal(t, txs[2], m.get(txs[2].Hash))
}
