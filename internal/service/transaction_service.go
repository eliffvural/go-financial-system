package service

import (
	"errors"
	"gofinancialsystem/internal/domain"
	"sync"
)

type TransactionRepositoryImpl struct {
	transactions map[int64]*domain.Transaction
	mu           sync.RWMutex
	nextID       int64
}

func NewTransactionRepository() *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{
		transactions: make(map[int64]*domain.Transaction),
		nextID:       1,
	}
}

func (r *TransactionRepositoryImpl) Create(tx *domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx.ID = r.nextID
	r.nextID++
	r.transactions[tx.ID] = tx
	return nil
}

func (r *TransactionRepositoryImpl) FindByID(id int64) (*domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if tx, exists := r.transactions[id]; exists {
		return tx, nil
	}
	return nil, errors.New("işlem bulunamadı")
}

func (r *TransactionRepositoryImpl) ListByUser(userID int64) ([]*domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Transaction
	for _, tx := range r.transactions {
		if (tx.FromUserID != nil && *tx.FromUserID == userID) || (tx.ToUserID != nil && *tx.ToUserID == userID) {
			result = append(result, tx)
		}
	}
	return result, nil
}
