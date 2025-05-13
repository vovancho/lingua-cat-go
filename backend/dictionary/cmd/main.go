package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vovancho/lingua-cat-go/dictionary/internal/wire"
)

// @title     Документация сервиса Dictionary
// @version   1.0
// @host      api.lingua-cat-go.localhost
// @BasePath  /dictionary
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Для авторизации используйте Bearer токен, полученный от Keycloak:
// @description 1. Получите access_token:
// @description    curl -X POST http://keycloak.localhost/realms/lingua-cat-go/protocol/openid-connect/token -H "Content-Type: application/x-www-form-urlencoded" -d 'grant_type=password&scope=openid&client_id=lingua-cat-go-dev&client_secret=GatPbS9gsEfplvCpiNitwBdmIRc0QqyQ&username=username&password=password'
// @description 2. Используйте access_token в заголовке Authorization: Bearer <token>
func main() {
	// Инициализация приложения с помощью Wire
	app, err := wire.InitializeApp()
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer app.DB.Close() // Закрытие соединения с базой данных

	// Завершение трейсера
	defer func() {
		if err := app.Tracer.Shutdown(context.Background()); err != nil {
			slog.Error("Failed to shutdown tracer provider", "error", err)
		}
	}()

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
