// Package repository price service
package repository

import (
	"context"
	"fmt"

	usProto "github.com/OVantsevich/User-Service/proto"
	"github.com/OVantsevich/proxy-service/internal/model"
)

// UserService entity
type UserService struct {
	client usProto.UserServiceClient
}

// NewUserServiceRepository user service repository constructor
func NewUserServiceRepository(usp usProto.UserServiceClient) *UserService {
	ps := &UserService{client: usp}
	return ps
}

// Signup user in user service
func (u *UserService) Signup(ctx context.Context, user *model.User) (*model.User, *model.TokenPair, error) {
	resp, err := u.client.Signup(ctx, &usProto.SignupRequest{
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Age:      user.Age,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("userService - Signup - Signup: %w", err)
	}
	return userFromGRPC(resp.User), &model.TokenPair{
		Access:  resp.AccessToken,
		Refresh: resp.RefreshToken,
	}, nil
}

// Login in user service
func (u *UserService) Login(ctx context.Context, login, password string) (*model.TokenPair, error) {
	resp, err := u.client.Login(ctx, &usProto.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("userService - Login - Login: %w", err)
	}
	return &model.TokenPair{
		Access:  resp.AccessToken,
		Refresh: resp.RefreshToken,
	}, nil
}

// Refresh refresh token pair
func (u *UserService) Refresh(ctx context.Context, userID, refresh string) (*model.TokenPair, error) {
	resp, err := u.client.Refresh(ctx, &usProto.RefreshRequest{
		Id:           userID,
		RefreshToken: refresh,
	})
	if err != nil {
		return nil, fmt.Errorf("userService - Refresh - Refresh: %w", err)
	}
	return &model.TokenPair{
		Access:  resp.AccessToken,
		Refresh: resp.RefreshToken,
	}, nil
}

// Update user data
func (u *UserService) Update(ctx context.Context, userID string, user *model.User) error {
	_, err := u.client.Update(ctx, &usProto.UpdateRequest{
		ID:    userID,
		Email: user.Email,
		Name:  user.Name,
		Age:   user.Age,
	})
	if err != nil {
		return fmt.Errorf("userService - Update - Update: %w", err)
	}
	return nil
}

// GetByID get user by ID
func (u *UserService) GetByID(ctx context.Context, userID string) (*model.User, error) {
	resp, err := u.client.UserById(ctx, &usProto.UserByIdRequest{
		ID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("userService - UserById - UserById: %w", err)
	}
	return userFromGRPC(resp.User), nil
}

func userFromGRPC(user *usProto.User) *model.User {
	return &model.User{
		ID:    user.Id,
		Login: user.Login,
		Email: user.Email,
		Name:  user.Name,
		Age:   user.Age,
	}
}
