package main

import (
	"context"
	"github.com/vovancho/lingua-cat-go/analytics/internal/wire"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// @title     Документация сервиса Analytics
// @version   1.0
// @host      api.lingua-cat-go.localhost
// @BasePath  /analytics
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case msg, ok := <-app.ConsumerMessages:
				if !ok {
					slog.Info("Kafka message channel closed")
					return
				}
				slog.Info("Received message:", "uuid", msg.UUID)

				if err := app.ConsumerHandler.Handle(msg); err != nil {
					slog.Error("Failed to handle message", "error", err)

					msg.Nack()
					continue
				}

				msg.Ack()
			case <-ctx.Done():
				slog.Info("Stopping Kafka message processing")
				return
			}
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

	// Закрытие subscriber
	if err := app.Consumer.Close(); err != nil {
		slog.Error("Error closing subscriber", "error", err)
	}

	// Ожидание завершения горутин с таймаутом
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("All servers stopped")
	case <-time.After(15 * time.Second):
		slog.Error("Timeout waiting for servers to stop")
	}
}
