// Package repository account service
package repository

import (
	"context"
	"fmt"

	psProto "github.com/OVantsevich/Payment-Service/proto"
	"github.com/OVantsevich/proxy-service/internal/model"
)

// PaymentService entity
type PaymentService struct {
	client psProto.PaymentServiceClient
}

// NewPaymentServiceRepository payment service repository constructor
func NewPaymentServiceRepository(psp psProto.PaymentServiceClient) *PaymentService {
	ps := &PaymentService{client: psp}
	return ps
}

// CreateAccount create account
func (p *PaymentService) CreateAccount(ctx context.Context, userID string) (*model.Account, error) {
	resp, err := p.client.CreateAccount(ctx, &psProto.CreateAccountRequest{
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("paymentService - CreateAccount - CreateAccount: %w", err)
	}
	return accountFromGRPC(resp.Account), nil
}

// GetAccount get user account
func (p *PaymentService) GetAccount(ctx context.Context, userID string) (*model.Account, error) {
	resp, err := p.client.GetAccount(ctx, &psProto.GetAccountRequest{
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("paymentService - GetAccount - GetAccount: %w", err)
	}
	return accountFromGRPC(resp.Account), nil
}

// IncreaseAmount increase amount of user account
func (p *PaymentService) IncreaseAmount(ctx context.Context, accountID string, amount float64) error {
	_, err := p.client.IncreaseAmount(ctx, &psProto.AmountRequest{
		Amount:    amount,
		AccountID: accountID,
	})
	if err != nil {
		return fmt.Errorf("paymentService - IncreaseAmount - IncreaseAmount: %w", err)
	}
	return nil
}

// DecreaseAmount decrease amount of user account
func (p *PaymentService) DecreaseAmount(ctx context.Context, accountID string, amount float64) error {
	_, err := p.client.DecreaseAmount(ctx, &psProto.AmountRequest{
		Amount:    amount,
		AccountID: accountID,
	})
	if err != nil {
		return fmt.Errorf("paymentService - DecreaseAmount - DecreaseAmount: %w", err)
	}
	return nil
}

func accountFromGRPC(acc *psProto.Account) *model.Account {
	modelPos := &model.Account{
		ID:     acc.ID,
		Amount: acc.Amount,
		User:   acc.UserID,
	}
	return modelPos
}
