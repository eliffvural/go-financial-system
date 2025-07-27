package api

import (
	"encoding/json"
	"gofinancialsystem/internal/domain"
	"gofinancialsystem/internal/service"
	"net/http"
)

// AuthHandler, auth işlemleri için servisleri tutar
type AuthHandler struct {
	UserService service.UserService
}

// Kullanıcı kaydı endpoint'i (POST /api/v1/auth/register)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz istek"))
		return
	}
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
	}
	if err := h.UserService.Register(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Kayıt başarısız: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Kayıt başarılı"))
}

// Kullanıcı girişi endpoint'i (POST /api/v1/auth/login)
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz istek"))
		return
	}
	user, err := h.UserService.Authenticate(req.Username, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Giriş başarısız: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Giriş başarılı, kullanıcı: " + user.Username))
}
