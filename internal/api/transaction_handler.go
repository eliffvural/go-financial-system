package api

import (
	"encoding/json"
	"fmt"
	"gofinancialsystem/internal/domain"
	"net/http"
	"strconv"
)

// TransactionHandler, transaction işlemleri için servisleri tutar
type TransactionHandler struct {
	TransactionService domain.TransactionService
	BalanceService     domain.BalanceService
}

// Para yatırma işlemi (POST /api/v1/transactions/credit)
func (h *TransactionHandler) Credit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz istek"))
		return
	}

	if err := h.TransactionService.Credit(req.UserID, req.Amount); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Para yatırma başarısız: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Para yatırma başarılı: %.2f", req.Amount)))
}

// Para çekme işlemi (POST /api/v1/transactions/debit)
func (h *TransactionHandler) Debit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz istek"))
		return
	}

	if err := h.TransactionService.Debit(req.UserID, req.Amount); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Para çekme başarısız: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Para çekme başarılı: %.2f", req.Amount)))
}

// Transfer işlemi (POST /api/v1/transactions/transfer)
func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromUserID int64   `json:"from_user_id"`
		ToUserID   int64   `json:"to_user_id"`
		Amount     float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz istek"))
		return
	}

	if err := h.TransactionService.Transfer(req.FromUserID, req.ToUserID, req.Amount); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Transfer başarısız: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Transfer başarılı: %.2f", req.Amount)))
}

// Transaction geçmişi (GET /api/v1/transactions/history)
func (h *TransactionHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
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

	transactions, err := h.TransactionService.ListByUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Transaction geçmişi alınamadı"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

// Belirli bir transaction'ı getir (GET /api/v1/transactions/{id})
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Transaction ID gerekli"))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Geçersiz transaction ID"))
		return
	}

	transaction, err := h.TransactionService.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Transaction bulunamadı"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}
