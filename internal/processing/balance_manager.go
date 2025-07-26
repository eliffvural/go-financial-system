package processing

import (
	"gofinancialsystem/internal/domain"
	"sync"
	"sync/atomic"
)

// BalanceManager, bakiyeleri thread-safe şekilde yöneten yapıdır
// Ayrıca transaction istatistiklerini atomic sayaçlarla tutar
type BalanceManager struct {
	balances map[int64]*domain.Balance // Kullanıcı ID -> Balance
	mu       sync.RWMutex              // Okuma/yazma için RWMutex

	totalTransactions uint64 // Toplam işlenen transaction sayısı (atomic)
	totalAmount       uint64 // Toplam işlenen miktar (atomic, kuruş cinsinden)
}

// Yeni bir BalanceManager oluşturur
func NewBalanceManager() *BalanceManager {
	return &BalanceManager{
		balances: make(map[int64]*domain.Balance),
	}
}

// Kullanıcıya ait bakiyeyi thread-safe şekilde döndürür
func (bm *BalanceManager) GetBalance(userID int64) *domain.Balance {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.balances[userID]
}

// Kullanıcıya ait bakiyeyi thread-safe şekilde günceller
func (bm *BalanceManager) UpdateBalance(userID int64, amount float64) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bal, ok := bm.balances[userID]
	if !ok {
		bal = &domain.Balance{UserID: userID}
		bm.balances[userID] = bal
	}
	bal.Add(amount)
	// Toplam transaction sayaçlarını güncelle
	atomic.AddUint64(&bm.totalTransactions, 1)
	atomic.AddUint64(&bm.totalAmount, uint64(amount*100)) // kuruş cinsinden
}

// Transaction istatistiklerini thread-safe şekilde döndürür
func (bm *BalanceManager) Stats() (totalTx uint64, totalAmount float64) {
	tx := atomic.LoadUint64(&bm.totalTransactions)
	amt := atomic.LoadUint64(&bm.totalAmount)
	return tx, float64(amt) / 100.0
}
