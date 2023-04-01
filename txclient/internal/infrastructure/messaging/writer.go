package messaging

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/timickb/txclient/internal/interfaces"
	"time"
)

type Producer struct {
	writer *kafka.Writer
	logger interfaces.Logger
}

func NewProducer(addr, topic string, logger interfaces.Logger) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(addr),
			Topic:        topic,
			WriteTimeout: time.Second * 15,
		},
		logger: logger,
	}
}

func (p *Producer) PushMessage(ctx context.Context, msg kafka.Message) error {
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (p *Producer) PushMessages(ctx context.Context, msgChan chan kafka.Message) error {
	p.logger.Info("Start kafka messages writing loop")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-msgChan:
			err := p.writer.WriteMessages(ctx, kafka.Message{
				Value: msg.Value,
			})
			if err != nil {
				return err
			}
		}
	}
}

func (p *Producer) Close() error {
	p.logger.Info("Closed kafka writer for topic ", p.writer.Topic)
	return p.writer.Close()
}
