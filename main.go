// Package main
package main

import (
	"context"
	"fmt"
	"net/http"

	pasProto "github.com/OVantsevich/Payment-Service/proto"
	prsProto "github.com/OVantsevich/Price-Service/proto"
	tsProto "github.com/OVantsevich/Trading-Service/proto"
	usProto "github.com/OVantsevich/User-Service/proto"
	_ "github.com/OVantsevich/proxy-service/docs"
	"github.com/OVantsevich/proxy-service/internal/config"
	"github.com/OVantsevich/proxy-service/internal/handler"
	"github.com/OVantsevich/proxy-service/internal/model"
	"github.com/OVantsevich/proxy-service/internal/repository"
	"github.com/OVantsevich/proxy-service/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CustomValidator echo validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate echo method
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

// @title Swagger Trading service API
// @version 1.0
// @description trading service.

// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	cfg, err := config.NewMainConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	connUser, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort), opts...)
	if err != nil {
		logrus.Fatal("Fatal Dial: ", err)
	}
	usClient := usProto.NewUserServiceClient(connUser)
	userRepository := repository.NewUserServiceRepository(usClient)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, cfg.JwtKey)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	noAuthentication := e.Group("/auth")
	noAuthentication.POST("/signup", userHandler.Signup)
	noAuthentication.POST("/login", userHandler.Login)
	noAuthentication.GET("/refresh", userHandler.Refresh)

	withAuthentication := e.Group("")

	withAuthentication.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/swagger/*"
		},
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtKey), nil
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return model.CustomClaims{}
		},
	}))

	withAuthentication.POST("/update", userHandler.Update)
	withAuthentication.POST("/userByID", userHandler.UserByID)

	connPayment, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.PaymentServiceHost, cfg.PaymentServicePort), opts...)
	if err != nil {
		logrus.Fatal("Fatal Dial: ", err)
	}
	psClient := pasProto.NewPaymentServiceClient(connPayment)
	accountRepository := repository.NewPaymentServiceRepository(psClient)
	accountService := service.NewAccountService(accountRepository)
	accountHandler := handler.NewAccountHandler(accountService)

	withAuthentication.GET("/createAccount", accountHandler.CreateAccount)
	withAuthentication.GET("/getUserAccount", accountHandler.GetUserAccount)
	withAuthentication.POST("/increaseAmount", accountHandler.IncreaseAmount)
	withAuthentication.POST("/decreaseAmount", accountHandler.DecreaseAmount)

	connPrice, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.PriceServiceHost, cfg.PriceServicePort), opts...)
	if err != nil {
		logrus.Fatal("Fatal Dial: ", err)
	}
	prsClient := prsProto.NewPriceServiceClient(connPrice)
	priceRepository, err := repository.NewPriceServiceRepository(context.Background(), prsClient)
	if err != nil {
		logrus.Fatal(err)
	}
	priceService := service.NewPriceService(context.Background(), priceRepository, repository.NewListenersRepository())
	priceHandler := handler.NewPriceHandler(priceService)

	withAuthentication.GET("/getCurrentPrices", priceHandler.GetCurrentPrices)
	withAuthentication.GET("/subscribe", priceHandler.Subscribe)

	connTrading, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.TradingServiceHost, cfg.TradingServicePort), opts...)
	if err != nil {
		logrus.Fatal("Fatal Dial: ", err)
	}
	tsClient := tsProto.NewTradingServiceClient(connTrading)
	tradingRepository, err := repository.NewTradingServiceRepository(tsClient)
	if err != nil {
		logrus.Fatal(err)
	}
	tradingService := service.NewTradingService(tradingRepository)
	tradingHandler := handler.NewTradingHandler(tradingService)

	withAuthentication.POST("/openPosition", tradingHandler.OpenPosition)
	withAuthentication.GET("/getUserPositions", tradingHandler.GetUserPositions)
	withAuthentication.GET("/getPositionByID", tradingHandler.GetPositionByID)
	withAuthentication.POST("/setTakeProfit", tradingHandler.SetTakeProfit)
	withAuthentication.POST("/setStopLoss", tradingHandler.SetStopLoss)
	withAuthentication.POST("/closePosition", tradingHandler.ClosePosition)

	logrus.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
