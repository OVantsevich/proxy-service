// Package handler account handler
package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/OVantsevich/proxy-service/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// AccountService service interface for account handler
//
//go:generate mockery --name=AccountService --case=underscore --output=./mocks
type AccountService interface {
	CreateAccount(ctx context.Context, userID string) (*model.Account, error)
	GetAccount(ctx context.Context, userID string) (*model.Account, error)
	IncreaseAmount(ctx context.Context, accountID string, amount float64) error
	DecreaseAmount(ctx context.Context, accountID string, amount float64) error
}

// Account handler
type Account struct {
	accountService AccountService
}

// NewAccountHandler new account handler
func NewAccountHandler(s AccountService) *Account {
	return &Account{accountService: s}
}

// CreateAccount godoc
//
// @Summary      creating account for user
// @Tags         accounts
// @Produce      json
// @Success      201	{object}	model.Account
// @Failure      500	{object}	echo.HTTPError
// @Router       /createAccount [post]
// @Security Bearer
func (a *Account) CreateAccount(c echo.Context) (err error) {
	id := idFromContext(c)

	account, err := a.accountService.CreateAccount(c.Request().Context(), id)
	if err != nil {
		logrus.Error(fmt.Errorf("account - CreateAccount - CreateAccount: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusCreated, account)
}

// GetUserAccount godoc
//
// @Summary      getting account by user id
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200	{object}	model.Account
// @Failure      500	{object}	echo.HTTPError
// @Router       /getUserAccount [get]
// @Security Bearer
func (a *Account) GetUserAccount(c echo.Context) error {
	id := idFromContext(c)

	account, err := a.accountService.GetAccount(c.Request().Context(), id)
	if err != nil {
		logrus.Error(fmt.Errorf("account - GetAccount - GetAccount: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, account)
}

// AmountRequest inc dec amount request
type AmountRequest struct {
	Amount    float64 `json:"amount" validate:"required,gte=0"`
	AccountID string  `json:"accountID" validate:"required"`
}

// IncreaseAmount godoc
//
// @Summary      increase account amount
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        amount	body 		AmountRequest  true  "Amount of operation"
// @Success      200
// @Failure      500	{object}	echo.HTTPError
// @Router       /increaseAmount [post]
// @Security Bearer
//
//nolint:dupl //just because
func (a *Account) IncreaseAmount(c echo.Context) (err error) {
	amount := &AmountRequest{}
	err = c.Bind(amount)
	if err != nil {
		logrus.Error(fmt.Errorf("account - IncreaseAmount - Bind: %w", err))
		return err
	}

	err = c.Validate(amount)
	if err != nil {
		err = fmt.Errorf("account - IncreaseAmount - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	err = a.accountService.IncreaseAmount(c.Request().Context(), amount.AccountID, amount.Amount)
	if err != nil {
		err = fmt.Errorf("account - IncreaseAmount - IncreaseAmount: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, "")
}

// DecreaseAmount godoc
//
// @Summary      decrease account amount
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        amount	body  		AmountRequest  true  "Amount of operation"
// @Success      200
// @Failure      500	{object}	echo.HTTPError
// @Router       /decreaseAmount [post]
// @Security Bearer
//
//nolint:dupl //just because
func (a *Account) DecreaseAmount(c echo.Context) (err error) {
	amount := &AmountRequest{}
	err = c.Bind(amount)
	if err != nil {
		logrus.Error(fmt.Errorf("account - DecreaseAmount - Bind: %w", err))
		return err
	}

	err = c.Validate(amount)
	if err != nil {
		err = fmt.Errorf("account - DecreaseAmount - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	err = a.accountService.DecreaseAmount(c.Request().Context(), amount.AccountID, amount.Amount)
	if err != nil {
		err = fmt.Errorf("account - DecreaseAmount - IncreaseAmount: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, "")
}
