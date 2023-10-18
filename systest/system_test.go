// go:build pow

package systest

import (
	"fmt"
	"os"
	"testing"

	"kobla/blockchain/core/chain"
	"kobla/blockchain/core/consensus/pow"
	"kobla/blockchain/core/types"

	"github.com/stretchr/testify/assert"
)

func TestSystem(t *testing.T) {
	defer func() {
		r := recover()
		assert.Empty(t, r)
	}()

	var node []*chain.Blockchain
	for i := 1; i <= 2; i++ {
		node = append(node, newNode(i))
		defer dropDB(i)
	}

	var acc []types.Account
	for i := 0; i < 3; i++ {
		account, err := types.NewAccount()
		assert.NoError(t, err)
		acc = append(acc, account)
	}

	tx := types.NewTransaction(acc[0].Address(), acc[1].Address(), 100, []byte("Hello Bob"))
	err := tx.WithSignature(acc[0])
	assert.NoError(t, err)

	ch := make(chan *types.Void)
	subID := node[0].SubscribeMempoolUpdates(ch)
	defer node[0].UnsubscribeMempoolUpdates(subID)

	node[1].SendTx(tx)
	assert.Equal(t, node[1].MempoolSize(), 1)

	assert.Equal(t, node[0].MempoolSize(), 0)
	<-ch
	assert.Equal(t, node[0].MempoolSize(), 1)

	txs := node[0].TopMempoolTxs(1)
	assert.NoError(t, err)
	assert.Equal(t, tx, txs[0])
	tx2, err := node[0].TxByHashFromMempool(tx.Hash)
	assert.NoError(t, err)
	assert.Equal(t, tx, tx2)

	block, err := node[1].BlockByNumber(0)
	assert.NoError(t, err)
	assert.NotEmpty(t, block)
	block2, err := node[0].BlockByHash(block.Hash)
	assert.NoError(t, err)
	assert.Equal(t, block, block2)
	tx2, err = node[0].TxByHash(block.Transactions[0].Hash)
	assert.NoError(t, err)
	assert.Equal(t, block.Transactions[0], tx2)

	err = node[1].MineBlock([]*types.Transaction{tx}, acc[1].Address())
	assert.NoError(t, err)
	assert.Equal(t, node[1].MempoolSize(), 0)
	assert.Equal(t, node[0].MempoolSize(), 1)
	<-ch
	assert.Equal(t, node[0].MempoolSize(), 0)

	block, err = node[0].BlockByNumber(1)
	assert.NoError(t, err)
	assert.Equal(t, block.Transactions[0].Hash, tx.Hash)
}

func newNode(id int) *chain.Blockchain {
	node, err := chain.New(&chain.Config{
		DBPath:    buildDBPath(id),
		Consensus: pow.New(),
		Url:       fmt.Sprintf("localhost:809%d", id),
		Genesis:   true,
		Nodes:     []string{fmt.Sprintf("localhost:809%d", id-1)},
	})
	if err != nil {
		panic(err)
	}
	return node
}

func dropDB(id int) {
	if err := os.RemoveAll(buildDBPath(id)); err != nil {
		panic(err)
	}
}

func buildDBPath(id int) string {
	return fmt.Sprintf("./db%d", id)
}
