package interfaces

import (
	"context"
	"github.com/segmentio/kafka-go"
	"io"
)

type Consumer interface {
	io.Closer
	PullMessages(ctx context.Context, messages chan<- kafka.Message) error
	CommitMessages(ctx context.Context, commits <-chan kafka.Message) error
}
type Producer interface {
	io.Closer
	PushMessages(ctx context.Context, msgChan chan kafka.Message, commitChan chan kafka.Message) error
}
