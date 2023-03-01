// Package service trading service
package service

import (
	"context"

	"github.com/OVantsevich/proxy-service/internal/model"
)

// TradingRepository repository interface for trading service
//
//go:generate mockery --name=TradingRepository --case=underscore --output=./mocks
type TradingRepository interface {
	OpenPosition(ctx context.Context, position *model.Position) (*model.Position, error)
	GetPositionByID(ctx context.Context, positionID string) (*model.Position, error)
	GetUserPositions(ctx context.Context, userID string) ([]*model.Position, error)
	SetStopLoss(ctx context.Context, positionID string, stopLoss float64) error
	SetTakeProfit(ctx context.Context, positionID string, takeProfit float64) error
	ClosePosition(ctx context.Context, positionID string) error
}

// Trading service
type Trading struct {
	tradingRepository TradingRepository
}

// NewTradingService new trading service
func NewTradingService(rps TradingRepository) *Trading {
	return &Trading{tradingRepository: rps}
}

// OpenPosition open new position
func (t *Trading) OpenPosition(ctx context.Context, position *model.Position) (*model.Position, error) {
	return t.tradingRepository.OpenPosition(ctx, position)
}

// GetPositionByID position by id
func (t *Trading) GetPositionByID(ctx context.Context, positionID string) (*model.Position, error) {
	return t.tradingRepository.GetPositionByID(ctx, positionID)
}

// GetUserPositions get all users position
func (t *Trading) GetUserPositions(ctx context.Context, userID string) ([]*model.Position, error) {
	return t.tradingRepository.GetUserPositions(ctx, userID)
}

// SetStopLoss set stop loss for selected position
func (t *Trading) SetStopLoss(ctx context.Context, positionID string, stopLoss float64) error {
	return t.tradingRepository.SetStopLoss(ctx, positionID, stopLoss)
}

// SetTakeProfit set take profit for selected position
func (t *Trading) SetTakeProfit(ctx context.Context, positionID string, takeProfit float64) error {
	return t.tradingRepository.SetTakeProfit(ctx, positionID, takeProfit)
}

// ClosePosition close position
func (t *Trading) ClosePosition(ctx context.Context, positionID string) error {
	return t.tradingRepository.ClosePosition(ctx, positionID)
}
