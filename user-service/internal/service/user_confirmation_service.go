package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Gergenus/commerce/user-service/internal/models"
)

func (u *UserService) ConfirmUser(ctx context.Context, uuid string) error {
	const op = "service.ConfirmUser"
	log := u.log.With(slog.String("op", op))
	user, err := u.repo.GetUserByUUID(ctx, uuid)
	if err != nil {
		log.Error("getting user by uuid error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	kafkaRequest := models.ConfirmationRequest{Email: user.Email}

	req, err := json.Marshal(kafkaRequest)
	if err != nil {
		log.Error("request marshaling error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	part, offset, err := u.EventProducer.SendMailOrder("mail", req)
	if err != nil {
		log.Error("sending request to kafka error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("data", slog.Int("partition", int(part)), slog.Int("offset", int(offset)))
	return nil
}
