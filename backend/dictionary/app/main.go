package main

import (
	"github.com/vovancho/lingua-cat-go/dictionary/internal/validator"

	_dictionaryGrpc "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc"
	pb "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc/gen"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/auth"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/config"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/db"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/response"
	"google.golang.org/grpc"
	"log/slog"
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
		panic(err)
	}

	dbConn, err := db.NewDB(cfg.DBDSN)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		panic(err)
	}
	defer dbConn.Close()

	validate, trans, err := validator.NewValidator()
	if err != nil {
		slog.Error("Failed to initialize validator", "error", err)
		panic(err)
	}

	authService, err := auth.NewAuthService(cfg.AuthPublicKeyPath)
	if err != nil {
		slog.Error("Failed to initialize auth service", "error", err)
		panic(err)
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

	go func() {
		lis, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			slog.Error("Failed to listen for gRPC", "error", err)
			panic(err)
		}
		slog.Info("gRPC server is listening", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("Failed to serve gRPC", "error", err)
			panic(err)
		}
	}()

	slog.Info("HTTP server is listening", "port", cfg.HTTPPort)
	if err := httpServer.ListenAndServe(); err != nil {
		slog.Error("Failed to serve HTTP", "error", err)
		panic(err)
	}
}
