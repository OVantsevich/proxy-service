// Package handler user handler
package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/OVantsevich/proxy-service/internal/model"

	"github.com/go-playground/validator/v10"
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

	val *validator.Validate
}

// NewUserHandler new user handler
func NewUserHandler(s UserService) *User {
	return &User{userService: s, val: validator.New()}
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
// @Param        body	body     model.User  true  "New user object"
// @Success      201	{object}	SignupResponse
// @Failure      400
// @Failure      500
// @Router       /signup [post]
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

// Login godoc
//
// @Summary      Login user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body		body    model.User	true  "login and password"
// @Success      201	{object}	model.TokenPair
// @Failure      500
// @Router       /login [get]
func (u *User) Login(c echo.Context) (err error) {
	user := &model.User{}
	err = c.Bind(user)
	if err != nil {
		logrus.Error(fmt.Errorf("user - Login - Bind: %w", err))
		return err
	}

	var tokenPair *model.TokenPair
	tokenPair, err = u.userService.Login(c.Request().Context(), user.Login, user.Password)
	if err != nil {
		logrus.Error(fmt.Errorf("user - Login - Login: %w", err))
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
// @Success      201	{object}	model.TokenPair
// @Failure      500
// @Router       /refresh [get]
// @Security Bearer
func (u *User) Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh")
	if err != nil {
		logrus.Error(fmt.Errorf("user - Refresh - Cookie: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	var userID string
	err = c.Bind(userID)
	if err != nil {
		logrus.Error(fmt.Errorf("user - Refresh - Bind: %w", err))
		return err
	}

	refresh := cookie.Value

	var tokenPair *model.TokenPair
	tokenPair, err = u.userService.Refresh(c.Request().Context(), userID, refresh)
	if err != nil {
		logrus.Error(fmt.Errorf("user - Refresh - Refresh: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, tokenPair)
}

// Update godoc
//
// @Summary      Update info about user
// @Tags         users
// @Produce      json
// @Param		 body	body	model.User	 true	"New data"
// @Success      201
// @Failure      400
// @Failure      500
// @Router       /update [put]
// @Security Bearer
func (u *User) Update(c echo.Context) (err error) {
	user := &model.User{}
	err = c.Bind(user)
	if err != nil {
		logrus.Error(fmt.Errorf("user - Update - Bind: %w", err))
		return err
	}

	if user.Name != "" {
		if err = u.val.Var(user.Name, "alpha,gt=2,lte=25"); err != nil {
			err = fmt.Errorf("user - Update - Var: %w", err)
		}
	}
	if err != nil {
		err = fmt.Errorf("user - Update - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	if user.Email != "" {
		if err = u.val.Var(user.Email, "email"); err != nil {
			err = fmt.Errorf("user - Update - Var: %w", err)
		}
	}
	if err != nil {
		err = fmt.Errorf("user - Update - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	if user.Age != 0 {
		if err = u.val.Var(user.Age, "gt=0,lte=100"); err != nil {
			err = fmt.Errorf("user - Update - Var: %w", err)
		}
	}
	if err != nil {
		err = fmt.Errorf("user - Update - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	id := tokenFromContext(c)
	err = u.userService.Update(c.Request().Context(), id, user)
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
// @Param        id	 header   string	true  "login"
// @Success      201	object	model.User
// @Failure      403
// @Failure      500
// @Router       /admin/userByLogin [get]
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

func tokenFromContext(c echo.Context) (ID string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims
	return claims.(*model.CustomClaims).ID
}
