package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_URL              string
	LogLevel                  string
	GRPCCartServiceAddress    string
	GRPCProductServiceAddress string
	HTTPPort                  string
}

func InitConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	return Config{
		POSTGRES_URL:              os.Getenv("POSTGRES_URL"),
		LogLevel:                  os.Getenv("LOG_LEVEL"),
		GRPCCartServiceAddress:    os.Getenv("GRPC_CART_SERVICE_ADDRESS"),
		GRPCProductServiceAddress: os.Getenv("GRPC_PRODUCT_SERVICE_ADDRESS"),
		HTTPPort:                  os.Getenv("HTTP_PORT"),
	}
}
