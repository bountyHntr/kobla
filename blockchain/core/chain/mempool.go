package chain

import (
	"container/list"
	"kobla/blockchain/core/types"
	"sync"
)

const (
	defaultMempoolCap = 64
	defaultBlockSize  = 5
)

type memoryPool struct {
	mu         sync.RWMutex
	orderedTxs *list.List
	txs        map[types.Hash]*list.Element
}

func newMempool() *memoryPool {
	return &memoryPool{
		orderedTxs: list.New(),
		txs:        make(map[types.Hash]*list.Element, defaultMempoolCap),
	}
}

func (mp *memoryPool) add(tx *types.Transaction) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if _, ok := mp.txs[tx.Hash]; ok {
		return
	}

	var (
		cost      = tx.Cost()
		txElement *list.Element
	)
	for t := mp.orderedTxs.Front(); t != nil; t = t.Next() {
		tCost := t.Value.(*types.Transaction).Cost()
		if tCost < cost {
			txElement = mp.orderedTxs.InsertBefore(tx, t)
			break
		}
	}

	if txElement == nil {
		txElement = mp.orderedTxs.PushBack(tx)
	}

	mp.txs[tx.Hash] = txElement
}

func (mp *memoryPool) remove(hash types.Hash) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	txElement, ok := mp.txs[hash]
	if !ok {
		return
	}

	mp.orderedTxs.Remove(txElement)
	delete(mp.txs, hash)
}

func (mp *memoryPool) contains(hash types.Hash) bool {
	mp.mu.RLock()
	_, contains := mp.txs[hash]
	mp.mu.Unlock()

	return contains
}

func (mp *memoryPool) size() int {
	mp.mu.RLock()
	size := len(mp.txs)
	mp.mu.RUnlock()

	return size
}

func (mp *memoryPool) top(n int) (txs []*types.Transaction) {
	mp.mu.RLock()

	tx := mp.orderedTxs.Front()
	for i := 0; i < n; i++ {
		txs = append(txs, tx.Value.(*types.Transaction))

		if tx = tx.Next(); tx == nil {
			break
		}
	}

	mp.mu.RUnlock()

	return
}
