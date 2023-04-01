package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/timickb/txserver/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAccount(acc domain.Account) error {
	st := `INSERT INTO accounts (id, balance) VALUES ($1, $2)`

	_, err := r.db.Exec(st, acc.Id, acc.Balance)
	if err != nil {
		return fmt.Errorf("db insertion err: %w", err)
	}

	return nil
}

func (r *Repository) UpdateAccount(acc domain.Account) error {
	st := `UPDATE accounts SET balance=$1 WHERE id=$2`

	result, err := r.db.Exec(st, acc.Balance, acc.Id)
	if err != nil {
		return fmt.Errorf("db update err: %w", err)
	}

	if affected, _ := result.RowsAffected(); affected != 1 {
		return errors.New("err invalid affected rows count")
	}

	return nil
}

func (r *Repository) FindAccount(id string) (*domain.Account, error) {
	row := r.db.QueryRow(`SELECT * FROM accounts WHERE id=$1`, id)

	acc := domain.Account{}

	if err := row.Scan(&acc.Id, &acc.Balance); err != nil {
		return nil, err
	}

	return &acc, nil
}

func (r *Repository) PerformTransaction(t domain.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var sndBal int64
	var rcvBal int64

	row := tx.QueryRow(`SELECT balance FROM accounts WHERE id=$1`, t.SenderId)
	if err := row.Scan(&sndBal); err != nil {
		_ = tx.Rollback()
		return err
	}

	if sndBal < t.Amount {
		_ = tx.Rollback()
		return errors.New("err not enough money")
	}

	row = tx.QueryRow(`SELECT balance FROM accounts WHERE id=$1`, t.ReceiverId)
	if err := row.Scan(&rcvBal); err != nil {
		_ = tx.Rollback()
		return err
	}

	updSndBal := sndBal - t.Amount
	updRcvBal := rcvBal + t.Amount

	result, err := tx.Exec(`UPDATE accounts SET balance=$1 WHERE id=$2`, updSndBal, t.SenderId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if affected, _ := result.RowsAffected(); affected != 1 {
		return errors.New("err invalid affected rows count")
	}

	result, err = tx.Exec(`UPDATE accounts SET balance=$1 WHERE id=$2`, updRcvBal, t.ReceiverId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if affected, _ := result.RowsAffected(); affected != 1 {
		return errors.New("err invalid affected rows count")
	}

	result, err = tx.Exec(`INSERT INTO transactions (id, sender_id, receiver_id, amount) 
		VALUES ($1, $2, $3, $4)`, t.Id, t.SenderId, t.ReceiverId, t.Amount)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if affected, _ := result.RowsAffected(); affected != 1 {
		return errors.New("err invalid affected rows count")
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
