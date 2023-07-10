package chain

import (
	"kobla/blockchain/core/types"
	"sync"
)

type Copier[T any] interface {
	Copy() *T
}

type SubscriptionID int

type subscriptionManager[T types.Block] struct {
	mu   sync.RWMutex
	subs map[SubscriptionID]chan *T
	id   SubscriptionID
}

func newSubscription[T types.Block]() *subscriptionManager[T] {
	return &subscriptionManager[T]{
		subs: make(map[SubscriptionID]chan *T),
	}
}

func (s *subscriptionManager[T]) subscribe(ch chan *T) SubscriptionID {
	s.mu.Lock()
	id := s.id
	s.id++
	s.subs[id] = ch
	s.mu.Unlock()

	return id
}

func (s *subscriptionManager[T]) unsubscribe(id SubscriptionID) {
	s.mu.Lock()
	delete(s.subs, id)
	s.mu.Unlock()
}

func (s *subscriptionManager[T]) notify(value Copier[T]) {
	s.mu.RLock()
	for _, ch := range s.subs {
		select {
		case ch <- value.Copy():
		default:
		}
	}
	s.mu.RUnlock()
}
