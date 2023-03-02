// Package repository repos
package repository

import (
	"sync"

	"github.com/OVantsevich/proxy-service/internal/model"

	"github.com/google/uuid"
)

// Subscribers storing subscribers
type Subscribers map[uuid.UUID]chan *model.Price

// Listeners websocket for grpc stream
type Listeners struct {
	MU     sync.RWMutex
	prices map[string]Subscribers
}

// NewListenersRepository constructor
func NewListenersRepository() *Listeners {
	return &Listeners{prices: make(map[string]Subscribers)}
}

// GetPrices get all currently active prices
func (l *Listeners) GetPrices() []string {
	l.MU.RLock()
	keys := make([]string, 0, len(l.prices))
	for k := range l.prices {
		keys = append(keys, k)
	}
	l.MU.RUnlock()
	return keys
}

// Update add new pairs: price-stream
func (l *Listeners) Update(listenerID uuid.UUID, listenerChan chan *model.Price, prices []string) {
	l.MU.Lock()
	for _, p := range prices {
		cp, ok := l.prices[p]
		if ok {
			cp[listenerID] = listenerChan
			continue
		}
		l.prices[p] = make(Subscribers)
		l.prices[p][listenerID] = listenerChan
	}
	l.MU.Unlock()
}

// Delete remove pairs: price-stream
func (l *Listeners) Delete(listenerID uuid.UUID) {
	l.MU.Lock()
	for _, p := range l.prices {
		delete(p, listenerID)
	}
	l.MU.Unlock()
}

// Send prices to streams
func (l *Listeners) Send(prices []*model.Price) {
	l.MU.RLock()
	for _, p := range prices {
		for _, c := range l.prices[p.Name] {
			select {
			case c <- p:
			default:
				go func(goP *model.Price, goC chan *model.Price) {
					goC <- goP
				}(p, c)
			}
		}
	}
	l.MU.RUnlock()
}
