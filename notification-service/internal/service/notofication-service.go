package service

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gergenus/commerce/notification-service/internal/brocker"
	"github.com/Gergenus/commerce/notification-service/internal/email"
	"github.com/Gergenus/commerce/notification-service/internal/models"
	"github.com/Gergenus/commerce/notification-service/pkg/jwtpkg"
	"github.com/IBM/sarama"
)

type NotificationService struct {
	log     *slog.Logger
	mail    email.EmailSenderInterface
	brock   brocker.BrockerInterface
	jwToken jwtpkg.JWTpkgInterface
}

func NewNotificationService(log *slog.Logger, mail email.EmailSenderInterface, brock brocker.BrockerInterface, jwToken jwtpkg.JWTpkgInterface) NotificationService {
	return NotificationService{log: log, mail: mail, brock: brock, jwToken: jwToken}
}

func (n NotificationService) ListenToTopic(topic string) error {
	const op = "service.ListenToTopic"
	log := n.log.With(slog.String("op", op), slog.String("topic", topic))
	consumer, err := n.brock.RecieveMessages(topic)
	if err != nil {
		log.Error("recieving error", slog.String("error", err.Error()))
		return err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error("consumer error", slog.String("error", err.Error()))
			case msg := <-consumer.Messages():
				n.sendVerificationEmail(msg)
			case <-sigchan:
				done <- struct{}{}
			}
		}

	}()

	<-done
	err = consumer.Close()
	if err != nil {
		log.Error("closing consumer error", slog.String("error", err.Error()))
	}
	return nil
}

func (n NotificationService) sendVerificationEmail(msg *sarama.ConsumerMessage) error {
	const op = "service.ListenToTopic.sendEmail"
	log := n.log.With(slog.String("op", op), slog.String("topic", msg.Topic))
	var data models.ConfirmationResponse

	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		log.Error("unmarshaling msg error", slog.String("error", err.Error()))
		return err
	}
	token, err := n.jwToken.GenerateToken(data.Email)
	if err != nil {
		log.Error("generating token error", slog.String("error", err.Error()))
		return err
	}

	err = n.mail.SendVerificationEmail(data.Email, "profile confirmation", "verification_email", map[string]string{"VerificationLink": "http://localhost:8081/api/v1/users/auth/verification?token=" + token})
	if err != nil {
		log.Error("email sending error", slog.String("error", err.Error()))
		return err
	}
	return nil
}
