package main

import (
	"fmt"
	"gofinancialsystem/internal/domain"
	"gofinancialsystem/internal/repository"
	"gofinancialsystem/internal/service"
	"log"
)

func testSystem() {
	fmt.Println("=== Go Financial System Test ===")

	// 1. Repository'leri oluştur
	userRepo := repository.NewUserRepository()

	// 2. Service'leri oluştur
	userService := service.NewUserService(userRepo)

	// 3. Test: Kullanıcı kaydı
	fmt.Println("\n1. Kullanıcı kaydı testi:")
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "123456",
		Role:     "user",
	}

	if err := userService.Register(user); err != nil {
		log.Fatalf("Kullanıcı kaydı hatası: %v", err)
	}
	fmt.Printf("Kullanıcı başarıyla kaydedildi: ID=%d, Username=%s\n", user.ID, user.Username)

	// 4. Test: Kullanıcı girişi
	fmt.Println("\n2. Kullanıcı girişi testi:")
	authUser, err := userService.Authenticate("testuser", "123456")
	if err != nil {
		log.Fatalf("Kullanıcı girişi hatası: %v", err)
	}
	fmt.Printf("Kullanıcı başarıyla giriş yaptı: ID=%d, Role=%s\n", authUser.ID, authUser.Role)

	// 5. Test: Yetkilendirme
	fmt.Println("\n3. Yetkilendirme testi:")
	if userService.Authorize(authUser, "user") {
		fmt.Println("Kullanıcı 'user' rolüne sahip")
	}
	if !userService.Authorize(authUser, "admin") {
		fmt.Println("Kullanıcı 'admin' rolüne sahip değil")
	}

	// 6. Test: Domain model validasyonu
	fmt.Println("\n4. Domain model validasyonu testi:")
	invalidUser := &domain.User{
		Username: "",
		Email:    "invalid-email",
		Password: "",
		Role:     "",
	}
	if err := invalidUser.Validate(); err != nil {
		fmt.Printf("Validasyon hatası (beklenen): %v\n", err)
	}

	fmt.Println("\n=== Test tamamlandı! ===")
}
