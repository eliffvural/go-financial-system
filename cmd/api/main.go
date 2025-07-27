package main

import (
	"gofinancialsystem/internal/api"
	"net/http"
)

func main() {
	// Router oluştur
	router := api.NewRouter()

	// Middleware'leri ekle
	router.Use(api.LoggingMiddleware)
	router.Use(api.CORSMiddleware)
	router.Use(api.SecurityHeadersMiddleware)
	router.Use(api.RateLimitMiddleware)

	// Basit bir health check endpoint'i ekle
	router.Handle("GET", "/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Sunucuyu başlat
	api.StartServer(":8080", router)
}
