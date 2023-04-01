package interfaces

import (
	"context"
	"github.com/segmentio/kafka-go"
	"io"
)

type Producer interface {
	io.Closer
	PushMessage(ctx context.Context, msg kafka.Message) error
	PushMessages(ctx context.Context, msgChan chan kafka.Message) error
}
