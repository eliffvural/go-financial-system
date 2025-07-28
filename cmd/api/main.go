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
	userService := service.NewUserService(userRepo)

	// Auth handler'ı oluştur
	authHandler := &api.AuthHandler{UserService: userService}

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

	// Sunucuyu başlat
	api.StartServer(":8080", router)
}
