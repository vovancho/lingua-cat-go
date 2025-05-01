package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vovancho/lingua-cat-go/dictionary/internal/wire"
	"net"
	"net/http"
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
		lis, err := net.Listen("tcp", app.Config.GRPCPort)
		if err != nil {
			slog.Error("Failed to listen for gRPC", "error", err)
			stop()
		}
		slog.Info("gRPC server is listening", "port", app.Config.GRPCPort)
		if err := app.GRPCServer.Serve(lis); err != nil {
			slog.Error("Failed to serve gRPC", "error", err)
			stop()
		}
	}()

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

	// Завершение gRPC-сервера
	app.GRPCServer.GracefulStop()
	slog.Info("gRPC server stopped")

	wg.Wait()
	slog.Info("All servers stopped")
}
