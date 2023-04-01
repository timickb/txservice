package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/timickb/txserver/internal/domain"
	"github.com/timickb/txserver/internal/interfaces"
	"sync"
)

type Storage interface {
	UpdateAccount(acc domain.Account) error
	PerformTransaction(t domain.Transaction) error
	FindAccount(id string) (*domain.Account, error)
}

type TransactionUseCase struct {
	store Storage
	log   interfaces.Logger
	mu    sync.Mutex // Several goroutines can use one usecase instance
}

func NewTransaction(store Storage, log interfaces.Logger) *TransactionUseCase {
	return &TransactionUseCase{store: store, log: log}
}

func (u *TransactionUseCase) PerformTransaction(sndId, rcvId string, amount int64) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if err := u.primaryTxValid(sndId, rcvId, amount); err != nil {
		return "", err
	}

	tx := domain.Transaction{
		Id:         uuid.NewString(),
		SenderId:   sndId,
		ReceiverId: rcvId,
		Amount:     amount,
	}

	if err := u.store.PerformTransaction(tx); err != nil {
		return "", err
	}

	return tx.Id, nil
}

func (u *TransactionUseCase) primaryTxValid(sndId, rcvId string, amount int64) error {
	if amount <= 0 {
		return errors.New("err non-positive amount")
	}
	if sndId == rcvId {
		return errors.New("err sender equaled to receiver")
	}

	sender, err := u.store.FindAccount(sndId)
	if err != nil {
		return err
	}

	receiver, err := u.store.FindAccount(rcvId)
	if err != nil {
		return err
	}

	u.log.Info("sender: ", sender)
	u.log.Info("receiver: ", receiver)

	return nil
}
