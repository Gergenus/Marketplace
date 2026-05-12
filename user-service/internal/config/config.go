package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL   string
	LogLevel      string
	HTTPPort      string
	KafkaURL      string
	JWTSecret     string
	JWTMailSecret string
	AccessTTL     time.Duration
	RefreshTTl    time.Duration
	EmailTTL      time.Duration
}

func InitConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	AccessTTL, err := time.ParseDuration(os.Getenv("ACCESS_TTL"))
	if err != nil {
		panic(err)
	}
	RefreshTTl, err := time.ParseDuration(os.Getenv("REFRESH_TTL"))
	if err != nil {
		panic(err)
	}
	EmailTTL, err := time.ParseDuration(os.Getenv("EMAIL_TTL"))
	if err != nil {
		panic(err)
	}
	return Config{
		PostgresURL:   os.Getenv("PostgresURL"),
		LogLevel:      os.Getenv("LogLevel"),
		HTTPPort:      os.Getenv("HTTPPort"),
		KafkaURL:      os.Getenv("KAFKA_URL"),
		JWTSecret:     os.Getenv("JWTSecret"),
		AccessTTL:     AccessTTL,
		RefreshTTl:    RefreshTTl,
		EmailTTL:      EmailTTL,
		JWTMailSecret: os.Getenv("JWT_MAIL_SECRET"),
	}
}
