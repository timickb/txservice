package service

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/timickb/txclient/internal/domain"
	"github.com/timickb/txclient/internal/interfaces"
)

type SendService struct {
	prod interfaces.Producer
	log  interfaces.Logger
}

func NewSender(prod interfaces.Producer, log interfaces.Logger) *SendService {
	return &SendService{
		prod: prod,
		log:  log,
	}
}

func (w *SendService) SendRequest(ctx context.Context, r domain.ClientRequest) error {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err
	}

	w.log.Info("Sending message to kafka")

	err = w.prod.PushMessage(ctx, kafka.Message{
		Value: bytes,
	})
	if err != nil {
		return err
	}

	return nil
}

func (w *SendService) Close() error {
	return w.prod.Close()
}
