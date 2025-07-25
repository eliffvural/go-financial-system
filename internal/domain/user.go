package domain

import (
	"errors"
	"regexp"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

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

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
