package repository

import (
	"errors"
	"sync"
	"gofinancialsystem/internal/domain"
)

// UserRepositoryImpl, UserRepository arayüzünün basit implementasyonudur (test için)
type UserRepositoryImpl struct {
	users map[int64]*domain.User // In-memory storage
	mu    sync.RWMutex
	nextID int64
}

// Yeni bir UserRepositoryImpl oluşturur
func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{
		users: make(map[int64]*domain.User),
		nextID: 1,
	}
}

// Kullanıcı oluşturur
func (r *UserRepositoryImpl) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = user
	return nil
}

// ID ile kullanıcı bulur
func (r *UserRepositoryImpl) FindByID(id int64) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("kullanıcı bulunamadı")
}

// Kullanıcı adı ile kullanıcı bulur
func (r *UserRepositoryImpl) FindByUsername(username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("kullanıcı bulunamadı")
} 