package processing

import (
	"fmt"
	"sync"
)

// BatchProcessor, birden fazla işlemi (ör. transaction) aynı anda işlemek için kullanılır
type BatchProcessor struct {
	wg sync.WaitGroup // Batch işlemlerinin bitmesini beklemek için
}

// İşleri eşzamanlı (concurrent) olarak işler
// Her iş için ayrı bir goroutine başlatır
func (bp *BatchProcessor) ProcessBatch(jobs []func()) {
	bp.wg.Add(len(jobs))
	for i, job := range jobs {
		go func(idx int, fn func()) {
			defer bp.wg.Done()
			fmt.Printf("Batch %d. iş başlatıldı\n", idx)
			fn()
			fmt.Printf("Batch %d. iş tamamlandı\n", idx)
		}(i, job)
	}
	bp.wg.Wait() // Tüm işler bitene kadar bekle
}
