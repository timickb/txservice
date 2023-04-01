package messaging

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/timickb/txserver/internal/interfaces"
)

type Consumer struct {
	reader *kafka.Reader
	logger interfaces.Logger
}

func NewConsumer(kafkaUrl, groupId string, topics []string, logger interfaces.Logger) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{kafkaUrl},
			GroupTopics: topics,
			GroupID:     groupId,
		}),
		logger: logger,
	}
}

func (c *Consumer) PullMessages(ctx context.Context, messages chan<- kafka.Message) error {
	c.logger.Info("Start kafka messages fetching loop")

	for {
		c.logger.Info("Waiting for message...")
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case messages <- msg:
			c.logger.Info("Kafka message fetched and pushed to channel: ", string(msg.Value))
		}
	}
}

func (c *Consumer) CommitMessages(ctx context.Context, commits <-chan kafka.Message) error {
	c.logger.Info("Start kafka messages committing loop")

	for {
		select {
		case msg := <-commits:
			err := c.reader.CommitMessages(ctx, msg)
			if err != nil {
				return err
			}
			c.logger.Info("Kafka message committed: ", string(msg.Value))
		case <-ctx.Done():
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
