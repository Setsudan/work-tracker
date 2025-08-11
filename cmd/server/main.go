package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"work-tracker/internal/config"
	"work-tracker/internal/httpserver"
	"work-tracker/internal/scheduler"
	"work-tracker/internal/secret"
	"work-tracker/internal/store"
)

func main() {
	cfg := config.Load()

	sec, err := secret.New(cfg.DataEncKeyB64)
	if err != nil {
		log.Fatalf("invalid DATA_ENCRYPTION_KEY: %v", err)
	}

	redisClient, err := store.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer func() { _ = redisClient.Close() }()

	stores := store.NewStores(redisClient, sec)

	srv := httpserver.NewServer(cfg, stores)

	cr := scheduler.New(cfg, stores)
	cr.Start()
	defer cr.Stop()

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr(),
		Handler:           srv,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr())
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
}
