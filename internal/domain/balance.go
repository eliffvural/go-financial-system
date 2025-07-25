package domain

import (
	"sync"
	"time"
)

type Balance struct {
	UserID        int64     `json:"user_id"`
	Amount        float64   `json:"amount"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	mu            sync.Mutex
}

func (b *Balance) Add(amount float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Amount += amount
	b.LastUpdatedAt = time.Now()
}

func (b *Balance) Subtract(amount float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.Amount < amount {
		return false
	}
	b.Amount -= amount
	b.LastUpdatedAt = time.Now()
	return true
}
