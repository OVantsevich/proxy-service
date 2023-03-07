// Package service price service
package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/OVantsevich/proxy-service/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// bufferSize number of messages stored for every grpc stream
const bufferSize = 1000

// PriceRepository repository interface for price service
//
//go:generate mockery --name=PriceRepository --case=underscore --output=./mocks
type PriceRepository interface {
	GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error)

	GetPrices() ([]*model.Price, error)
	UpdateSubscription(names []string) error
}

// ListenersRepository repository of channels from websocket to stream
//
//go:generate mockery --name=ListenersRepository --case=underscore --output=./mocks
type ListenersRepository interface {
	GetPrices() []string
	Send(prices []*model.Price)
	Update(streamID uuid.UUID, streamChan chan *model.Price, prices []string)
	Delete(streamID uuid.UUID)
}

// Price service
type Price struct {
	priceRepository PriceRepository

	lisRepos ListenersRepository
	sMap     sync.Map
}

// NewPriceService new price service
func NewPriceService(ctx context.Context, rps PriceRepository, lr ListenersRepository) *Price {
	price := &Price{priceRepository: rps, lisRepos: lr}
	price.cycle(ctx)
	return price
}

// GetCurrentPrices get current prices from price service
func (p *Price) GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error) {
	return p.priceRepository.GetCurrentPrices(ctx, names)
}

// GetPrices get price update from price service
func (p *Price) GetPrices() ([]*model.Price, error) {
	return p.priceRepository.GetPrices()
}

// Subscribe allocating new channel for grpc stream with id and returning it
func (p *Price) Subscribe(streamID uuid.UUID) chan *model.Price {
	streamChan := make(chan *model.Price, bufferSize)
	p.sMap.Store(streamID, streamChan)
	return streamChan
}

// UpdateSubscription update list of price for subscriptions
func (p *Price) UpdateSubscription(socketID uuid.UUID, names []string) error {
	streamChan, ok := p.sMap.Load(socketID)
	if !ok {
		return fmt.Errorf("not found")
	}
	p.lisRepos.Delete(socketID)
	p.lisRepos.Update(socketID, streamChan.(chan *model.Price), names)
	err := p.priceRepository.UpdateSubscription(p.lisRepos.GetPrices())
	if err != nil {
		return fmt.Errorf("price - UpdateSubscription - UpdateSubscription: %w", err)
	}
	return nil
}

// DeleteSubscription delete websocket subscription and close it's chan
func (p *Price) DeleteSubscription(streamID uuid.UUID) error {
	streamChan, ok := p.sMap.Load(streamID)
	if !ok {
		return fmt.Errorf("not found")
	}
	p.lisRepos.Delete(streamID)
	close(streamChan.(chan *model.Price))
	p.sMap.Delete(streamID)
	return nil
}

// cycle getting data from price service and sending to websocket subscribers
func (p *Price) cycle(ctx context.Context) {
	var prices []*model.Price
	var err error

	for {
		select {
		case <-ctx.Done():
			return
		default:
			prices, err = p.priceRepository.GetPrices()
			if err != nil {
				logrus.Fatalf("prices - cycle - GetPrices: %v", err)
				return
			}
			p.lisRepos.Send(prices)
		}
	}
}
