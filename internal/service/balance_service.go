package service

import (
	"errors"
	"gofinancialsystem/internal/domain"
	"sync"
	"time"
)

// BalanceServiceImpl, BalanceService interface'ini implement eder
type BalanceServiceImpl struct {
	balanceRepo domain.BalanceRepository
	// Thread-safe balance cache
	balanceCache map[int64]*domain.Balance
	cacheMutex   sync.RWMutex
	// Historical balance tracking
	balanceHistory map[int64][]*domain.Balance
	historyMutex   sync.RWMutex
}

// NewBalanceService, yeni bir BalanceService instance'ı oluşturur
func NewBalanceService(balanceRepo domain.BalanceRepository) domain.BalanceService {
	return &BalanceServiceImpl{
		balanceRepo:    balanceRepo,
		balanceCache:   make(map[int64]*domain.Balance),
		balanceHistory: make(map[int64][]*domain.Balance),
	}
}

// GetBalance, kullanıcının mevcut bakiyesini getirir
func (s *BalanceServiceImpl) GetBalance(userID int64) (*domain.Balance, error) {
	// Önce cache'den kontrol et
	s.cacheMutex.RLock()
	if balance, exists := s.balanceCache[userID]; exists {
		s.cacheMutex.RUnlock()
		return balance, nil
	}
	s.cacheMutex.RUnlock()

	// Cache'de yoksa repository'den al
	balance, err := s.balanceRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Cache'e ekle
	s.cacheMutex.Lock()
	s.balanceCache[userID] = balance
	s.cacheMutex.Unlock()

	return balance, nil
}

// UpdateBalance, kullanıcının bakiyesini günceller
func (s *BalanceServiceImpl) UpdateBalance(userID int64, amount float64) error {
	// Thread-safe balance update
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Mevcut balance'ı al veya oluştur
	balance, exists := s.balanceCache[userID]
	if !exists {
		balance = &domain.Balance{
			UserID:        userID,
			Amount:        0,
			LastUpdatedAt: time.Now(),
		}
		s.balanceCache[userID] = balance
	}

	// Balance'ı güncelle
	balance.Add(amount)
	balance.LastUpdatedAt = time.Now()

	// Repository'yi güncelle
	if err := s.balanceRepo.Update(userID, amount); err != nil {
		return err
	}

	// Historical tracking
	s.historyMutex.Lock()
	if s.balanceHistory[userID] == nil {
		s.balanceHistory[userID] = make([]*domain.Balance, 0)
	}
	// Historical copy oluştur
	historicalBalance := &domain.Balance{
		UserID:        balance.UserID,
		Amount:        balance.Amount,
		LastUpdatedAt: balance.LastUpdatedAt,
	}
	s.balanceHistory[userID] = append(s.balanceHistory[userID], historicalBalance)
	s.historyMutex.Unlock()

	return nil
}

// GetBalanceHistory, kullanıcının bakiye geçmişini getirir
func (s *BalanceServiceImpl) GetBalanceHistory(userID int64) ([]*domain.Balance, error) {
	s.historyMutex.RLock()
	defer s.historyMutex.RUnlock()

	history, exists := s.balanceHistory[userID]
	if !exists {
		return []*domain.Balance{}, nil
	}

	return history, nil
}

// GetBalanceAtTime, belirli bir zamandaki bakiyeyi getirir (basit implementasyon)
func (s *BalanceServiceImpl) GetBalanceAtTime(userID int64, targetTime time.Time) (*domain.Balance, error) {
	s.historyMutex.RLock()
	defer s.historyMutex.RUnlock()

	history, exists := s.balanceHistory[userID]
	if !exists {
		return nil, errors.New("bakiye geçmişi bulunamadı")
	}

	// En yakın zamandaki balance'ı bul
	var closestBalance *domain.Balance
	var minDiff time.Duration

	for _, balance := range history {
		diff := targetTime.Sub(balance.LastUpdatedAt)
		if diff < 0 {
			diff = -diff
		}
		if closestBalance == nil || diff < minDiff {
			closestBalance = balance
			minDiff = diff
		}
	}

	if closestBalance == nil {
		return nil, errors.New("belirtilen zamanda bakiye bulunamadı")
	}

	return closestBalance, nil
}

// CalculateBalance, kullanıcının toplam bakiyesini hesaplar (optimizasyon için)
func (s *BalanceServiceImpl) CalculateBalance(userID int64) (float64, error) {
	balance, err := s.GetBalance(userID)
	if err != nil {
		return 0, err
	}
	return balance.Amount, nil
}
