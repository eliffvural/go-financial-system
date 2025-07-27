package service

import (
	"errors"
	"gofinancialsystem/internal/domain"
)

// TransactionServiceImpl, TransactionService arayüzünün gerçek implementasyonudur
type TransactionServiceImpl struct {
	transactionRepo domain.TransactionRepository // Transaction veritabanı işlemleri için repository
	balanceRepo     domain.BalanceRepository     // Bakiye işlemleri için repository
}

// Yeni bir TransactionServiceImpl oluşturur
func NewTransactionService(txRepo domain.TransactionRepository, balRepo domain.BalanceRepository) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		transactionRepo: txRepo,
		balanceRepo:     balRepo,
	}
}

// Kullanıcıya kredi (para ekleme) işlemi
func (s *TransactionServiceImpl) Credit(userID int64, amount float64) error {
	tx := &domain.Transaction{
		ToUserID: &userID,
		Amount:   amount,
		Type:     domain.TransactionDeposit,
		Status:   domain.TransactionPending,
	}
	if err := s.balanceRepo.Update(userID, amount); err != nil {
		tx.Fail()
		return err
	}
	tx.Complete()
	return s.transactionRepo.Create(tx)
}

// Kullanıcıdan debit (para çekme) işlemi
func (s *TransactionServiceImpl) Debit(userID int64, amount float64) error {
	tx := &domain.Transaction{
		FromUserID: &userID,
		Amount:     amount,
		Type:       domain.TransactionWithdraw,
		Status:     domain.TransactionPending,
	}
	if err := s.balanceRepo.Update(userID, -amount); err != nil {
		tx.Fail()
		return err
	}
	tx.Complete()
	return s.transactionRepo.Create(tx)
}

// Hesaplar arası transfer işlemi
func (s *TransactionServiceImpl) Transfer(fromUserID, toUserID int64, amount float64) error {
	tx := &domain.Transaction{
		FromUserID: &fromUserID,
		ToUserID:   &toUserID,
		Amount:     amount,
		Type:       domain.TransactionTransfer,
		Status:     domain.TransactionPending,
	}
	// Önce gönderenin bakiyesinden düş
	if err := s.balanceRepo.Update(fromUserID, -amount); err != nil {
		tx.Fail()
		return err
	}
	// Sonra alıcının bakiyesine ekle
	if err := s.balanceRepo.Update(toUserID, amount); err != nil {
		// Rollback: Gönderenin bakiyesini geri yükle
		s.balanceRepo.Update(fromUserID, amount)
		tx.Fail()
		return err
	}
	tx.Complete()
	return s.transactionRepo.Create(tx)
}

// Transaction rollback işlemi (örnek: işlemi geri almak)
func (s *TransactionServiceImpl) Rollback(txID int64) error {
	tx, err := s.transactionRepo.FindByID(txID)
	if err != nil {
		return errors.New("işlem bulunamadı")
	}
	if tx.Status != domain.TransactionCompleted {
		return errors.New("sadece tamamlanmış işlemler geri alınabilir")
	}
	// Rollback işlemi: para hareketini tersine çevir
	switch tx.Type {
	case domain.TransactionDeposit:
		if tx.ToUserID != nil {
			s.balanceRepo.Update(*tx.ToUserID, -tx.Amount)
		}
	case domain.TransactionWithdraw:
		if tx.FromUserID != nil {
			s.balanceRepo.Update(*tx.FromUserID, tx.Amount)
		}
	case domain.TransactionTransfer:
		if tx.FromUserID != nil && tx.ToUserID != nil {
			s.balanceRepo.Update(*tx.FromUserID, tx.Amount)
			s.balanceRepo.Update(*tx.ToUserID, -tx.Amount)
		}
	}
	tx.Status = domain.TransactionFailed
	return nil
}
