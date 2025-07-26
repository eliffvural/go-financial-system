package service

import (
	"errors"
	"gofinancialsystem/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// UserServiceImpl, UserService arayüzünün gerçek implementasyonudur
type UserServiceImpl struct {
	userRepo domain.UserRepository // Kullanıcı veritabanı işlemleri için repository
}

// Yeni bir UserServiceImpl oluşturur
func NewUserService(userRepo domain.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{userRepo: userRepo}
}

// Kullanıcı kaydı (şifre hash'lenir)
func (s *UserServiceImpl) Register(user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("şifre hashlenemedi")
	}
	user.Password = string(hash)
	return s.userRepo.Create(user)
}

// Kullanıcı adı ve şifre ile giriş (authentication)
func (s *UserServiceImpl) Authenticate(username, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("kullanıcı bulunamadı")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("şifre hatalı")
	}
	return user, nil
}

// Kullanıcının rolünü kontrol eder (ör: admin yetkisi var mı?)
func (s *UserServiceImpl) Authorize(user *domain.User, requiredRole string) bool {
	return user.Role == requiredRole
}
