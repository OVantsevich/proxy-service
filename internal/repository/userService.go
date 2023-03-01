// Package repository price service
package repository

import (
	"context"
	usProto "github.com/OVantsevich/User-Service/proto"
	"github.com/OVantsevich/proxy-service/internal/model"
)

// UserService entity
type UserService struct {
	ctx    context.Context
	client usProto.UserServiceClient
}

// NewUserServiceRepository user service repository constructor
func NewUserServiceRepository(ctx context.Context, pspp usProto.UserServiceClient) (*UserService, error) {
	ps := &UserService{client: pspp, ctx: ctx}
	return ps, nil
}

func (u *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	u.client.Signup(ctx, usProto.SignupRequest{
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Age:      user.Age,
	})
}

func (u *UserService) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) UpdateUser(ctx context.Context, userID string, user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) RefreshUser(ctx context.Context, userID, token string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) DeleteUser(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}
