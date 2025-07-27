package service

import (
	"gofinancialsystem/internal/domain"
	"sync"
	"time"
)

// BalanceHistory, bakiye geçmişini takip etmek için kullanılır
type BalanceHistory struct {
	UserID    int64     `json:"user_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// BalanceServiceImpl, BalanceService arayüzünün gerçek implementasyonudur
type BalanceServiceImpl struct {
	balanceRepo domain.BalanceRepository   // Bakiye veritabanı işlemleri için repository
	mu          sync.RWMutex               // Thread-safe işlemler için RWMutex
	history     map[int64][]BalanceHistory // Kullanıcı bakiye geçmişi (cache)
}

// Yeni bir BalanceServiceImpl oluşturur
func NewBalanceService(balanceRepo domain.BalanceRepository) *BalanceServiceImpl {
	return &BalanceServiceImpl{
		balanceRepo: balanceRepo,
		history:     make(map[int64][]BalanceHistory),
	}
}

// Kullanıcının bakiyesini thread-safe şekilde döndürür
func (s *BalanceServiceImpl) GetByUserID(userID int64) (*domain.Balance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.balanceRepo.GetByUserID(userID)
}

// Kullanıcının bakiyesini thread-safe şekilde günceller
func (s *BalanceServiceImpl) Update(userID int64, amount float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Bakiye güncelleme
	if err := s.balanceRepo.Update(userID, amount); err != nil {
		return err
	}

	// Geçmiş kaydı ekle
	historyEntry := BalanceHistory{
		UserID:    userID,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	s.history[userID] = append(s.history[userID], historyEntry)

	// Geçmiş kayıtlarını optimize et (son 100 kayıt tut)
	if len(s.history[userID]) > 100 {
		s.history[userID] = s.history[userID][len(s.history[userID])-100:]
	}

	return nil
}

// Kullanıcının bakiye geçmişini döndürür
func (s *BalanceServiceImpl) GetHistory(userID int64) []BalanceHistory {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if history, exists := s.history[userID]; exists {
		return history
	}
	return []BalanceHistory{}
}

// Bakiye hesaplama optimizasyonu: belirli bir tarihten sonraki değişiklikleri hesaplar
func (s *BalanceServiceImpl) CalculateBalanceFromDate(userID int64, fromDate time.Time) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var totalChange float64
	if history, exists := s.history[userID]; exists {
		for _, entry := range history {
			if entry.Timestamp.After(fromDate) {
				totalChange += entry.Amount
			}
		}
	}
	return totalChange
}
