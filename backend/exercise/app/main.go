package main

import (
	"context"
	"github.com/google/wire"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Инициализация приложения с помощью Wire
	app, err := wire.InitializeApp()
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer app.DB.Close() // Закрытие соединения с базой данных

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("HTTP server is listening", "port", app.Config.HTTPPort)
		if err := app.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to serve HTTP", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("Initiating graceful shutdown...")

	// Завершение HTTP-сервера с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := app.HTTPServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Failed to shutdown HTTP server", "error", err)
	} else {
		slog.Info("HTTP server stopped")
	}

	wg.Wait()
	slog.Info("All servers stopped")
}
