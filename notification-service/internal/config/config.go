package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel          string
	JWTMailSecret     string
	TokenTTL          time.Duration
	FromEmail         string
	FromEmailPassword string
	FromEmailSMTP     string
	SMTPAddr          string
	KafkaURL          string
}

func InitConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	TokenTTL, err := time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		panic(err)
	}

	return Config{
		LogLevel:          os.Getenv("LOG_LEVEL"),
		JWTMailSecret:     os.Getenv("JWT_MAIL_SECRET"),
		TokenTTL:          TokenTTL,
		FromEmail:         os.Getenv("FROM_EMAIL"),
		FromEmailPassword: os.Getenv("FROM_EMAIL_PASSWORD"),
		FromEmailSMTP:     os.Getenv("FROM_EMAIL_SMTP"),
		SMTPAddr:          os.Getenv("SMTP_ADDR"),
		KafkaURL:          os.Getenv("KAFKA_URL"),
	}
}
