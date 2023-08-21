package chain

import (
	"container/list"
	"kobla/blockchain/core/types"
	"sync"

	log "github.com/sirupsen/logrus"
)

const defaultMempoolCap = 64

type memoryPool struct {
	mu         sync.RWMutex
	orderedTxs *list.List
	txs        map[types.Hash]*list.Element
	subs       *subscriptionManager[types.Void]
}

func newMempool() *memoryPool {
	return &memoryPool{
		orderedTxs: list.New(),
		txs:        make(map[types.Hash]*list.Element, defaultMempoolCap),
		subs:       newSubscription[types.Void](),
	}
}

func (mp *memoryPool) add(tx *types.Transaction) bool {

	ok, err := tx.Sender.Verify(tx.Hash, tx.Signature)
	if err != nil || !ok {
		log.WithField("service", "mempool").
			WithError(err).
			WithField("ok", ok).
			WithField("hash", tx.Hash.String()).
			Debug("skip")

		return false
	}

	mp.mu.Lock()
	defer mp.mu.Unlock()

	if _, ok := mp.txs[tx.Hash]; ok {
		return false
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
	mp.subs.notify(types.Void{})
	return true
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
	mp.subs.notify(types.Void{})
}

func (mp *memoryPool) size() int {
	mp.mu.RLock()
	size := len(mp.txs)
	mp.mu.RUnlock()

	return size
}

func (mp *memoryPool) top(n int) (txs []*types.Transaction) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	tx := mp.orderedTxs.Front()
	if tx == nil {
		return
	}

	for i := 0; i < n; i++ {
		txs = append(txs, tx.Value.(*types.Transaction))

		if tx = tx.Next(); tx == nil {
			break
		}
	}

	return
}

func (mp *memoryPool) get(hash types.Hash) *types.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	txElement := mp.txs[hash]
	if txElement == nil {
		return nil
	}

	return txElement.Value.(*types.Transaction)
}

func (mp *memoryPool) subscribeUpdates(subCh chan *types.Void) SubscriptionID {
	return mp.subs.subscribe(subCh)
}

func (mp *memoryPool) unsubscribeUpdates(id SubscriptionID) {
	mp.subs.unsubscribe(id)
}
