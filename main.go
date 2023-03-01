// Package main
package main

import (
	"fmt"
	"github.com/OVantsevich/proxy-service/internal/handler"
	"github.com/OVantsevich/proxy-service/internal/service"
	"google.golang.org/grpc/credentials/insecure"

	pasProto "github.com/OVantsevich/Payment-Service/proto"
	usProto "github.com/OVantsevich/User-Service/proto"
	"github.com/OVantsevich/proxy-service/internal/config"
	"github.com/OVantsevich/proxy-service/internal/model"
	"github.com/OVantsevich/proxy-service/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// @title Trading service API
// @version 1.0
// @description trading service.

// @host localhost:99999
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	e := echo.New()

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
	userHandler := handler.NewUserHandler(userService)

	noAuthentication := e.Group("/auth")
	noAuthentication.POST("/signup", userHandler.Signup)
	noAuthentication.GET("/login", userHandler.Login)
	noAuthentication.GET("/refresh", userHandler.Refresh)

	withAuthentication := e.Group("")
	withAuthentication.Use(echojwt.WithConfig(echojwt.Config{
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtKey), nil
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return model.CustomClaims{}
		},
	}))

	withAuthentication.POST("/update", userHandler.Update)
	withAuthentication.POST("/UserByID", userHandler.UserByID)

	connPayment, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.PaymentServiceHost, cfg.PaymentServicePort), opts...)
	if err != nil {
		logrus.Fatal("Fatal Dial: ", err)
	}
	psClient := pasProto.NewPaymentServiceClient(connPayment)
	accountRepository := repository.NewPaymentServiceRepository(psClient)
	accountService := service.NewAccountService(accountRepository)
	accountHandler := handler.NewAccountHandler(accountService)

	withAuthentication.GET("/getAccount", accountHandler.GetAccount)
	withAuthentication.POST("/increaseAmount", accountHandler.IncreaseAmount)
	withAuthentication.POST("/decreaseAmount", accountHandler.DecreaseAmount)
}
