// Package service account service
package service

import (
	"context"

	"github.com/OVantsevich/proxy-service/internal/model"
)

// AccountRepository repository interface for account handler
//
//go:generate mockery --name=AccountRepository --case=underscore --output=./mocks
type AccountRepository interface {
	CreateAccount(ctx context.Context, userID string) (*model.Account, error)
	GetAccount(ctx context.Context, userID string) (*model.Account, error)
	IncreaseAmount(ctx context.Context, accountID string, amount float64) error
	DecreaseAmount(ctx context.Context, accountID string, amount float64) error
}

// Account service
type Account struct {
	accountRepository AccountRepository
}

// NewAccountService new account service
func NewAccountService(rps AccountRepository) *Account {
	return &Account{accountRepository: rps}
}

// CreateAccount create account
func (a *Account) CreateAccount(ctx context.Context, userID string) (*model.Account, error) {
	return a.accountRepository.CreateAccount(ctx, userID)
}

// GetAccount get user account
func (a *Account) GetAccount(ctx context.Context, userID string) (*model.Account, error) {
	return a.accountRepository.GetAccount(ctx, userID)
}

// IncreaseAmount increase account amount
func (a *Account) IncreaseAmount(ctx context.Context, accountID string, amount float64) error {
	return a.accountRepository.IncreaseAmount(ctx, accountID, amount)
}

// DecreaseAmount decrease account amount
func (a *Account) DecreaseAmount(ctx context.Context, accountID string, amount float64) error {
	return a.accountRepository.DecreaseAmount(ctx, accountID, amount)
}
