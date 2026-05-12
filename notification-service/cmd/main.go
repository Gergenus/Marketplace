package main

import (
	"github.com/Gergenus/commerce/notification-service/internal/brocker"
	"github.com/Gergenus/commerce/notification-service/internal/config"
	"github.com/Gergenus/commerce/notification-service/internal/email"
	"github.com/Gergenus/commerce/notification-service/internal/service"
	"github.com/Gergenus/commerce/notification-service/pkg/jwtpkg"
	"github.com/Gergenus/commerce/notification-service/pkg/kafka"
	"github.com/Gergenus/commerce/notification-service/pkg/logger"
)

func main() {
	cfg := config.InitConfig()
	log := logger.SetUp(cfg.LogLevel)
	mail := email.NewEmailSender(cfg.FromEmail, cfg.FromEmailPassword, cfg.FromEmailSMTP, cfg.SMTPAddr, log)
	log.Info("starting the project")
	jwtPkg := jwtpkg.NewJWTpkg(cfg.JWTMailSecret, cfg.TokenTTL)

	kafkaConsumer := kafka.ConnectConsumer([]string{cfg.KafkaURL})
	kafkaBrocker := brocker.NewKafkaBrocker(kafkaConsumer)
	servce := service.NewNotificationService(log, &mail, kafkaBrocker, jwtPkg)

	if err := servce.ListenToTopic("mail"); err != nil {
		panic(err)
	}

	log.Info("stopping the project")
}
