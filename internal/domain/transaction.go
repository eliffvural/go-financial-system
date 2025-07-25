package domain

import (
	"errors"
	"time"
)

type TransactionStatus string

type TransactionType string

const (
	TransactionPending   TransactionStatus = "pending"
	TransactionCompleted TransactionStatus = "completed"
	TransactionFailed    TransactionStatus = "failed"

	TransactionDeposit  TransactionType = "deposit"
	TransactionWithdraw TransactionType = "withdraw"
	TransactionTransfer TransactionType = "transfer"
)

type Transaction struct {
	ID         int64             `json:"id"`
	FromUserID *int64            `json:"from_user_id,omitempty"`
	ToUserID   *int64            `json:"to_user_id,omitempty"`
	Amount     float64           `json:"amount"`
	Type       TransactionType   `json:"type"`
	Status     TransactionStatus `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
}

func (t *Transaction) Complete() error {
	if t.Status != TransactionPending {
		return errors.New("sadece bekleyen işlemler tamamlanabilir")
	}
	t.Status = TransactionCompleted
	return nil
}

func (t *Transaction) Fail() error {
	if t.Status != TransactionPending {
		return errors.New("sadece bekleyen işlemler başarısız yapılabilir")
	}
	t.Status = TransactionFailed
	return nil
}
