package main

import (
	"fmt"
	"gofinancialsystem/internal/domain"
	"gofinancialsystem/internal/processing"
	"gofinancialsystem/internal/repository"
	"gofinancialsystem/internal/service"
	"log"
	"sync"
	"time"
)

func testSystem() {
	fmt.Println("=== Go Financial System Full Test ===")

	// 1. Repository ve servisleri oluştur
	userRepo := repository.NewUserRepository()
	balanceRepo := repository.NewBalanceRepository()
	transactionRepo := repository.NewTransactionRepository()

	userService := service.NewUserService(userRepo)
	balanceService := service.NewBalanceService(balanceRepo)
	transactionService := service.NewTransactionService(transactionRepo, balanceRepo)

	// 2. Kullanıcı oluştur ve kaydet
	user1 := &domain.User{Username: "alice", Email: "alice@example.com", Password: "pass1", Role: "user"}
	user2 := &domain.User{Username: "bob", Email: "bob@example.com", Password: "pass2", Role: "user"}
	if err := userService.Register(user1); err != nil {
		log.Fatalf("Kullanıcı kaydı hatası: %v", err)
	}
	if err := userService.Register(user2); err != nil {
		log.Fatalf("Kullanıcı kaydı hatası: %v", err)
	}
	fmt.Printf("Kullanıcılar oluşturuldu: %s (ID=%d), %s (ID=%d)\n", user1.Username, user1.ID, user2.Username, user2.ID)

	// 3. Para yatırma ve çekme işlemleri
	fmt.Println("\n--- Para yatırma/çekme işlemleri ---")
	if err := transactionService.Credit(user1.ID, 1000); err != nil {
		log.Fatalf("Para yatırma hatası: %v", err)
	}
	if err := transactionService.Debit(user1.ID, 200); err != nil {
		log.Fatalf("Para çekme hatası: %v", err)
	}
	bal, _ := balanceService.GetByUserID(user1.ID)
	fmt.Printf("%s bakiyesi: %.2f\n", user1.Username, bal.Amount)

	// 4. Transfer işlemi
	fmt.Println("\n--- Transfer işlemi ---")
	if err := transactionService.Transfer(user1.ID, user2.ID, 300); err != nil {
		log.Fatalf("Transfer hatası: %v", err)
	}
	bal1, _ := balanceService.GetByUserID(user1.ID)
	bal2, _ := balanceService.GetByUserID(user2.ID)
	fmt.Printf("%s bakiyesi: %.2f, %s bakiyesi: %.2f\n", user1.Username, bal1.Amount, user2.Username, bal2.Amount)

	// 5. Worker pool ile toplu transaction işleme
	fmt.Println("\n--- Worker Pool ile toplu işlem ---")
	workerPool := processing.NewWorkerPool(3, 10)
	workerPool.Start(func(job processing.TransactionJob) {
		// Her transaction'ı işleyip transactionRepo'ya ekle
		tx := job.Transaction
		tx.Complete()
		transactionRepo.Create(tx)
	})
	for i := 0; i < 5; i++ {
		tx := &domain.Transaction{
			FromUserID: &user1.ID,
			ToUserID:   &user2.ID,
			Amount:     float64(10 * (i + 1)),
			Type:       domain.TransactionTransfer,
			Status:     domain.TransactionPending,
			CreatedAt:  time.Now(),
		}
		workerPool.Enqueue(processing.TransactionJob{Transaction: tx})
	}
	workerPool.Stop()
	fmt.Println("Toplu işlemler tamamlandı.")

	// 6. Batch processor ile toplu görev işleme
	fmt.Println("\n--- Batch Processor ile toplu görev ---")
	batch := processing.BatchProcessor{}
	var wg sync.WaitGroup
	jobs := []func(){
		func() { fmt.Println("Görev 1 çalıştı") },
		func() { fmt.Println("Görev 2 çalıştı") },
		func() { fmt.Println("Görev 3 çalıştı") },
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		batch.ProcessBatch(jobs)
	}()
	wg.Wait()
	fmt.Println("Batch görevler tamamlandı.")

	// 7. Transaction geçmişi ve validasyon
	fmt.Println("\n--- Transaction geçmişi ve validasyon ---")
	txs, _ := transactionRepo.ListByUser(user1.ID)
	fmt.Printf("%s kullanıcısının toplam işlemi: %d\n", user1.Username, len(txs))
	invalidUser := &domain.User{Username: "", Email: "invalid", Password: "", Role: ""}
	if err := invalidUser.Validate(); err != nil {
		fmt.Printf("Geçersiz kullanıcı validasyon hatası (beklenen): %v\n", err)
	}

	fmt.Println("\n=== Tüm testler başarıyla tamamlandı! ===")
}
