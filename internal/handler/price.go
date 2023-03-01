// Package handler account handler
package handler

import (
	"context"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"

	"github.com/OVantsevich/proxy-service/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// PriceService service interface for price service
//
//go:generate mockery --name=PriceService --case=underscore --output=./mocks
type PriceService interface {
	GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error)

	GetPrices() ([]*model.Price, error)
	UpdateSubscription(names []string) error
}

type PriceRequest struct {
	Prices []string `json:"prices" validate:"required,dive,alpha,gte=2,lte=25"`
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
func (p *Price) Subscribe(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		in := make(chan interface{})
		out := make(chan *PriceRequest)

		go sendPriceResponse(ws, in)
		go getPriceRequest(ws, out)

		for {
			select {}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func sendPriceResponse(ws *websocket.Conn, in chan interface{}) {
	for {
		select {
		case data, ok := <-in:
			if !ok {
				return
			}
			err := websocket.Message.Send(ws, data)
			if err != nil {
				logrus.Errorf("price - Subscribe - sendPriceResponse - Send: %v", err)
				return
			}
		}
	}
}

func getPriceRequest(ws *websocket.Conn, out chan *PriceRequest) {
	for {
		select {
		default:
			var priceRequest *PriceRequest
			err := websocket.Message.Receive(ws, priceRequest)
			if err != nil {
				logrus.Errorf("price - Subscribe - getPriceRequest - Receive: %v", err)
				return
			}
			out <- priceRequest
		}
	}
}

// GetCurrentPrices godoc
//
// @Summary      decrease account amount
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body	body  PriceRequest  true  "Prices list"
// @Success      200
// @Failure      500
// @Router       /decreaseAmount [post]
func (p *Price) GetCurrentPrices(c echo.Context) (err error) {
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

	err = p.priceService.GetCurrentPrices(c.Request().Context(), id, amount.Amount)
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
