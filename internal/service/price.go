// Package service price service
package service

import (
	"context"
	"fmt"

	"github.com/OVantsevich/proxy-service/internal/model"
)

// PriceRepository repository interface for price service
//
//go:generate mockery --name=PriceRepository --case=underscore --output=./mocks
type PriceRepository interface {
	GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error)

	GetPrices() ([]*model.Price, error)
	UpdateSubscription(names []string) error
}

// Sockets repository interface for price service
//
//go:generate mockery --name=Sockets --case=underscore --output=./mocks
type Sockets interface {
	AddSocket(ctx context.Context, positions []*model.Position, prices map[string]*model.Price) error
	RemoveSocket(position *model.Position) error

	SendPrices(prices []*model.Price)
}

// Price service
type Price struct {
	priceRepository PriceRepository
}

// NewPriceService new price service
func NewPriceService(rps PriceRepository) *Price {
	return &Price{priceRepository: rps}
}

func (p *Price) GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error) {
	return p.priceRepository.GetCurrentPrices(ctx, names)
}

func (p *Price) GetPrices() ([]*model.Price, error) {
	return p.priceRepository.GetPrices()
}

func (p *Price) UpdateSubscription(names []string) error {
	return p.priceRepository.UpdateSubscription(names)
}

func getPricesListener(ctx context.Context, t *Price, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
		default:
			prices, err := t.priceRepository.GetPrices()
			if err != nil {
				errChan <- fmt.Errorf("price - getPricesListener - GetPrices: %w", err)
				continue
			}
		}
	}
}
