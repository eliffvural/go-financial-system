package service

import (
	"errors"
	"gofinancialsystem/internal/domain"
	"sync"
)

type BalanceRepositoryImpl struct {
	balances map[int64]*domain.Balance
	mu       sync.RWMutex
}

func NewBalanceRepository() *BalanceRepositoryImpl {
	return &BalanceRepositoryImpl{
		balances: make(map[int64]*domain.Balance),
	}
}

func (r *BalanceRepositoryImpl) GetByUserID(userID int64) (*domain.Balance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if bal, exists := r.balances[userID]; exists {
		return bal, nil
	}
	return nil, errors.New("bakiye bulunamadı")
}

func (r *BalanceRepositoryImpl) Update(userID int64, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	bal, ok := r.balances[userID]
	if !ok {
		bal = &domain.Balance{UserID: userID}
		r.balances[userID] = bal
	}
	if amount < 0 && bal.Amount < -amount {
		return errors.New("yetersiz bakiye")
	}
	bal.Amount += amount
	bal.LastUpdatedAt = bal.LastUpdatedAt.Add(0) // Sadece örnek için
	return nil
}
