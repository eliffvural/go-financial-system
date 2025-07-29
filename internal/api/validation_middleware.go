package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ValidationMiddleware, JSON request'lerini validate eder
func ValidationMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Sadece POST/PUT request'lerini validate et
		if r.Method == "POST" || r.Method == "PUT" {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Content-Type application/json olmalı"))
				return
			}

			// Request body'sini kontrol et
			if r.Body == nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Request body gerekli"))
				return
			}
		}

		next(w, r)
	}
}

// ErrorHandlingMiddleware, panic'leri yakalar ve hata response'u döner
func ErrorHandlingMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error (gerçek implementasyonda logger kullanılmalı)
				fmt.Printf("Panic yakalandı: %v\n", err)
				
				// JSON error response
				errorResponse := map[string]interface{}{
					"error":   "Internal Server Error",
					"message": "Beklenmeyen bir hata oluştu",
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errorResponse)
			}
		}()

		next(w, r)
	}
}

// PerformanceMonitoringMiddleware, request süresini ölçer
func PerformanceMonitoringMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Response writer'ı wrap et
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: 200}
		
		next(wrappedWriter, r)
		
		// Performance metrics
		duration := time.Since(start)
		fmt.Printf("[PERFORMANCE] %s %s - %d - %v\n", 
			r.Method, r.URL.Path, wrappedWriter.statusCode, duration)
	}
}

// responseWriter, status code'u yakalamak için wrapper
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestSizeMiddleware, request boyutunu kontrol eder
func RequestSizeMiddleware(maxSize int64) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxSize {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				w.Write([]byte("Request boyutu çok büyük"))
				return
			}
			next(w, r)
		}
	}
}

// TimeoutMiddleware, request timeout'u kontrol eder
func TimeoutMiddleware(timeout time.Duration) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Context timeout ekle
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			
			// Channel ile timeout kontrolü
			done := make(chan bool, 1)
			go func() {
				next(w, r.WithContext(ctx))
				done <- true
			}()
			
			select {
			case <-done:
				// Request tamamlandı
			case <-ctx.Done():
				w.WriteHeader(http.StatusRequestTimeout)
				w.Write([]byte("Request timeout"))
			}
		}
	}
} 