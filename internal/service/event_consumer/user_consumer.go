package eventconsumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/google/uuid"
)

type UserSaver interface {
	CreateUser(ctx context.Context, userID uuid.UUID, email string) (model.User, error)
}

type Consumer interface {
	RunConsume(ctx context.Context, handler func([]byte) error) error
	Stop() error
}

type UserEventConsumer struct {
	log       *slog.Logger
	userSaver UserSaver
	consumer  Consumer
}

func New(log *slog.Logger, userSaver UserSaver, consumer Consumer) *UserEventConsumer {
	return &UserEventConsumer{
		log:       log,
		userSaver: userSaver,
		consumer:  consumer,
	}
}

func (u *UserEventConsumer) MustRun(ctx context.Context) {
	if err := u.consumer.RunConsume(ctx, u.userHandler); err != nil {
		panic(err)
	}
}

func (u *UserEventConsumer) userHandler(data []byte) error {
	const op = "eventconsumer.UserEventConsumer.userHandler"
	log := u.log.With(slog.String("op", op))

	var userPayload model.UserEvent
	if err := json.Unmarshal(data, &userPayload); err != nil {
		log.Error("unmarshal failed", sl.Err(err))
		return err
	}

	if _, err := u.userSaver.CreateUser(context.TODO(), userPayload.ID, userPayload.Email); err != nil {
		log.Error("creating user failed", sl.Err(err))
		return err
	}

	return nil
}

func (u *UserEventConsumer) Stop() error {
	return u.consumer.Stop()
}
