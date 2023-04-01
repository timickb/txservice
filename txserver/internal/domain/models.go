package domain

import "github.com/google/uuid"

const (
	TxStatusHandling = 0
	TxStatusOk       = 1
	TxStatusFailed   = 2
)

type (
	TxStatus byte
)

type Account struct {
	Id string `json:"id,omitempty"`

	// Balance in minimal currency units
	Balance int64 `json:"balance,omitempty"`
}

type Transaction struct {
	Id         string   `json:"id,omitempty"`
	SenderId   string   `json:"sender_id,omitempty"`
	ReceiverId string   `json:"receiver_id,omitempty"`
	Amount     int64    `json:"amount,omitempty"`
	Status     TxStatus `json:"status,omitempty"`
}

type ClientRequest struct {
	SenderId   uuid.UUID `json:"sender_id,omitempty"`
	ReceiverId uuid.UUID `json:"receiver_id,omitempty"`
	Amount     int64     `json:"amount,omitempty"`
}
