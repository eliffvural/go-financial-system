package domain

import (
	"errors"
	"regexp"
)

// User, sistemdeki kullanıcıyı temsil eder
// Kullanıcı adı, e-posta, şifre ve rol bilgilerini içerir
type User struct {
	ID       int64  `json:"id"`       // Kullanıcının benzersiz ID'si
	Username string `json:"username"` // Kullanıcı adı
	Email    string `json:"email"`    // E-posta adresi
	Password string `json:"password"` // Şifre (hash'lenmiş olarak tutulmalı)
	Role     string `json:"role"`     // Kullanıcı rolü (ör: admin, user)
}

// Kullanıcı verisinin geçerli olup olmadığını kontrol eder
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username boş olamaz")
	}
	if u.Email == "" {
		return errors.New("email boş olamaz")
	}
	if !isValidEmail(u.Email) {
		return errors.New("geçersiz email formatı")
	}
	if u.Password == "" {
		return errors.New("şifre boş olamaz")
	}
	if u.Role == "" {
		return errors.New("rol boş olamaz")
	}
	return nil
}

// E-posta adresinin geçerli olup olmadığını kontrol eden yardımcı fonksiyon
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
