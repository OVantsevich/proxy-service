// Package repository price service
package repository

import (
	"context"
	"fmt"

	tsProto "github.com/OVantsevich/Trading-Service/proto"
	"github.com/OVantsevich/proxy-service/internal/model"
)

// TradingService entity
type TradingService struct {
	context.Context
	client tsProto.TradingServiceClient
}

// NewTradingServiceRepository trading service repository constructor
func NewTradingServiceRepository(trp tsProto.TradingServiceClient) (*TradingService, error) {
	ps := &TradingService{client: trp}
	return ps, nil
}

// OpenPosition create position
func (t *TradingService) OpenPosition(ctx context.Context, position *model.Position) (*model.Position, error) {
	resp, err := t.client.OpenPosition(ctx, &tsProto.OpenPositionRequest{
		UserID:        position.User,
		Name:          position.Name,
		Amount:        position.Amount,
		ShortPosition: position.ShortPosition,
	})
	if err != nil {
		return nil, fmt.Errorf("tradingService - OpenPosition - OpenPosition: %w", err)
	}
	return positionFromGRPC(resp.Position), nil
}

// GetPositionByID get pos by ID
func (t *TradingService) GetPositionByID(ctx context.Context, positionID string) (*model.Position, error) {
	resp, err := t.client.GetPositionByID(ctx, &tsProto.GetPositionByIDRequest{
		PositionID: positionID,
	})
	if err != nil {
		return nil, fmt.Errorf("tradingService - GetPositionByID - GetPositionByID: %w", err)
	}
	return positionFromGRPC(resp.Position), nil
}

// GetUserPositions get all user positions
func (t *TradingService) GetUserPositions(ctx context.Context, userID string) ([]*model.Position, error) {
	resp, err := t.client.GetUserPositions(ctx, &tsProto.GetUserPositionsRequest{
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("tradingService - GetUserPositions - GetUserPositions: %w", err)
	}

	response := make([]*model.Position, len(resp.Position))
	for i, p := range resp.Position {
		response[i] = positionFromGRPC(p)
	}

	return response, nil
}

// SetStopLoss set stop loss for user position
func (t *TradingService) SetStopLoss(ctx context.Context, positionID string, stopLoss float64) error {
	_, err := t.client.StopLoss(ctx, &tsProto.StopLossRequest{
		PositionID: positionID,
		Price:      stopLoss,
	})
	if err != nil {
		return fmt.Errorf("tradingService - SetStopLoss - SetStopLoss: %w", err)
	}
	return nil
}

// SetTakeProfit set take profit for user position
func (t *TradingService) SetTakeProfit(ctx context.Context, positionID string, takeProfit float64) error {
	_, err := t.client.TakeProfit(ctx, &tsProto.TakeProfitRequest{
		PositionID: positionID,
		Price:      takeProfit,
	})
	if err != nil {
		return fmt.Errorf("tradingService - SetTakeProfit - SetTakeProfit: %w", err)
	}
	return nil
}

// ClosePosition close user position
func (t *TradingService) ClosePosition(ctx context.Context, positionID string) error {
	_, err := t.client.ClosePosition(ctx, &tsProto.ClosePositionRequest{
		PositionID: positionID,
	})
	if err != nil {
		return fmt.Errorf("tradingService - ClosePosition - ClosePosition: %w", err)
	}
	return nil
}

func positionFromGRPC(pos *tsProto.Position) *model.Position {
	modelPos := &model.Position{
		ID:            pos.Id,
		Name:          pos.Name,
		Amount:        pos.Amount,
		Closed:        pos.Closed,
		ShortPosition: pos.ShortPosition,
		SellingPrice:  pos.SellingPrice,
		PurchasePrice: pos.PurchasePrice,
	}
	if pos.StopLoss != nil {
		modelPos.StopLoss = *pos.StopLoss
	}
	if pos.TakeProfit != nil {
		modelPos.TakeProfit = *pos.TakeProfit
	}
	return modelPos
}
