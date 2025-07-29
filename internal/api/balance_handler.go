package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"gofinancialsystem/internal/domain"
)

// BalanceHandler, balance işlemleri için servisleri tutar
type BalanceHandler struct {
	BalanceService domain.BalanceService
}

// Mevcut bakiye (GET /api/v1/balances/current)
func (h *BalanceHandler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Kullanıcı ID gerekli"))
		return
	}
	
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz kullanıcı ID"))
		return
	}
	
	balance, err := h.BalanceService.GetBalance(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Bakiye bulunamadı"))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

// Bakiye geçmişi (GET /api/v1/balances/historical)
func (h *BalanceHandler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Kullanıcı ID gerekli"))
		return
	}
	
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz kullanıcı ID"))
		return
	}
	
	history, err := h.BalanceService.GetBalanceHistory(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Bakiye geçmişi alınamadı"))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}

// Belirli zamandaki bakiye (GET /api/v1/balances/at-time)
func (h *BalanceHandler) GetBalanceAtTime(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	timestampStr := r.URL.Query().Get("timestamp")
	
	if userIDStr == "" || timestampStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Kullanıcı ID ve timestamp gerekli"))
		return
	}
	
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz kullanıcı ID"))
		return
	}
	
	// Timestamp'i parse et (RFC3339 formatında)
	targetTime, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz timestamp formatı"))
		return
	}
	
	balance, err := h.BalanceService.GetBalanceAtTime(userID, targetTime)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Belirtilen zamanda bakiye bulunamadı"))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

// Bakiye hesaplama (GET /api/v1/balances/calculate)
func (h *BalanceHandler) CalculateBalance(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Kullanıcı ID gerekli"))
		return
	}
	
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz kullanıcı ID"))
		return
	}
	
	amount, err := h.BalanceService.CalculateBalance(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Bakiye hesaplanamadı"))
		return
	}
	
	response := map[string]interface{}{
		"user_id": userID,
		"amount":  amount,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
} 