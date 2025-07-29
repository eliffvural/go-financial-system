package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"gofinancialsystem/internal/domain"
)

// AuthMiddleware, JWT token kontrolü yapar
func AuthMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorization header'ını kontrol et
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization header gerekli"))
			return
		}

		// Bearer token formatını kontrol et
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Geçersiz token formatı"))
			return
		}

		// Token'ı çıkar
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Token boş olamaz"))
			return
		}

		// Token'ı doğrula (basit implementasyon - gerçek JWT validation yapılmalı)
		userID, err := validateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Geçersiz token"))
			return
		}

		// User ID'yi context'e ekle
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
	}
}

// RoleMiddleware, belirli roller için erişim kontrolü yapar
func RoleMiddleware(requiredRoles ...string) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// User ID'yi context'ten al
			userID, ok := r.Context().Value("user_id").(int64)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Kullanıcı kimliği bulunamadı"))
				return
			}

			// User'ı getir (gerçek implementasyonda database'den alınmalı)
			user, err := getUserByID(userID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Kullanıcı bulunamadı"))
				return
			}

			// Role kontrolü
			hasRequiredRole := false
			for _, requiredRole := range requiredRoles {
				if user.Role == requiredRole {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Bu işlem için yetkiniz yok"))
				return
			}

			next(w, r)
		}
	}
}

// AdminOnlyMiddleware, sadece admin kullanıcılar için
func AdminOnlyMiddleware(next HandlerFunc) HandlerFunc {
	return RoleMiddleware("admin")(next)
}

// UserOrAdminMiddleware, user veya admin kullanıcılar için
func UserOrAdminMiddleware(next HandlerFunc) HandlerFunc {
	return RoleMiddleware("user", "admin")(next)
}

// Basit token validation (gerçek implementasyonda JWT kullanılmalı)
func validateToken(token string) (int64, error) {
	// Bu basit implementasyon - gerçek JWT validation yapılmalı
	// Şimdilik token'ı user ID olarak kabul ediyoruz
	var userID int64
	_, err := fmt.Sscanf(token, "%d", &userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Basit user lookup (gerçek implementasyonda database'den alınmalı)
func getUserByID(userID int64) (*domain.User, error) {
	// Bu basit implementasyon - gerçek database lookup yapılmalı
	return &domain.User{
		ID:       userID,
		Username: "user",
		Role:     "user", // Varsayılan role
	}, nil
} 