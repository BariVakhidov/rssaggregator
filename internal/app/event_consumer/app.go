package kafkaconsumerapp

import (
	"context"
	"log/slog"

	"github.com/BariVakhidov/rssaggregator/internal/kafka"
	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	eventconsumer "github.com/BariVakhidov/rssaggregator/internal/service/event_consumer"
)

type App struct {
	userConsumer *eventconsumer.UserEventConsumer
	log          *slog.Logger
}

func New(log *slog.Logger,
	userSaver eventconsumer.UserSaver,
	brokers []string,
	topic string,
) *App {
	consumer := kafka.NewConsumer(log, brokers, "user-consumer-group", topic)
	userConsumer := eventconsumer.New(log, userSaver, consumer)

	return &App{userConsumer: userConsumer, log: log}
}

func (a *App) MustRun(ctx context.Context) {
	a.userConsumer.MustRun(ctx)
}

func (a *App) Stop() {
	const op = "app.eventconsumerapp.Stop"
	log := a.log.With(slog.String("op", op))

	log.Info("stopping user consumer...")

	if err := a.userConsumer.Stop(); err != nil {
		log.Error("user consumer failed to stop", sl.Err(err))
	}
}
