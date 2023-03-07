// Package service user service
package service

import (
	"context"

	"github.com/OVantsevich/proxy-service/internal/model"
)

// UserRepository user repository interface for user service
//
//go:generate mockery --name=UserRepository --case=underscore --output=./mocks
type UserRepository interface {
	Signup(ctx context.Context, user *model.User) (*model.User, *model.TokenPair, error)
	Login(ctx context.Context, login, password string) (*model.TokenPair, error)
	Refresh(ctx context.Context, userID string, refresh string) (*model.TokenPair, error)

	Update(ctx context.Context, userID string, user *model.User) error
	GetByID(ctx context.Context, userID string) (*model.User, error)
}

// User service
type User struct {
	userRepository UserRepository
}

// NewUserService new user service
func NewUserService(rps UserRepository) *User {
	return &User{userRepository: rps}
}

// Signup user Signup
func (u *User) Signup(ctx context.Context, user *model.User) (*model.User, *model.TokenPair, error) {
	return u.userRepository.Signup(ctx, user)
}

// Login user Login
func (u *User) Login(ctx context.Context, login, password string) (*model.TokenPair, error) {
	return u.userRepository.Login(ctx, login, password)
}

// Refresh user Refresh
func (u *User) Refresh(ctx context.Context, userID, refresh string) (*model.TokenPair, error) {
	return u.userRepository.Refresh(ctx, userID, refresh)
}

// Update user Update
func (u *User) Update(ctx context.Context, userID string, user *model.User) error {
	return u.userRepository.Update(ctx, userID, user)
}

// GetByID user GetByID
func (u *User) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return u.userRepository.GetByID(ctx, userID)
}
