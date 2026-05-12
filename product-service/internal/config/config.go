package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL              string
	LogLevel                 string
	HTTPPort                 string
	JWTSecret                string
	GRPCProductServerAddress string
	ElasticAddress           string
	ElasticUser              string
	ElasticPassword          string
	ElasticCrt               string
}

func InitConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	return Config{
		PostgresURL:              os.Getenv("PostgresURL"),
		LogLevel:                 os.Getenv("LogLevel"),
		HTTPPort:                 os.Getenv("HTTPPort"),
		JWTSecret:                os.Getenv("JWTSecret"),
		GRPCProductServerAddress: os.Getenv("GRPC_PRODUCT_SERVER_ADDRESS"),
		ElasticAddress:           os.Getenv("ELASTIC_ADDRESS"),
		ElasticUser:              os.Getenv("ELASTIC_USER"),
		ElasticPassword:          os.Getenv("ELASTIC_PASSWORD"),
		ElasticCrt:               os.Getenv("ELASTIC_CRT"),
	}
}
