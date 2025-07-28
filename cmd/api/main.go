// The main function initializes repositories, services, handlers, and middleware for a financial system API in Go.
package main

import (
	"gofinancialsystem/internal/api"
	"gofinancialsystem/internal/repository"
	"gofinancialsystem/internal/service"
	"net/http"
)

func main() {
	// Repository ve servisleri başlat
	userRepo := repository.NewUserRepository()
	balanceRepo := repository.NewBalanceRepository()
	transactionRepo := repository.NewTransactionRepository()

	userService := service.NewUserService(userRepo)
	balanceService := service.NewBalanceService(balanceRepo)
	transactionService := service.NewTransactionService(transactionRepo, balanceRepo)

	// Handler'ları oluştur
	authHandler := &api.AuthHandler{UserService: userService}
	userHandler := &api.UserHandler{UserService: userService}
	transactionHandler := &api.TransactionHandler{
		TransactionService: transactionService,
		BalanceService:     balanceService,
	}

	// Router oluştur
	router := api.NewRouter()

	// Middleware'leri ekle
	router.Use(api.LoggingMiddleware)
	router.Use(api.CORSMiddleware)
	router.Use(api.SecurityHeadersMiddleware)
	router.Use(api.RateLimitMiddleware)

	// Health endpoint
	router.Handle("GET", "/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Auth endpointleri
	router.Handle("POST", "/api/v1/auth/register", authHandler.Register)
	router.Handle("POST", "/api/v1/auth/login", authHandler.Login)

	// User Management endpointleri
	router.Handle("GET", "/api/v1/users", userHandler.ListUsers)
	router.Handle("GET", "/api/v1/users/get", userHandler.GetUser)
	router.Handle("PUT", "/api/v1/users/update", userHandler.UpdateUser)
	router.Handle("DELETE", "/api/v1/users/delete", userHandler.DeleteUser)

	// Transaction endpointleri
	router.Handle("POST", "/api/v1/transactions/credit", transactionHandler.Credit)
	router.Handle("POST", "/api/v1/transactions/debit", transactionHandler.Debit)
	router.Handle("POST", "/api/v1/transactions/transfer", transactionHandler.Transfer)
	router.Handle("GET", "/api/v1/transactions/history", transactionHandler.GetHistory)
	router.Handle("GET", "/api/v1/transactions/get", transactionHandler.GetTransaction)

	// Sunucuyu başlat
	api.StartServer(":8080", router)
}
