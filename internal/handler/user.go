// Package handler user handler
package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/OVantsevich/proxy-service/internal/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

// passwordStrength strength of password
const passwordStrength = 50

// UserService service interface for user handler
//
//go:generate mockery --name=UserService --case=underscore --output=./mocks
type UserService interface {
	Signup(ctx context.Context, user *model.User) (*model.User, *model.TokenPair, error)
	Login(ctx context.Context, login, password string) (*model.TokenPair, error)
	Refresh(ctx context.Context, userID string, refresh string) (*model.TokenPair, error)

	Update(ctx context.Context, userID string, user *model.User) error
	GetByID(ctx context.Context, userID string) (*model.User, error)
}

// User handler
type User struct {
	userService UserService

	jwtKey string
}

// NewUserHandler new user handler
func NewUserHandler(s UserService, key string) *User {
	return &User{userService: s, jwtKey: key}
}

// SignupRequest signup request
type SignupRequest struct {
	Login    string `json:"login" validate:"required,alphanum,gte=5,lte=20" example:"User123"`
	Email    string `json:"email" validate:"required,email" format:"email" example:"user@usermail.com"`
	Password string `json:"password" validate:"required" example:"strongPassword@123"`
	Name     string `json:"name" validate:"required,alpha,gte=2,lte=25" example:"userNoNum"`
	Age      int32  `json:"age" validate:"required,gte=0,lte=100" example:"20"`
}

// SignupResponse signup response
type SignupResponse struct {
	*model.User
	*model.TokenPair
}

// Signup godoc
//
// @Summary      Add new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user	body    	SignupRequest  true  "New user"
// @Success      201	{object}	SignupResponse
// @Failure      400	{object}	echo.HTTPError
// @Failure      500	{object}	echo.HTTPError
// @Router       /auth/signup [post]
func (u *User) Signup(c echo.Context) (err error) {
	user := &model.User{}
	err = c.Bind(user)
	if err != nil {
		logrus.Error(fmt.Errorf("user - Signup - Bind: %w", err))
		return err
	}

	err = c.Validate(user)
	if err != nil {
		err = fmt.Errorf("user - Signup - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	if err = passwordvalidator.Validate(user.Password, passwordStrength); err != nil {
		err = fmt.Errorf("user - Signup - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var tokenPair *model.TokenPair
	var userResponse *model.User
	userResponse, tokenPair, err = u.userService.Signup(c.Request().Context(), user)
	if err != nil {
		logrus.Error(fmt.Errorf("user - Signup - Signup: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusCreated,
		SignupResponse{
			userResponse,
			tokenPair,
		})
}

// LoginRequest login request
type LoginRequest struct {
	Login    string `json:"login" validate:"required,alphanum,gte=5,lte=20" example:"User123"`
	Password string `json:"password" validate:"required" example:"strongPassword@123"`
}

// Login godoc
//
// @Summary      Login user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        data		body    	LoginRequest	true	"login and password"
// @Success      200		{object}	model.TokenPair
// @Failure      400		{object}	echo.HTTPError
// @Failure      500		{object}	echo.HTTPError
// @Router       /auth/login [post]
func (u *User) Login(c echo.Context) (err error) {
	user := &LoginRequest{}
	err = c.Bind(user)
	if err != nil {
		err = fmt.Errorf("user - Login - Bind: %w", err)
		logrus.Error(err)
		return err
	}

	err = c.Validate(user)
	if err != nil {
		err = fmt.Errorf("user - Login - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var tokenPair *model.TokenPair
	tokenPair, err = u.userService.Login(c.Request().Context(), user.Login, user.Password)
	if err != nil {
		err = fmt.Errorf("user - Login - Login: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, tokenPair)
}

// Refresh godoc
//
// @Summary      Refresh accessToken and refreshToken
// @Tags         users
// @Produce      json
// @Param 		 Cookie 	header 		string  		true	"refresh token"
// @Success      200		{object}	model.TokenPair
// @Failure      500		{object}	echo.HTTPError
// @Router       /auth/refresh [get]
func (u *User) Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh")
	if err != nil {
		logrus.Error(fmt.Errorf("user - Refresh - Cookie: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	refresh := cookie.Value

	var id string
	id, err = idFromToken(refresh, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtKey), nil
	})
	if err != nil {
		logrus.Error(fmt.Errorf("user - Refresh - idFromToken: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var tokenPair *model.TokenPair
	tokenPair, err = u.userService.Refresh(c.Request().Context(), id, refresh)
	if err != nil {
		logrus.Error(fmt.Errorf("user - Refresh - Refresh: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, tokenPair)
}

// UpdateRequest update request
type UpdateRequest struct {
	Email string `json:"email" validate:"required,email" format:"email" example:"user@usermail.com"`
	Name  string `json:"name" validate:"required,alpha,gte=2,lte=25" example:"userNoNum"`
	Age   int32  `json:"age" validate:"required,gte=0,lte=100" example:"20"`
}

// Update godoc
//
// @Summary      Update info about user
// @Tags         users
// @Produce      json
// @Param		 user	body	UpdateRequest	 true	"New user info"
// @Success      200
// @Failure      400	{object}	echo.HTTPError
// @Failure      500	{object}	echo.HTTPError
// @Router       /update [put]
// @Security Bearer
func (u *User) Update(c echo.Context) (err error) {
	user := &UpdateRequest{}
	err = c.Bind(user)
	if err != nil {
		logrus.Error(fmt.Errorf("user - Update - Bind: %w", err))
		return err
	}

	err = c.Validate(user)
	if err != nil {
		err = fmt.Errorf("user - Update - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	id := idFromContext(c)
	err = u.userService.Update(c.Request().Context(), id, &model.User{
		Email: user.Email,
		Name:  user.Name,
		Age:   user.Age,
	})
	if err != nil {
		logrus.Error(fmt.Errorf("user - Update - Update: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, "")
}

// UserByID godoc
//
// @Summary		 getting user by id
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id	 	header   	string		true  "id"
// @Success      200	object		model.User
// @Failure      403	{object}	echo.HTTPError
// @Failure      500	{object}	echo.HTTPError
// @Router       /userByID [get]
// @Security Bearer
func (u *User) UserByID(c echo.Context) (err error) {
	id := c.Request().Header.Get("id")

	var user *model.User
	user, err = u.userService.GetByID(c.Request().Context(), id)
	if err != nil {
		logrus.Error(fmt.Errorf("user - UserByID - GetByID: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, user)
}

func idFromContext(c echo.Context) (ID string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims
	return claims.(*model.CustomClaims).ID
}

func idFromToken(token string, keyFunc func(token *jwt.Token) (interface{}, error)) (id string, err error) {
	claims := &model.CustomClaims{}

	_, err = jwt.ParseWithClaims(
		token,
		claims,
		keyFunc,
	)
	if err != nil {
		err = fmt.Errorf("invalid token: %w", err)
	}

	return claims.ID, err
}
