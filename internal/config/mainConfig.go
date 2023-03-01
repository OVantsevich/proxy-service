// Package config main config
package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

// MainConfig with init data
type MainConfig struct {
	JwtKey string `env:"JWT_KEY,notEmpty" envDefault:"874967EC3EA3490F8F2EF6478B72A756"`
	Port   string `env:"PORT,notEmpty" envDefault:"99999"`

	PriceServicePort string `env:"PRICE_SERVICE_PORT,notEmpty" envDefault:"10000"`
	PriceServiceHost string `env:"PRICE_SERVICE_HOST,notEmpty" envDefault:"localhost"`

	PaymentServicePort string `env:"PAYMENT_SERVICE_PORT,notEmpty" envDefault:"30000"`
	PaymentServiceHost string `env:"PAYMENT_SERVICE_HOST,notEmpty" envDefault:"localhost"`

	UserServicePort string `env:"PAYMENT_SERVICE_PORT,notEmpty" envDefault:"20000"`
	UserServiceHost string `env:"PAYMENT_SERVICE_HOST,notEmpty" envDefault:"localhost"`

	TradingServicePort string `env:"PAYMENT_SERVICE_PORT,notEmpty" envDefault:"40000"`
	TradingServiceHost string `env:"PAYMENT_SERVICE_HOST,notEmpty" envDefault:"localhost"`
}

// NewMainConfig parsing config from environment
func NewMainConfig() (*MainConfig, error) {
	mainConfig := &MainConfig{}

	err := env.Parse(mainConfig)
	if err != nil {
		return nil, fmt.Errorf("config - NewMainConfig - Parse:%w", err)
	}

	return mainConfig, nil
}
