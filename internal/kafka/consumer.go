package kafka

import (
	"context"
	"log/slog"

	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	log    *slog.Logger
}

func NewConsumer(log *slog.Logger, brokers []string, groupID string, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
		log: log,
	}
}

// RunConsume consumes messages from the Kafka topic and processes them using the handler function
func (c *Consumer) RunConsume(ctx context.Context, handler func([]byte) error) error {
	const op = "kafka.Consumer.Run"
	log := c.log.With(slog.String("op", op), slog.String("topic", c.reader.Config().Topic), slog.Any("brokers", c.reader.Config().Brokers))

	log.Info("starting consume messages")

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Error("error reading message", sl.Err(err))
			continue
		}

		log.Info("received", slog.String("message", string(msg.Value)))

		// Call the handler function with the message payload
		if err := handler(msg.Value); err != nil {
			log.Error("handler failed", sl.Err(err))
		}
	}
}

func (c *Consumer) Stop() error {
	return c.reader.Close()
}
