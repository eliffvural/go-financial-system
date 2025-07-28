package api

import (
	"encoding/json"
	"fmt"
	"gofinancialsystem/internal/domain"
	"net/http"
	"strconv"
)

// UserHandler, kullanıcı yönetimi işlemleri için servisleri tutar
type UserHandler struct {
	UserService domain.UserService
}

// Tüm kullanıcıları listeler (GET /api/v1/users)
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Basit bir örnek - gerçek uygulamada repository'den tüm kullanıcıları çekersiniz
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Kullanıcı listesi (henüz implement edilmedi)"))
}

// Belirli bir kullanıcıyı getirir (GET /api/v1/users/{id})
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// URL'den ID'yi çıkar (basit implementasyon)
	// Gerçek uygulamada path parameter parsing yaparsınız
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Kullanıcı ID gerekli"))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz kullanıcı ID"))
		return
	}

	user, err := h.UserService.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Kullanıcı bulunamadı"))
		return
	}

	// JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Kullanıcı bilgilerini günceller (PUT /api/v1/users/{id})
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz istek"))
		return
	}

	// Basit bir örnek - gerçek uygulamada kullanıcıyı güncelleme işlemi yaparsınız
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Kullanıcı güncellendi"))
}

// Kullanıcıyı siler (DELETE /api/v1/users/{id})
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Kullanıcı ID gerekli"))
		return
	}

	// Basit bir örnek - gerçek uygulamada kullanıcı silme işlemi yaparsınız
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Kullanıcı %s silindi", idStr)))
}
