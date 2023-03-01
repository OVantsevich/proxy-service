// Package handler account handler
package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/OVantsevich/proxy-service/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// AccountService service interface for account handler
//
//go:generate mockery --name=AccountService --case=underscore --output=./mocks
type AccountService interface {
	GetAccount(ctx context.Context, userID string) (*model.Account, error)
	IncreaseAmount(ctx context.Context, accountID string, amount float64) error
	DecreaseAmount(ctx context.Context, accountID string, amount float64) error
}

type AmountRequest struct {
	Amount float64 `json:"amount" validate:"required,gte=0"`
}

// Account handler
type Account struct {
	accountService AccountService

	val *validator.Validate
}

// NewAccountHandler new account handler
func NewAccountHandler(s AccountService) *Account {
	return &Account{accountService: s, val: validator.New()}
}

// GetAccount godoc
//
// @Summary      Get user account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200	{object}	model.Account
// @Failure      500
// @Router       /getAccount [get]
func (a *Account) GetAccount(c echo.Context) error {
	_, id := tokenFromContext(c)

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

// IncreaseAmount godoc
//
// @Summary      increase account amount
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body	body  AmountRequest  true  "Amount of operation"
// @Success      200
// @Failure      500
// @Router       /increaseAmount [post]
func (a *Account) IncreaseAmount(c echo.Context) (err error) {
	_, id := tokenFromContext(c)
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

	err = a.accountService.IncreaseAmount(c.Request().Context(), id, amount.Amount)
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
// @Param        body	body  AmountRequest  true  "Amount of operation"
// @Success      200
// @Failure      500
// @Router       /decreaseAmount [post]
func (a *Account) DecreaseAmount(c echo.Context) (err error) {
	_, id := tokenFromContext(c)
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

	err = a.accountService.DecreaseAmount(c.Request().Context(), id, amount.Amount)
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
