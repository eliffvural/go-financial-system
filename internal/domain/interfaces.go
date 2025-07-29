package domain

import "time"

// UserService, TransactionService, BalanceService gibi servisler için temel arayüzler
type UserService interface {
	Register(user *User) error
	Authenticate(username, password string) (*User, error)
	GetByID(id int64) (*User, error)
}

type TransactionService interface {
	Create(tx *Transaction) error
	GetByID(id int64) (*Transaction, error)
	ListByUser(userID int64) ([]*Transaction, error)
	Credit(userID int64, amount float64) error
	Debit(userID int64, amount float64) error
	Transfer(fromUserID, toUserID int64, amount float64) error
}

type BalanceService interface {
	GetBalance(userID int64) (*Balance, error)
	UpdateBalance(userID int64, amount float64) error
	GetBalanceHistory(userID int64) ([]*Balance, error)
	GetBalanceAtTime(userID int64, targetTime time.Time) (*Balance, error)
	CalculateBalance(userID int64) (float64, error)
}

// Repository arayüzleri
type UserRepository interface {
	Create(user *User) error
	FindByID(id int64) (*User, error)
	FindByUsername(username string) (*User, error)
}

type TransactionRepository interface {
	Create(tx *Transaction) error
	FindByID(id int64) (*Transaction, error)
	ListByUser(userID int64) ([]*Transaction, error)
}

type BalanceRepository interface {
	GetByUserID(userID int64) (*Balance, error)
	Update(userID int64, amount float64) error
}
