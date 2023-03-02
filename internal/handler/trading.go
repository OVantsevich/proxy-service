// Package handler trading handler
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

// GetPositionByIDRequest id request
type GetPositionByIDRequest struct {
	ID string `json:"id" validate:"required"`
}

// SetThresholdRequest request for SL and TP set
type SetThresholdRequest struct {
	ID     string  `json:"id" validate:"required"`
	Amount float64 `json:"amount" validate:"required,gte=0"`
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
	id := tokenFromContext(c)

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
// @Summary      Get position by ID
// @Tags         trading
// @Accept       json
// @Produce      json
// @Param        id		header   	  	string  true  "Position ID"
// @Success      200	{object}	model.Position
// @Failure      500
// @Router       /getPositionByID [get]
//
//nolint:dupl //just because
func (t *Trading) GetPositionByID(c echo.Context) error {
	request := &GetPositionByIDRequest{}
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
		logrus.Error(fmt.Errorf("trading - GetPositionByID - GetPositionByID: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, positionResponse)
}

// GetUserPositions godoc
//
// @Summary      Get al user positions
// @Tags         trading
// @Accept       json
// @Produce      json
// @Param        id		header   	  	string  true  "User ID"
// @Success      200	{array}	model.Position
// @Failure      500
// @Router       /getUserPosition [get]
//
//nolint:dupl //just because
func (t *Trading) GetUserPositions(c echo.Context) error {
	request := &GetPositionByIDRequest{}
	err := c.Bind(request)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - GetUserPositions - Bind: %w", err))
		return err
	}

	err = c.Validate(request)
	if err != nil {
		err = fmt.Errorf("trading - GetUserPositions - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	positionResponse, err := t.tradingService.GetUserPositions(c.Request().Context(), request.ID)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - GetUserPositions - GetUserPositions: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, positionResponse)
}

// SetStopLoss godoc
//
// @Summary      Set stop loss for position
// @Tags         trading
// @Accept       json
// @Produce      json
// @Param        id		body   	  	SetThresholdRequest  true  "ID and threshold of position"
// @Success      200
// @Failure      500
// @Router       /setStopLoss [post]
//
//nolint:dupl //just because
func (t *Trading) SetStopLoss(c echo.Context) error {
	request := &SetThresholdRequest{}
	err := c.Bind(request)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - SetStopLoss - Bind: %w", err))
		return err
	}

	err = c.Validate(request)
	if err != nil {
		err = fmt.Errorf("trading - SetStopLoss - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	err = t.tradingService.SetStopLoss(c.Request().Context(), request.ID, request.Amount)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - SetStopLoss - SetStopLoss: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, "")
}

// SetTakeProfit godoc
//
// @Summary      Set take profit for position
// @Tags         trading
// @Accept       json
// @Produce      json
// @Param        id		body   	  	SetThresholdRequest  true  "ID and threshold of position"
// @Success      201	{object}	model.Position
// @Failure      500
// @Router       /setTakeProfit [post]
//
//nolint:dupl //just because
func (t *Trading) SetTakeProfit(c echo.Context) error {
	request := &SetThresholdRequest{}
	err := c.Bind(request)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - SetTakeProfit - Bind: %w", err))
		return err
	}

	err = c.Validate(request)
	if err != nil {
		err = fmt.Errorf("trading - SetTakeProfit - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	err = t.tradingService.SetTakeProfit(c.Request().Context(), request.ID, request.Amount)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - SetTakeProfit - SetTakeProfit: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, "")
}

// ClosePosition godoc
//
// @Summary      Close position
// @Tags         trading
// @Accept       json
// @Produce      json
// @Param        id		header   	  	string  true  "Position ID"
// @Success      200
// @Failure      500
// @Router       /closePosition [post]
func (t *Trading) ClosePosition(c echo.Context) error {
	var request string
	err := c.Bind(request)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - ClosePosition - Bind: %w", err))
		return err
	}

	err = t.tradingService.ClosePosition(c.Request().Context(), request)
	if err != nil {
		logrus.Error(fmt.Errorf("trading - ClosePosition - ClosePosition: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, "")
}
