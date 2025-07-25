package domain

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
}

type BalanceService interface {
	GetByUserID(userID int64) (*Balance, error)
	Update(userID int64, amount float64) error
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
