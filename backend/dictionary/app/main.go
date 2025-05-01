package main

import (
	"context"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/translator"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/validator"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_dictionaryGrpc "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc"
	pb "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc/gen"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/auth"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/config"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/db"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/response"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"

	_dictionaryHttp "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/http"
	"github.com/vovancho/lingua-cat-go/dictionary/dictionary/repository/postgres"
	"github.com/vovancho/lingua-cat-go/dictionary/dictionary/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	dbConn, err := db.NewDB(cfg.DBDSN)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	trans, err := translator.NewTranslator()
	if err != nil {
		slog.Error("Failed to initialize translator", "error", err)
		os.Exit(1)
	}

	validate, err := validator.NewValidator(trans)
	if err != nil {
		slog.Error("Failed to initialize validator", "error", err)
		os.Exit(1)
	}

	authService, err := auth.NewAuthService(cfg.AuthPublicKeyPath)
	if err != nil {
		slog.Error("Failed to initialize auth service", "error", err)
		os.Exit(1)
	}

	dictionaryRepo := postgres.NewPostgresDictionaryRepository(dbConn)
	dictionaryUcase := usecase.NewDictionaryUseCase(dictionaryRepo, validate, time.Duration(cfg.Timeout)*time.Second)

	router := http.NewServeMux()
	_dictionaryHttp.NewDictionaryHandler(router, validate, dictionaryUcase)
	httpServer := http.Server{
		Addr:    cfg.HTTPPort,
		Handler: response.ErrorMiddleware(authService.AuthMiddleware(router), trans),
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authService.AuthInterceptor))
	pb.RegisterDictionaryServiceServer(grpcServer, _dictionaryGrpc.NewDictionaryHandler(validate, dictionaryUcase))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			slog.Error("Failed to listen for gRPC", "error", err)
			stop()
		}
		slog.Info("gRPC server is listening", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("Failed to serve gRPC", "error", err)
			stop()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("HTTP server is listening", "port", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to serve HTTP", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("Initiating graceful shutdown...")

	// Завершение HTTP-сервера с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Failed to shutdown HTTP server", "error", err)
	} else {
		slog.Info("HTTP server stopped")
	}

	// Завершение gRPC-сервера
	grpcServer.GracefulStop()
	slog.Info("gRPC server stopped")

	wg.Wait()
	slog.Info("All servers stopped")
}
