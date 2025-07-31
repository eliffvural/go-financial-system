package api

import (
	"encoding/json"
	"fmt"
	"gofinancialsystem/internal/domain"
	"net/http"
)

// AuthHandler, auth işlemleri için servisleri tutar
type AuthHandler struct {
	UserService domain.UserService
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

	response := map[string]interface{}{
		"message": "Kayıt başarılı",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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

	// Generate a simple token (in production, use JWT)
	token := fmt.Sprintf("token_%s_%d", user.Username, user.ID)

	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
