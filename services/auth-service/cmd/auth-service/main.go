package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PomozD/monitoring-platform/services/auth-service/internal/config"
	"github.com/PomozD/monitoring-platform/services/auth-service/internal/database"
	"github.com/PomozD/monitoring-platform/services/auth-service/internal/handler"
)

func main() {
	logger := log.New(os.Stdout, "auth-service: ", log.LstdFlags)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("load config: %v", err)
	}

	ctx := context.Background()

	dbPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("connect to postgres: %v", err)
	}
	defer dbPool.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)

	go func() {
		logger.Printf("HTTP server listening on %s", cfg.AppPort)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverErr:
		logger.Fatalf("HTTP server error: %v", err)

	case <-signalCtx.Done():
		logger.Println("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("HTTP server shutdown error: %v", err)
	}

	logger.Println("auth-service stopped")
}
