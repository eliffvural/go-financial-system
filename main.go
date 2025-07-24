package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gofinancialsystem/internal/config"
	"gofinancialsystem/internal/logger"
)

func main() {
	// Config yükle
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Config yüklenemedi:", err)
		os.Exit(1)
	}

	// Logger başlat
	log := logger.New(cfg.Env)
	log.Info().Msg("Uygulama başlatıldı")

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Info().Msg("Kapanıyor...")
	// Burada kaynakları temizle (db, vs)
	// ...

	time.Sleep(1 * time.Second)
	log.Info().Msg("Çıkış yapıldı.")
}
