package dummy

import (
	"errors"
	"github.com/timickb/txserver/internal/domain"
)

type Repository struct {
	// There is no sync.Map or mutexes cause this repo
	// created just for unit tests
	accounts     map[string]*domain.Account
	transactions map[string]*domain.Transaction
}

func New() *Repository {
	return &Repository{
		accounts:     make(map[string]*domain.Account),
		transactions: make(map[string]*domain.Transaction),
	}
}

func (r *Repository) CreateAccount(acc domain.Account) error {
	r.accounts[acc.Id] = &acc
	return nil
}

func (r *Repository) UpdateAccount(acc domain.Account) error {
	if _, ok := r.accounts[acc.Id]; !ok {
		return errors.New("err not found")
	}

	r.accounts[acc.Id] = &acc
	return nil
}

func (r *Repository) FindAccount(id string) (*domain.Account, error) {
	value, ok := r.accounts[id]
	if !ok {
		return nil, errors.New("err not found")
	}
	return value, nil
}

func (r *Repository) PerformTransaction(t domain.Transaction) error {
	sender, sndOk := r.accounts[t.SenderId]
	_, rcvOk := r.accounts[t.ReceiverId]

	if !sndOk {
		return errors.New("err sender not found")
	}
	if !rcvOk {
		return errors.New("err receiver not found")
	}

	if sender.Balance < t.Amount {
		return errors.New("err not enough money")
	}

	r.accounts[t.SenderId].Balance -= t.Amount
	r.accounts[t.ReceiverId].Balance += t.Amount
	r.transactions[t.Id] = &t

	return nil
}
