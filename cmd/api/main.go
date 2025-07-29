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
	balanceHandler := &api.BalanceHandler{BalanceService: balanceService}

	// Router oluştur
	router := api.NewRouter()

	// Middleware'leri ekle (sıralama önemli)
	router.Use(api.ErrorHandlingMiddleware)
	router.Use(api.PerformanceMonitoringMiddleware)
	router.Use(api.LoggingMiddleware)
	router.Use(api.CORSMiddleware)
	router.Use(api.SecurityHeadersMiddleware)
	router.Use(api.RateLimitMiddleware)
	router.Use(api.ValidationMiddleware)
	router.Use(api.RequestSizeMiddleware(1024 * 1024)) // 1MB limit

	// Health endpoint
	router.Handle("GET", "/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Auth endpointleri (auth middleware yok)
	router.Handle("POST", "/api/v1/auth/register", authHandler.Register)
	router.Handle("POST", "/api/v1/auth/login", authHandler.Login)

	// User Management endpointleri (auth gerekli)
	router.Handle("GET", "/api/v1/users", api.AuthMiddleware(userHandler.ListUsers))
	router.Handle("GET", "/api/v1/users/get", api.AuthMiddleware(userHandler.GetUser))
	router.Handle("PUT", "/api/v1/users/update", api.AuthMiddleware(userHandler.UpdateUser))
	router.Handle("DELETE", "/api/v1/users/delete", api.AdminOnlyMiddleware(userHandler.DeleteUser))

	// Transaction endpointleri (auth gerekli)
	router.Handle("POST", "/api/v1/transactions/credit", api.AuthMiddleware(transactionHandler.Credit))
	router.Handle("POST", "/api/v1/transactions/debit", api.AuthMiddleware(transactionHandler.Debit))
	router.Handle("POST", "/api/v1/transactions/transfer", api.AuthMiddleware(transactionHandler.Transfer))
	router.Handle("GET", "/api/v1/transactions/history", api.AuthMiddleware(transactionHandler.GetHistory))
	router.Handle("GET", "/api/v1/transactions/get", api.AuthMiddleware(transactionHandler.GetTransaction))

	// Balance endpointleri (auth gerekli)
	router.Handle("GET", "/api/v1/balances/current", api.AuthMiddleware(balanceHandler.GetCurrentBalance))
	router.Handle("GET", "/api/v1/balances/historical", api.AuthMiddleware(balanceHandler.GetBalanceHistory))
	router.Handle("GET", "/api/v1/balances/at-time", api.AuthMiddleware(balanceHandler.GetBalanceAtTime))
	router.Handle("GET", "/api/v1/balances/calculate", api.AuthMiddleware(balanceHandler.CalculateBalance))

	// Sunucuyu başlat
	api.StartServer(":8080", router)
}
