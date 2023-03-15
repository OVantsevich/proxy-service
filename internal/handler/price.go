// Package handler account handler
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OVantsevich/proxy-service/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// PriceService service interface for price service
//
//go:generate mockery --name=PriceService --case=underscore --output=./mocks
type PriceService interface {
	GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error)

	GetPrices() ([]*model.Price, error)
	Subscribe(streamID uuid.UUID) chan *model.Price
	UpdateSubscription(socketID uuid.UUID, names []string) error
	DeleteSubscription(streamID uuid.UUID) error
}

// PriceRequest websocket request
type PriceRequest struct {
	Names []string `json:"names" validate:"required,dive,alpha,gte=2,lte=25" example:"gold,google,tesla,oil"`
}

// PriceResponse websocket response
type PriceResponse struct {
	Prices []*model.Price `json:"prices"`
}

// Price handler
type Price struct {
	priceService PriceService

	val *validator.Validate
}

// NewPriceHandler new price handler
func NewPriceHandler(s PriceService) *Price {
	return &Price{priceService: s, val: validator.New()}
}

// Subscribe godoc
//
// @Summary      Subscribe for prices
// @Tags         prices
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       /subscribe [get]
// @Security Bearer
func (p *Price) Subscribe(c echo.Context) error {
	h := websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		socketID := uuid.New()
		priceChan := p.priceService.Subscribe(socketID)
		defer p.priceService.DeleteSubscription(socketID)

		out := make(chan *PriceRequest)
		go getPrice(ws, out)
		go sendPrice(ws, priceChan)

		for {
			priceRequest, ok := <-out
			if !ok {
				return
			}
			err := p.priceService.UpdateSubscription(socketID, priceRequest.Names)
			if err != nil {
				logrus.Errorf("price - Subscribe - UpdateSubscription: %v", err)
				return
			}
		}
	})
	s := websocket.Server{Handler: h, Handshake: nil}
	s.ServeHTTP(c.Response(), c.Request())
	return nil
}

func sendPrice(ws *websocket.Conn, in chan *model.Price) {
	for {
		data, ok := <-in
		if !ok {
			return
		}

		marshalData, err := json.Marshal(data)
		if err != nil {
			logrus.Errorf("price - Subscribe - sendPrice - Marshal: %v", err)
			return
		}

		err = websocket.Message.Send(ws, marshalData)
		if err != nil {
			logrus.Errorf("price - Subscribe - sendPriceResponse - Send: %v", err)
			return
		}
	}
}

func getPrice(ws *websocket.Conn, out chan *PriceRequest) {
	for {
		priceRequest := &PriceRequest{}
		var data []byte
		err := websocket.Message.Receive(ws, &data)
		if err != nil {
			close(out)
			logrus.Errorf("price - Subscribe - getPriceRequest - Receive: %v", err)
			return
		}

		err = json.Unmarshal(data, priceRequest)
		if err != nil {
			close(out)
			logrus.Errorf("price - Subscribe - getPriceRequest - Unmarshal: %v", err)
			return
		}
		out <- priceRequest
	}
}

// GetCurrentPriceResponse gcp response
type GetCurrentPriceResponse struct {
	Prices map[string]*model.Price `json:"prices"`
}

// GetCurrentPrices godoc
//
// @Summary      get current prices
// @Tags         prices
// @Accept       json
// @Produce      json
// @Param        names	body 		PriceRequest  true  "Prices list"
// @Success      200   	object		GetCurrentPriceResponse
// @Failure      500	{object}	echo.HTTPError
// @Router       /getCurrentPrices	[post]
// @Security Bearer
func (p *Price) GetCurrentPrices(c echo.Context) (err error) {
	names := &PriceRequest{}
	err = c.Bind(names)
	if err != nil {
		logrus.Error(fmt.Errorf("price - GetCurrentPrices - Bind: %w", err))
		return err
	}

	err = c.Validate(names)
	if err != nil {
		err = fmt.Errorf("price - GetCurrentPrices - Validate: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	prices, err := p.priceService.GetCurrentPrices(c.Request().Context(), names.Names)
	if err != nil {
		err = fmt.Errorf("price - GetCurrentPrices - GetCurrentPrices: %w", err)
		logrus.Error(err)
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, GetCurrentPriceResponse{Prices: prices})
}
