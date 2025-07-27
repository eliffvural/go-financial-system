package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// CORS için varsayılan başlıklar
var defaultCORSHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
	"Access-Control-Allow-Headers":     "Content-Type, Authorization",
	"Access-Control-Allow-Credentials": "true",
}

// CORS middleware'i: CORS başlıklarını ekler
func CORSMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for k, v := range defaultCORSHeaders {
			w.Header().Set(k, v)
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

// Güvenlik başlıkları middleware'i
func SecurityHeadersMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next(w, r)
	}
}

// Rate limiting middleware'i (IP başına basit limit)
type rateLimiter struct {
	mu      sync.Mutex
	clients map[string]time.Time
}

var limiter = &rateLimiter{clients: make(map[string]time.Time)}

func RateLimitMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter.mu.Lock()
		last, exists := limiter.clients[ip]
		now := time.Now()
		if exists && now.Sub(last) < 500*time.Millisecond {
			limiter.mu.Unlock()
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, "Çok fazla istek! Lütfen bekleyin.")
			return
		}
		limiter.clients[ip] = now
		limiter.mu.Unlock()
		next(w, r)
	}
}

// Request logging middleware'i
func LoggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		dur := time.Since(start)
		fmt.Printf("%s %s %s %v\n", r.RemoteAddr, r.Method, r.URL.Path, dur)
	}
}
