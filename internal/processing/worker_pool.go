package processing

import (
	"fmt"
	"gofinancialsystem/internal/domain"
	"sync"
)

// TransactionJob, işlenmek üzere kuyruğa alınan transaction'ı temsil eder
type TransactionJob struct {
	Transaction *domain.Transaction // İşlenecek transaction
}

// WorkerPool, belirli sayıda worker ile transaction'ları işler
type WorkerPool struct {
	JobQueue   chan TransactionJob // Transaction iş kuyruğu (channel)
	NumWorkers int                 // Worker sayısı
	wg         sync.WaitGroup      // Worker'ların bitişini beklemek için
}

// Yeni bir worker pool oluşturur
func NewWorkerPool(numWorkers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		JobQueue:   make(chan TransactionJob, queueSize),
		NumWorkers: numWorkers,
	}
}

// Worker pool'u başlatır ve worker'ları çalıştırır
func (wp *WorkerPool) Start(processFunc func(TransactionJob)) {
	for i := 0; i < wp.NumWorkers; i++ {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			for job := range wp.JobQueue {
				fmt.Printf("Worker %d: Transaction işleniyor: %v\n", workerID, job.Transaction.ID)
				processFunc(job)
			}
		}(i)
	}
}

// Kuyruğa yeni bir transaction ekler
func (wp *WorkerPool) Enqueue(job TransactionJob) {
	wp.JobQueue <- job
}

// Tüm işlerin bitmesini bekler ve worker'ları kapatır
func (wp *WorkerPool) Stop() {
	close(wp.JobQueue)
	wp.wg.Wait()
}
