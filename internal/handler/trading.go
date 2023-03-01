// Package handler trading handler
package handler

import (
	"context"
	"fmt"
	"github.com/OVantsevich/proxy-service/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

// TradingService service interface for trading handler
//
//go:generate mockery --name=TradingService --case=underscore --output=./mocks
type TradingService interface {
	OpenPosition(ctx context.Context, position *model.Position) (*model.Position, error)
	GetPositionByID(ctx context.Context, positionID string) (*model.Position, error)
	GetUserPositions(ctx context.Context, userID string) ([]*model.Position, error)
	SetStopLoss(ctx context.Context, positionID string, stopLoss float64) error
	SetTakeProfit(ctx context.Context, positionID string, takeProfit float64) error
	ClosePosition(ctx context.Context, positionID string) error
}

type GetPositionByIdRequest struct {
	ID string `json:"id" validate:"required"`
}

// Trading handler
type Trading struct {
	tradingService TradingService

	val *validator.Validate
}

// NewTradingHandler new trading handler
func NewTradingHandler(s TradingService) *Trading {
	return &Trading{tradingService: s, val: validator.New()}
}

// OpenPosition godoc
//
// @Summary      Open new position
// @Tags         trading
// @Accept       json
// @Produce      json
// @Param        body	body     	model.Position  true  "New position"
// @Success      201	{object}	model.Position
// @Failure      500
// @Router       /openPosition [post]
func (t *Trading) OpenPosition(c echo.Context) error {
	_, id := tokenFromContext(c)

	position := &model.Position{}
	err := c.Bind(position)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - OpenPosition - Bind: %w", err))
		return err
	}

	err = c.Validate(position)
	if err != nil {
		err = fmt.Errorf("trading - OpenPosition - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	position.User = id
	positionResponse, err := t.tradingService.OpenPosition(c.Request().Context(), position)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - OpenPosition - OpenPosition: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusCreated, positionResponse)
}

// GetPositionByID godoc
//
// @Summary      Open new position
// @Tags         trading
// @Accept       json
// @Produce      json
// @Param        id		head   	  	string  true  "Position ID"
// @Success      201	{object}	model.Position
// @Failure      500
// @Router       /getPosition [get]
func (t *Trading) GetPositionByID(c echo.Context) error {
	_, id := tokenFromContext(c)

	request := &GetPositionByIdRequest{}
	err := c.Bind(request)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - GetPositionByID - Bind: %w", err))
		return err
	}

	err = c.Validate(request)
	if err != nil {
		err = fmt.Errorf("trading - GetPositionByID - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	positionResponse, err := t.tradingService.GetPositionByID(c.Request().Context(), request.ID)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - GetPositionByID - OpenPosition: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusCreated, positionResponse)
}
