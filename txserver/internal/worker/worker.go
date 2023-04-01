package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/timickb/txserver/internal/domain"
	"github.com/timickb/txserver/internal/interfaces"
)

type TransactionUseCase interface {
	PerformTransaction(sndId, rcvId string, amount int64) (string, error)
}

type ReadWorker struct {
	log  interfaces.Logger
	cons interfaces.Consumer
	tuc  TransactionUseCase

	msgChan       chan kafka.Message
	msgCommitChan chan kafka.Message

	ctx    context.Context
	cancel context.CancelFunc
}

func New(cons interfaces.Consumer, tuc TransactionUseCase, log interfaces.Logger) *ReadWorker {
	return &ReadWorker{
		cons:          cons,
		log:           log,
		tuc:           tuc,
		msgChan:       make(chan kafka.Message),
		msgCommitChan: make(chan kafka.Message),
	}
}

func (w *ReadWorker) Start(ctx context.Context, errChan chan error) error {
	w.ctx, w.cancel = context.WithCancel(ctx)

	// Receiving kafka messages goroutine
	go func() {
		if err := w.cons.PullMessages(w.ctx, w.msgChan); err != nil {
			errChan <- err
		}
	}()

	// Committing goroutine. Commit means mark as received.
	go func() {
		if err := w.cons.CommitMessages(w.ctx, w.msgCommitChan); err != nil {
			errChan <- err
		}
	}()

	// Main worker loop
	for {
		select {
		case msg := <-w.msgChan:
			if err := w.handleMessage(msg); err != nil {
				errChan <- err
			}
			w.msgCommitChan <- msg
		case err := <-errChan:
			w.log.Error("Requests worker: %s", err.Error())
		case <-w.ctx.Done():
			return w.ctx.Err()
		}
	}
}

func (w *ReadWorker) Stop() error {
	w.cancel()

	close(w.msgChan)
	close(w.msgCommitChan)

	return w.cons.Close()
}

func (w *ReadWorker) handleMessage(msg kafka.Message) error {
	var tx domain.Transaction

	if err := json.Unmarshal(msg.Value, &tx); err != nil {
		return fmt.Errorf("err parse kafka msg: %w", err)
	}

	w.log.Info(fmt.Sprintf("Parsed request from topic %s: %v", msg.Topic, tx))

	w.log.Info(fmt.Sprintf(
		"start transaction performing with sender=%s, receiver=%s, amount=%d",
		tx.SenderId, tx.ReceiverId, tx.Amount))

	txId, err := w.tuc.PerformTransaction(
		tx.SenderId,
		tx.ReceiverId,
		tx.Amount)

	if err != nil {
		w.log.Error("transaction failed: ", err.Error())
		return err
	}

	w.log.Info("transaction succeed, id = ", txId)

	return nil
}
